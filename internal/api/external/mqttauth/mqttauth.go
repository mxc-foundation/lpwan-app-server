package mqttauth

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"text/template"

	"github.com/brocaar/lorawan"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
	app "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/data"
	device "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/data"
)

// Server defines the MosquittoAuth Service Server API structure
type Server struct {
	auth auth.Authenticator
	st   Store
	jwtv *jwt.Validator

	eventTypes         []string
	eventTopicRegexp   *regexp.Regexp
	commandTopicRegexp *regexp.Regexp

	eventTopicTemplate   *template.Template
	commandTopicTemplate *template.Template
}

// Store defines set of db APIs used in mqttauth service
type Store interface {
	GetDevice(ctx context.Context, devEUI lorawan.EUI64, forUpdate bool) (device.Device, error)
	GetApplicationWithIDAndOrganizationID(ctx context.Context, id, orgID int64) (app.Application, error)
}

// NewServer returns a new MosquittoAuth Service Server
func NewServer(st Store, auth auth.Authenticator, jwtv *jwt.Validator,
	eventTopic, commandTopic *regexp.Regexp, eventTopicTemp, commandTopicTemp *template.Template) *Server {
	return &Server{
		auth: auth,
		st:   st,
		jwtv: jwtv,
		eventTypes: []string{
			UplinkEvent,
			JoinEvent,
			AckEvent,
			ErrorEvent,
			StatusEvent,
			LocationEvent,
			TxAckEvent,
		},
		eventTopicRegexp:     eventTopic,
		commandTopicRegexp:   commandTopic,
		eventTopicTemplate:   eventTopicTemp,
		commandTopicTemplate: commandTopicTemp,
	}
}

// GetJWT returns JWT for mosquitto auth with given org id
// Only accessible for authenticated supernode user
func (s *Server) GetJWT(ctx context.Context, req *pb.GetJWTRequest) (*pb.GetJWTResponse, error) {
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	if !cred.IsOrgUser {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied: %v", err)
	}

	jwToken, err := s.jwtv.SignToken(jwt.Claims{
		UserID:         cred.UserID,
		Username:       cred.Username,
		OrganizationID: req.OrganizationId,
	}, 0, []string{"mosquitto-auth"})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't create a token: %v", err)
	}

	return &pb.GetJWTResponse{JwtMqttAuth: jwToken}, nil
}

// JWTAuthentication will be called by mosquitto auth plugin JWT backend, request and response are also defined there
func (s *Server) JWTAuthentication(ctx context.Context, req *pb.JWTAuthenticationRequest) (*pb.JWTAuthenticationResponse, error) {
	_, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithAudience("mosquitto-auth"))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	return &pb.JWTAuthenticationResponse{}, nil
}

const (
	// EventTopicTemplate defines topic template which will be published by appserver, read by user
	EventTopicTemplate = "application/{{ .ApplicationID }}/device/{{ .DevEUI }}/event/{{ .Type }}"
	// CommandTopicTemplate defines topic template which will be published by user, read by appserver
	CommandTopicTemplate = "application/{{ .ApplicationID }}/device/{{ .DevEUI }}/command/{{ .Type }}"
)

const (
	// UplinkEvent contains the data and meta-data for an uplink application payload.
	UplinkEvent string = "up"
	// JoinEvent is event published when a device joins the network.
	// Please note that this is sent after the first received uplink (data) frame.
	JoinEvent string = "join"
	// AckEvent is event published on downlink frame acknowledgements.
	AckEvent string = "ack"
	// ErrorEvent is event published in case of an error related to payload scheduling or handling.
	// E.g. in case when a payload could not be scheduled as it exceeds the maximum payload-size.
	ErrorEvent string = "error"
	// StatusEvent is event for battery and margin status received from devices.
	StatusEvent string = "status"
	// LocationEvent is event to set device location
	LocationEvent string = "location"
	// TxAckEvent is event published when a downlink frame has been acknowledged by the gateway for transmission.
	TxAckEvent string = "txack"
	// IntegrationEvent is used by LoRaCloud integration, which should not be implemented here, before refactor
	// the whole integration package, leave it here temporarily
	IntegrationEvent string = "integration"
)

// CompileRegexpFromTopicTemplate creates regexp object from template, must be called in initializing of the server
func CompileRegexpFromTopicTemplate(templateName string, topicTemplate string) (*regexp.Regexp, error) {
	topicBuffer := bytes.NewBuffer(nil)
	temp, err := template.New(templateName).Parse(topicTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse template %s error: %v", topicTemplate, err)
	}
	if err := temp.Execute(topicBuffer,
		struct {
			ApplicationID string
			DevEUI        string
			Type          string
		}{`(?P<application_id>\w+)`, `(?P<dev_eui>\w+|\+)`, `(?P<type>\w+|\+)`}); err != nil {
		return nil, fmt.Errorf("create topic from temp: %v", err)
	}
	topicRegexp, err := regexp.Compile(topicBuffer.String())
	if err != nil {
		return nil, fmt.Errorf("compile regexp error: %v", err)
	}
	return topicRegexp, nil
}

