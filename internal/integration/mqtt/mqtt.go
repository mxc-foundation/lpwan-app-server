package mqtt

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"sync"
	"text/template"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	pb "github.com/brocaar/chirpstack-api/go/v3/as/integration"
	"github.com/brocaar/lorawan"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/mqttauth"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/marshaler"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/models"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	rs "github.com/mxc-foundation/lpwan-app-server/internal/modules/redis"

	"github.com/mxc-foundation/lpwan-app-server/internal/integration/types"
)

const (
	downlinkLockTTL = time.Millisecond * 100
)

// Integration implements a MQTT integration.
type Integration struct {
	marshaler          marshaler.Type
	conn               mqtt.Client
	dataDownChan       chan models.DataDownPayload
	wg                 sync.WaitGroup
	config             types.IntegrationMQTTConfig
	eventTopicTemplate *template.Template
	commandTopicRegexp *regexp.Regexp

	downlinkTopic string
	retainEvents  bool
}

// New creates a new MQTT integration.
func New(m marshaler.Type, conf types.IntegrationMQTTConfig) (*Integration, error) {
	var err error
	i := Integration{
		marshaler:    m,
		dataDownChan: make(chan models.DataDownPayload),
		config:       conf,
	}

	i.retainEvents = i.config.RetainEvents
	i.eventTopicTemplate, err = template.New("event").Parse(mqttauth.EventTopicTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse event topic template error: %v", err)
	}
	i.commandTopicRegexp, err = mqttauth.CompileRegexpFromTopicTemplate("command",
		mqttauth.CommandTopicTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "compile command topic regexp error")
	}

	// generate downlink topic matching all applications and devices
	topicBuffer := bytes.NewBuffer(nil)
	downTopicTemp, err := template.New("downlink").Parse(mqttauth.CommandTopicTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse command topic template error: %v", err)
	}
	if err = downTopicTemp.Execute(topicBuffer,
		struct {
			ApplicationID string
			DevEUI        string
			Type          string
		}{"+", "+", "down"}); err != nil {
		return nil, fmt.Errorf("create downlink topic error: %v", err)
	}
	i.downlinkTopic = topicBuffer.String()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(i.config.Server)
	opts.SetUsername(i.config.Username)
	opts.SetPassword(i.config.Password)
	opts.SetCleanSession(i.config.CleanSession)
	opts.SetClientID(i.config.ClientID)
	opts.SetOnConnectHandler(i.onConnected)
	opts.SetConnectionLostHandler(i.onConnectionLost)
	opts.SetMaxReconnectInterval(i.config.MaxReconnectInterval)

	tlsconfig, err := newTLSConfig(i.config.CACert, i.config.TLSCert, i.config.TLSKey)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"ca_cert":  i.config.CACert,
			"tls_cert": i.config.TLSCert,
			"tls_key":  i.config.TLSKey,
		}).Fatalf("error loading mqtt certificate files")
	}
	if tlsconfig != nil {
		opts.SetTLSConfig(tlsconfig)
	}

	log.WithField("server", i.config.Server).Info("integration/mqtt: connecting to mqtt broker")
	i.conn = mqtt.NewClient(opts)
	for {
		if token := i.conn.Connect(); token.Wait() && token.Error() != nil {
			log.Errorf("integration/mqtt: connecting to broker error, will retry in 2s: %s", token.Error())
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	return &i, nil
}

func newTLSConfig(cafile, certFile, certKeyFile string) (*tls.Config, error) {
	// Here are three valid options:
	//   - Only CA
	//   - TLS cert + key
	//   - CA, TLS cert + key

	if cafile == "" && certFile == "" && certKeyFile == "" {
		log.Info("integration/mqtt: TLS config is empty")
		return nil, nil
	}

	tlsConfig := &tls.Config{}

	// Import trusted certificates from CAfile.pem.
	if cafile != "" {
		cacert, err := ioutil.ReadFile(cafile)
		if err != nil {
			log.Errorf("integration/mqtt: couldn't load cafile: %s", err)
			return nil, err
		}
		certpool := x509.NewCertPool()
		certpool.AppendCertsFromPEM(cacert)

		tlsConfig.RootCAs = certpool // RootCAs = certs used to verify server cert.
	}

	// Import certificate and the key
	if certFile != "" || certKeyFile != "" {
		kp, err := tls.LoadX509KeyPair(certFile, certKeyFile) // here raises error when the pair of cert and key are invalid (e.g. either one is empty)
		if err != nil {
			log.Errorf("integration/mqtt: couldn't load MQTT TLS key pair: %s", err)
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{kp}
	}

	return tlsConfig, nil
}

// Close stops the handler.
func (i *Integration) Close() error {
	log.Info("integration/mqtt: closing handler")
	log.WithField("topic", i.downlinkTopic).Info("integration/mqtt: unsubscribing from tx topic")
	if token := i.conn.Unsubscribe(i.downlinkTopic); token.Wait() && token.Error() != nil {
		return fmt.Errorf("integration/mqtt: unsubscribe from %s error: %s", i.downlinkTopic, token.Error())
	}
	log.Info("integration/mqtt: handling last items in queue")
	i.wg.Wait()
	close(i.dataDownChan)
	return nil
}

// HandleUplinkEvent sends an UplinkEvent.
func (i *Integration) HandleUplinkEvent(ctx context.Context, _ models.Integration, vars map[string]string, payload pb.UplinkEvent) error {
	return i.publish(ctx, payload.ApplicationId, payload.DevEui, mqttauth.UplinkEvent, &payload)
}

// HandleJoinEvent sends a JoinEvent.
func (i *Integration) HandleJoinEvent(ctx context.Context, _ models.Integration, vars map[string]string, payload pb.JoinEvent) error {
	return i.publish(ctx, payload.ApplicationId, payload.DevEui, mqttauth.JoinEvent, &payload)
}

// HandleAckEvent sends an AckEvent.
func (i *Integration) HandleAckEvent(ctx context.Context, _ models.Integration, vars map[string]string, payload pb.AckEvent) error {
	return i.publish(ctx, payload.ApplicationId, payload.DevEui, mqttauth.AckEvent, &payload)
}

// HandleErrorEvent sends an ErrorEvent.
func (i *Integration) HandleErrorEvent(ctx context.Context, _ models.Integration, vars map[string]string, payload pb.ErrorEvent) error {
	return i.publish(ctx, payload.ApplicationId, payload.DevEui, mqttauth.ErrorEvent, &payload)
}

// HandleStatusEvent sends a StatusEvent.
func (i *Integration) HandleStatusEvent(ctx context.Context, _ models.Integration, vars map[string]string, payload pb.StatusEvent) error {
	return i.publish(ctx, payload.ApplicationId, payload.DevEui, mqttauth.StatusEvent, &payload)
}

// HandleLocationEvent sends a LocationEvent.
func (i *Integration) HandleLocationEvent(ctx context.Context, _ models.Integration, vars map[string]string, payload pb.LocationEvent) error {
	return i.publish(ctx, payload.ApplicationId, payload.DevEui, mqttauth.LocationEvent, &payload)
}

// HandleTxAckEvent sends a TxAckEvent.
func (i *Integration) HandleTxAckEvent(ctx context.Context, _ models.Integration, vars map[string]string, payload pb.TxAckEvent) error {
	return i.publish(ctx, payload.ApplicationId, payload.DevEui, mqttauth.TxAckEvent, &payload)
}

// HandleIntegrationEvent sends an IntegrationEvent.
func (i *Integration) HandleIntegrationEvent(ctx context.Context, _ models.Integration, vars map[string]string, payload pb.IntegrationEvent) error {
	return i.publish(ctx, payload.ApplicationId, payload.DevEui, mqttauth.IntegrationEvent, &payload)
}

func (i *Integration) publish(ctx context.Context, applicationID uint64, devEUIB []byte, eventType string, msg proto.Message) error {
	var devEUI lorawan.EUI64
	copy(devEUI[:], devEUIB)

	topic, err := i.getTopic(applicationID, devEUI, eventType)
	if err != nil {
		return errors.Wrap(err, "get topic error")
	}

	retain := i.getRetainEvents()

	b, err := marshaler.Marshal(i.marshaler, msg)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"dev_eui": devEUI,
		"retain":  retain,
		"topic":   topic,
		"qos":     i.config.QOS,
		"ctx_id":  ctx.Value(logging.ContextIDKey),
	}).Info("integration/mqtt: publishing event")
	if token := i.conn.Publish(topic, i.config.QOS, retain, b); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	mqttEventCounter(eventType).Inc()

	return nil
}

