package gateway

import (
	"context"
	"database/sql/driver"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"

	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/provisionserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"
)

var (
	gatewayNameRegexp          = regexp.MustCompile(`^[\w-]+$`)
	serialNumberOldGWValidator = regexp.MustCompile(`^MX([A-Z1-9]){7}$`)
	serialNumberNewGWValidator = regexp.MustCompile(`^M2X([A-Z1-9]){8}$`)
)

// Value implements the driver.Valuer interface.
func (l GPSPoint) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s)", strconv.FormatFloat(l.Latitude, 'f', -1, 64), strconv.FormatFloat(l.Longitude, 'f', -1, 64)), nil
}

// Scan implements the sql.Scanner interface.
func (l *GPSPoint) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", src)
	}

	_, err := fmt.Sscanf(string(b), "(%f,%f)", &l.Latitude, &l.Longitude)
	return err
}

// Validate validates the gateway data.
func (g Gateway) Validate() error {
	if !gatewayNameRegexp.MatchString(g.Name) {
		return ErrGatewayInvalidName
	}

	if strings.HasPrefix(g.Model, "MX19") {
		if !serialNumberNewGWValidator.MatchString(g.SerialNumber) {
			return ErrGatewayInvalidSerialNumber
		}
	} else if g.Model != "" {
		if !serialNumberOldGWValidator.MatchString(g.SerialNumber) {
			return ErrGatewayInvalidSerialNumber
		}
	}

	return nil
}

var SupernodeAddr string

func UpdateFirmwareFromProvisioningServer(conf config.Config) error {
	log.WithFields(log.Fields{
		"provisioning-server": conf.ProvisionServer.ProvisionServer,
		"caCert":              conf.ProvisionServer.CACert,
		"tlsCert":             conf.ProvisionServer.TLSCert,
		"tlsKey":              conf.ProvisionServer.TLSKey,
		"schedule":            conf.ProvisionServer.UpdateSchedule,
	}).Info("Start schedule to update gateway firmware...")
	SupernodeAddr = os.Getenv("APPSERVER")
	if strings.HasPrefix(SupernodeAddr, "https://") {
		SupernodeAddr = strings.Replace(SupernodeAddr, "https://", "", -1)
	}
	if strings.HasPrefix(SupernodeAddr, "http://") {
		SupernodeAddr = strings.Replace(SupernodeAddr, "http://", "", -1)
	}
	if strings.HasSuffix(SupernodeAddr, ":8080") {
		SupernodeAddr = strings.Replace(SupernodeAddr, ":8080", "", -1)
	}
	SupernodeAddr = strings.Replace(SupernodeAddr, "/", "", -1)

	var bindPortOldGateway string
	var bindPortNewGateway string

	if strArray := strings.Split(conf.ApplicationServer.APIForGateway.OldGateway.Bind, ":"); len(strArray) != 2 {
		return errors.New(fmt.Sprintf("Invalid API Bind settings for OldGateway: %s", conf.ApplicationServer.APIForGateway.OldGateway.Bind))
	} else {
		bindPortOldGateway = strArray[1]
	}

	if strArray := strings.Split(conf.ApplicationServer.APIForGateway.NewGateway.Bind, ":"); len(strArray) != 2 {
		return errors.New(fmt.Sprintf("Invalid API Bind settings for NewGateway: %s", conf.ApplicationServer.APIForGateway.NewGateway.Bind))
	} else {
		bindPortNewGateway = strArray[1]
	}

	c := cron.New()
	err := c.AddFunc(conf.ProvisionServer.UpdateSchedule, func() {
		log.Info("Check firmware update...")
		gwFwList, err := GetGatewayFirmwareList(DB())
		if err != nil {
			log.WithError(err).Errorf("Failed to get gateway firmware list.")
			return
		}

		// send update
		psClient, err := provisionserver.CreateClientWithCert(conf.ProvisionServer.ProvisionServer,
			conf.ProvisionServer.CACert,
			conf.ProvisionServer.TLSCert,
			conf.ProvisionServer.TLSKey)
		if err != nil {
			log.WithError(err).Errorf("Create Provisioning server client error")
			return
		}

		for _, v := range gwFwList {
			res, err := psClient.GetUpdate(context.Background(), &psPb.GetUpdateRequest{
				Model:          v.Model,
				SuperNodeAddr:  SupernodeAddr,
				PortOldGateway: bindPortOldGateway,
				PortNewGateway: bindPortNewGateway,
			})
			if err != nil {
				log.WithError(err).Errorf("Failed to get update for gateway model: %s", v.Model)
				continue
			}

			var md5sum types.MD5SUM
			if err := md5sum.UnmarshalText([]byte(res.FirmwareHash)); err != nil {
				log.WithError(err).Errorf("Failed to unmarshal firmware hash: %s", res.FirmwareHash)
				continue
			}

			gatewayFw := GatewayFirmware{
				Model:        v.Model,
				ResourceLink: res.ResourceLink,
				FirmwareHash: md5sum,
			}

			model, err := UpdateGatewayFirmware(DB(), &gatewayFw)
			if model == "" {
				log.Warnf("No row updated for gateway_firmware at model=%s", v.Model)
			}

		}
	})
	if err != nil {
		log.Fatalf("Failed to set update schedule when set up provisioning server config: %s", err.Error())
	}

	go c.Start()

	return nil
}