// GetTopicVariables parses given topic string and extract variables from it
func GetTopicVariables(topicRegexp *regexp.Regexp, topic string) (TopicVariables, error) {
	var tv TopicVariables

	match := topicRegexp.FindStringSubmatch(topic)
	if len(match) != len(topicRegexp.SubexpNames()) {
		return tv, fmt.Errorf("topic regex match error")
	}

	result := make(map[string]string)
	for i, name := range topicRegexp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	if idStr, ok := result["application_id"]; ok {
		tv.ApplicationID = idStr
	}
	if devEUIStr, ok := result["dev_eui"]; ok {
		tv.DevEUI = devEUIStr
	}
	if typeStr, ok := result["type"]; ok {
		tv.Type = typeStr
	}
	return tv, nil
}

// TopicVariables includes all variables mentioned in different topics
type TopicVariables struct {
	// all aclReqType require
	ApplicationID string `mapstructure:"application_id"`
	// not required by subAllDevEvents
	DevEUI string `mapstructure:"dev_eui"`
	Type   string `mapstructure:"type"`
}

type aclReqType int32

const (
	subDeviceEvent  aclReqType = 0
	subAllEvents    aclReqType = 1
	subAllDevEvents aclReqType = 2
	readDeviceEvent aclReqType = 3
	pubCommand      aclReqType = 4
)

func (s *Server) verifyTopicVariables(ctx context.Context, orgID int64, variables *TopicVariables, acl aclReqType) error {
	var applicationID int64
	var devEUI lorawan.EUI64
	var err error

	applicationID, err = strconv.ParseInt(variables.ApplicationID, 10, 64)
	if err != nil {
		return fmt.Errorf("parse application id error")
	}

	// get application with id
	_, err = s.st.GetApplicationWithIDAndOrganizationID(ctx, applicationID, orgID)
	if err != nil {
		return fmt.Errorf("get application with id %d error: %v", applicationID, err)
	}

	if acl != subAllDevEvents {
		if err = devEUI.UnmarshalText([]byte(variables.DevEUI)); err != nil {
			return fmt.Errorf("parse deveui error: %v", err)
		}
		dev, err := s.st.GetDevice(ctx, devEUI, false)
		if err != nil {
			return fmt.Errorf("no such dev (%s) : %v", devEUI.String(), err)
		}
		if dev.ApplicationID != applicationID {
			return fmt.Errorf("dev (%s) is not under application %d", devEUI.String(), applicationID)
		}
	}

	return nil
}

// CheckACL will be called by mosquitto auth plugin JWT backend, request and response are also defined there
func (s *Server) CheckACL(ctx context.Context, req *pb.CheckACLRequest) (*pb.CheckACLResponse, error) {
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgIDFromToken().WithAudience("mosquitto-auth"))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	if !cred.IsDeviceAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	//acc = 1 is read, 2 is write, 3 is readwrite (not impelemented at the moment) , 4 is subscribe
	switch req.Acc {
	case 1:
		// read message from given topic
		if err := s.checkACLForRead(ctx, req.Topic, cred.OrgID); err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}
		return &pb.CheckACLResponse{}, nil
	case 4:
		// subscribe topic
		if err := s.checkACLForSubscribe(ctx, req.Topic, cred.OrgID); err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}
		return &pb.CheckACLResponse{}, nil
	case 2:
		// publish message to given topic
		if err := s.checkACLForWrite(ctx, req.Topic, cred.OrgID); err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}
		return &pb.CheckACLResponse{}, nil
	default:
		return nil, status.Errorf(codes.Unimplemented, "req.Acc is not supported: %d", req.Acc)
	}
}

func (s *Server) checkACLForSubscribe(ctx context.Context, topic string, orgID int64) error {
	var err error
	var tv TopicVariables

	// check topic application/{{ .ApplicationID }}/device/{{ .DevEUI }}/event/{{ .EventType }}
	tv, err = GetTopicVariables(s.eventTopicRegexp, topic)
	if err != nil {
		return err
	}

	if tv.DevEUI == "+" && tv.Type == "+" {
		// application/(?P<application_id>\w+)/device/+/event/+
		return s.verifyTopicVariables(ctx, orgID, &tv, subAllDevEvents)
	} else if tv.DevEUI != "" && tv.Type == "+" {
		// application/(?P<application_id>\w+)/device/(?P<dev_eui>\w+|\+)/event/+
		return s.verifyTopicVariables(ctx, orgID, &tv, subAllEvents)
	} else if tv.DevEUI != "" && tv.Type != "" {
		// application/(?P<application_id>\w+)/device/(?P<dev_eui>\w+|\+)/event/(?P<type>\w+|\+)
		return s.verifyTopicVariables(ctx, orgID, &tv, subDeviceEvent)
	}
	return fmt.Errorf("invalid topic to subscribe")
}

