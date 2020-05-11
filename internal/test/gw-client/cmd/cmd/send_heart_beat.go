package cmd

import (
	"context"

	api "github.com/mxc-foundation/lpwan-app-server/gw-client/api"
	"github.com/mxc-foundation/lpwan-app-server/gw-client/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var sendHeartbeatCmd = &cobra.Command{
	Use:   "heartbeat",
	Short: "Send heartbeat to supernode",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := internal.CreateClientWithCert(serverAddr, rootCAPath, clientCertPath, clientKeyPath)
		if err != nil {
			log.WithError(err).Fatal("failed to create appserver client")
		}

		resp, err := client.Heartbeat(context.Background(), &api.HeartbeatRequest{
			GatewayMac: gatewayMac,
			Model:      model,
			ConfigHash: "d91cfb21162f77e8e19a0c39ae684df0",
			OsVersion:  osVersion,
			Statistics: "",
		})
		if err != nil {
			log.WithError(err).Fatal("failed to call Connect")
		}

		log.WithFields(log.Fields{
			"new-firmware-link": resp.NewFirmwareLink,
			"config":            resp.Config,
		}).Info("Get response of calling API: Heartbeat")

		return nil
	},
}