// DataDownChan returns the channel containing the received DataDownPayload.
func (i *Integration) DataDownChan() chan models.DataDownPayload {
	return i.dataDownChan
}

func (i *Integration) txPayloadHandler(mqttc mqtt.Client, msg mqtt.Message) {
	i.wg.Add(1)
	defer i.wg.Done()

	log.WithField("topic", msg.Topic()).Info("integration/mqtt: downlink event received")
	tv, err := mqttauth.GetTopicVariables(i.commandTopicRegexp, msg.Topic())
	if err != nil {
		log.WithError(err).Warning("integration/mqtt: get variables from topic error")
		return
	}

	var pl models.DataDownPayload
	dec := json.NewDecoder(bytes.NewReader(msg.Payload()))
	if err := dec.Decode(&pl); err != nil {
		log.WithFields(log.Fields{
			"data_base64": base64.StdEncoding.EncodeToString(msg.Payload()),
		}).Errorf("integration/mqtt: tx payload unmarshal error: %s", err)
		return
	}
	pl.ApplicationID, err = strconv.ParseInt(tv.ApplicationID, 10, 64)
	if err != nil {
		log.WithError(err).Warning("integration/mqtt: parse application id error")
		return
	}
	if err = pl.DevEUI.UnmarshalText([]byte(tv.DevEUI)); err != nil {
		log.WithError(err).Warning("integration/mqtt: get dev eui error")
		return
	}

	if pl.FPort == 0 || pl.FPort > 224 {
		log.WithFields(log.Fields{
			"topic":   msg.Topic(),
			"dev_eui": pl.DevEUI,
			"f_port":  pl.FPort,
		}).Error("integration/mqtt: fPort must be between 1 - 224")
		return
	}

	// Since with MQTT all subscribers will receive the downlink messages sent
	// by the application, the first instance receiving the message must lock it,
	// so that other instances can ignore the message.
	key := fmt.Sprintf("lora:as:downlink:lock:%d:%s", pl.ApplicationID, pl.DevEUI)
	set, err := rs.RedisClient().SetNX(key, "lock", downlinkLockTTL).Result()
	if err != nil {
		log.WithError(err).Error("integration/mqtt: acquire lock error")
		return
	}

	// If we could not set, it means it is already locked by an other process.
	if !set {
		return
	}

	mqttCommandCounter("down").Inc()

	i.dataDownChan <- pl
}