func (s *Server) checkACLForRead(ctx context.Context, topic string, orgID int64) error {
	var tv TopicVariables
	var err error

	tv, err = GetTopicVariables(s.eventTopicRegexp, topic)
	if err != nil {
		return err
	}

	return s.verifyTopicVariables(ctx, orgID, &tv, readDeviceEvent)
}

func (s *Server) checkACLForWrite(ctx context.Context, topic string, orgID int64) error {
	var tv TopicVariables
	var err error

	tv, err = GetTopicVariables(s.commandTopicRegexp, topic)
	if err != nil {
		return err
	}

	return s.verifyTopicVariables(ctx, orgID, &tv, pubCommand)
}

// SubsribeDeviceEvents takes device eui as request parameter,
// returns topis that can be used to subscribe to all device events or one specific event
func (s *Server) SubsribeDeviceEvents(ctx context.Context, req *pb.SubsribeDeviceEventsRequest) (*pb.SubsribeDeviceEventsResponse, error) {
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsDeviceAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	var devEUI lorawan.EUI64
	if err := devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "parse dev eui error: %v", err)
	}
	dev, err := s.st.GetDevice(ctx, devEUI, false)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "device %s not found", devEUI.String())
	}
	application, err := s.st.GetApplicationWithIDAndOrganizationID(ctx, dev.ApplicationID, req.OrganizationId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound,
			"application %d not found with given organization %d", dev.ApplicationID, req.OrganizationId)
	}

	response := pb.SubsribeDeviceEventsResponse{}
	topicBuffer := bytes.NewBuffer(nil)
	// subscribe one device one event at a time
	for _, event := range s.eventTypes {
		if err := s.eventTopicTemplate.Execute(topicBuffer, struct {
			ApplicationID string
			DevEUI        string
			Type          string
		}{fmt.Sprintf("%d", application.ID), devEUI.String(), event}); err != nil {
			return nil, status.Errorf(codes.Internal, "couldn't get evet topic for %s: %event", event, err)
		}
		response.Topic = append(response.Topic, fmt.Sprintf("Topic for subscribing to device %s on event %s: '%s'",
			devEUI.String(), event, topicBuffer.String()))
		topicBuffer.Reset()
	}

	// subscribe one device all events
	if err := s.eventTopicTemplate.Execute(topicBuffer, struct {
		ApplicationID string
		DevEUI        string
		Type          string
	}{fmt.Sprintf("%d", application.ID), devEUI.String(), "+"}); err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't get evet topic for subscribing all events: %v", err)
	}
	response.Topic = append(response.Topic, fmt.Sprintf("Topic for subscribing to device %s on all events: '%s'",
		devEUI.String(), topicBuffer.String()))
	topicBuffer.Reset()

	return &response, nil
}

// SubsribeApplicationEvents takes application id as request parameter,
// returns topics that can be used to subscribe to all devices' events under same application
func (s *Server) SubsribeApplicationEvents(ctx context.Context, req *pb.SubsribeApplicationEventsRequest) (*pb.SubsribeApplicationEventsResponse, error) {
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsDeviceAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	application, err := s.st.GetApplicationWithIDAndOrganizationID(ctx, req.ApplicationId, req.OrganizationId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound,
			"application %d not found with given organization %d", req.ApplicationId, req.OrganizationId)
	}

	var response pb.SubsribeApplicationEventsResponse
	topicBuffer := bytes.NewBuffer(nil)
	// subscribe application
	if err := s.eventTopicTemplate.Execute(topicBuffer, struct {
		ApplicationID string
		DevEUI        string
		Type          string
	}{fmt.Sprintf("%d", application.ID), "+", "+"}); err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't get evet topic for subscribing application: %v", err)
	}
	response.Topic = fmt.Sprintf("Topic for subscribing to application %d: '%s'", application.ID, topicBuffer.String())

	return &response, nil
}

// SendCommandToDevice takes device eui as request paramter,
// returns topics that can be used to send command to a specific device
func (s *Server) SendCommandToDevice(ctx context.Context, req *pb.SendCommandToDeviceRequest) (*pb.SendCommandToDeviceResponse, error) {
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsDeviceAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	var devEUI lorawan.EUI64
	if err := devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "parse dev eui error: %v", err)
	}
	dev, err := s.st.GetDevice(ctx, devEUI, false)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "device %s not found", devEUI.String())
	}
	application, err := s.st.GetApplicationWithIDAndOrganizationID(ctx, dev.ApplicationID, req.OrganizationId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound,
			"application %d not found with given organization %d", dev.ApplicationID, req.OrganizationId)
	}

	var response pb.SendCommandToDeviceResponse
	topicBuffer := bytes.NewBuffer(nil)
	// send command
	if err := s.commandTopicTemplate.Execute(topicBuffer, struct {
		ApplicationID string
		DevEUI        string
		Type          string
	}{fmt.Sprintf("%d", application.ID), devEUI.String(), "down"}); err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't get evet topic for subscribing application: %v", err)
	}
	response.Topic = fmt.Sprintf("Topic for subscribing to application %d: '%s'", application.ID, topicBuffer.String())

	return &response, nil
}