func (i *Integration) onConnected(mqttc mqtt.Client) {
	log.Info("integration/mqtt: connected to mqtt broker")
	for {
		log.WithFields(log.Fields{
			"topic": i.downlinkTopic,
			"qos":   i.config.QOS,
		}).Info("integration/mqtt: subscribing to tx topic")
		if token := i.conn.Subscribe(i.downlinkTopic, i.config.QOS, i.txPayloadHandler); token.Wait() && token.Error() != nil {
			log.WithField("topic", i.downlinkTopic).Errorf("integration/mqtt: subscribe error: %s", token.Error())
			time.Sleep(time.Second)
			continue
		}
		return
	}
}

func (i *Integration) onConnectionLost(mqttc mqtt.Client, reason error) {
	log.Errorf("integration/mqtt: mqtt connection error: %s", reason)
}

func (i *Integration) getTopic(applicationID uint64, devEUI lorawan.EUI64, eventType string) (string, error) {
	topic := bytes.NewBuffer(nil)
	err := i.eventTopicTemplate.Execute(topic, struct {
		ApplicationID uint64
		DevEUI        lorawan.EUI64
		Type          string
	}{applicationID, devEUI, eventType})
	if err != nil {
		return "", errors.Wrap(err, "execute template error")
	}

	return topic.String(), nil
}

func (i *Integration) getRetainEvents() bool {
	return i.retainEvents
}
