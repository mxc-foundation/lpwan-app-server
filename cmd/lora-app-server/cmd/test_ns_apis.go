package cmd

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	pb "github.com/brocaar/chirpstack-api/go/v3/ns"
	ns "github.com/mxc-foundation/lpwan-app-server/internal/clients/networkserver"
)

var testNsAPICmd = &cobra.Command{
	Use:   "test-ns-api",
	Short: "Batch test all ns apis",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ns.Setup(); err != nil {
			panic(err)
		}

		s := &ns.NSStruct{
			Server:  "network-server:8000",
			CACert:  "",
			TLSCert: "",
			TLSKey:  "",
		}
		nsCli, err := s.GetNetworkServiceClient()
		if err != nil {
			panic(err)
		}

		ctx := context.TODO()

		// CreateServiceProfile creates the given service-profile.
		_, err = nsCli.CreateServiceProfile(ctx, &pb.CreateServiceProfileRequest{})
		log.Warnf("CreateServiceProfile: ", err)
		// GetServiceProfile returns the service-profile matching the given id.
		_, err = nsCli.GetServiceProfile(ctx, &pb.GetServiceProfileRequest{})
		log.Warnf("GetServiceProfile: ", err)
		// UpdateServiceProfile updates the given service-profile.
		_, err = nsCli.UpdateServiceProfile(ctx, &pb.UpdateServiceProfileRequest{})
		log.Warnf("UpdateServiceProfile: ", err)
		// DeleteServiceProfile deletes the service-profile matching the given id.
		_, err = nsCli.DeleteServiceProfile(ctx, &pb.DeleteServiceProfileRequest{})
		log.Warnf("DeleteServiceProfile: ", err)
		// CreateRoutingProfile creates the given routing-profile.
		_, err = nsCli.CreateRoutingProfile(ctx, &pb.CreateRoutingProfileRequest{})
		log.Warnf("CreateRoutingProfile: ", err)
		// GetRoutingProfile returns the routing-profile matching the given id.
		_, err = nsCli.GetRoutingProfile(ctx, &pb.GetRoutingProfileRequest{})
		log.Warnf("GetRoutingProfile: ", err)
		// UpdateRoutingProfile updates the given routing-profile.
		_, err = nsCli.UpdateRoutingProfile(ctx, &pb.UpdateRoutingProfileRequest{})
		log.Warnf("UpdateRoutingProfile: ", err)
		// DeleteRoutingProfile deletes the routing-profile matching the given id.
		_, err = nsCli.DeleteRoutingProfile(ctx, &pb.DeleteRoutingProfileRequest{})
		log.Warnf("DeleteRoutingProfile: ", err)
		// CreateDeviceProfile creates the given device-profile.
		_, err = nsCli.CreateDeviceProfile(ctx, &pb.CreateDeviceProfileRequest{})
		log.Warnf("CreateDeviceProfile: ", err)
		// GetDeviceProfile returns the device-profile matching the given id.
		_, err = nsCli.GetDeviceProfile(ctx, &pb.GetDeviceProfileRequest{})
		log.Warnf("GetDeviceProfile: ", err)
		// UpdateDeviceProfile updates the given device-profile.
		_, err = nsCli.UpdateDeviceProfile(ctx, &pb.UpdateDeviceProfileRequest{})
		log.Warnf("UpdateDeviceProfile: ", err)
		// DeleteDeviceProfile deletes the device-profile matching the given id.
		_, err = nsCli.DeleteDeviceProfile(ctx, &pb.DeleteDeviceProfileRequest{})
		log.Warnf("DeleteDeviceProfile: ", err)
		// CreateDevice creates the given device.
		_, err = nsCli.CreateDevice(ctx, &pb.CreateDeviceRequest{})
		log.Warnf("CreateDevice: ", err)
		// GetDevice returns the device matching the given DevEUI.
		_, err = nsCli.GetDevice(ctx, &pb.GetDeviceRequest{})
		log.Warnf("GetDevice: ", err)
		// UpdateDevice updates the given device.
		_, err = nsCli.UpdateDevice(ctx, &pb.UpdateDeviceRequest{})
		log.Warnf("UpdateDevice: ", err)
		// DeleteDevice deletes the device matching the given DevEUI.
		_, err = nsCli.DeleteDevice(ctx, &pb.DeleteDeviceRequest{})
		log.Warnf("DeleteDevice: ", err)
		// ActivateDevice activates a device (ABP{})
		_, err = nsCli.ActivateDevice(ctx, &pb.ActivateDeviceRequest{})
		log.Warnf("ActivateDevice: ", err)
		// DeactivateDevice de-activates a device.
		_, err = nsCli.DeactivateDevice(ctx, &pb.DeactivateDeviceRequest{})
		log.Warnf("DeactivateDevice: ", err)
		// GetDeviceActivation returns the device activation details.
		_, err = nsCli.GetDeviceActivation(ctx, &pb.GetDeviceActivationRequest{})
		log.Warnf("GetDeviceActivation: ", err)
		// CreateDeviceQueueItem creates the given device-queue item.
		_, err = nsCli.CreateDeviceQueueItem(ctx, &pb.CreateDeviceQueueItemRequest{})
		log.Warnf("CreateDeviceQueueItem: ", err)
		// FlushDeviceQueueForDevEUI flushes the device-queue for the given DevEUI.
		_, err = nsCli.FlushDeviceQueueForDevEUI(ctx, &pb.FlushDeviceQueueForDevEUIRequest{})
		log.Warnf("FlushDeviceQueueForDevEUI: ", err)
		// GetDeviceQueueItemsForDevEUI returns all device-queue items for the given DevEUI.
		_, err = nsCli.GetDeviceQueueItemsForDevEUI(ctx, &pb.GetDeviceQueueItemsForDevEUIRequest{})
		log.Warnf("GetDeviceQueueItemsForDevEUI: ", err)
		// GetNextDownlinkFCntForDevEUI returns the next FCnt that must be used.
		// This also takes device-queue items for the given DevEUI into consideration.
		_, err = nsCli.GetNextDownlinkFCntForDevEUI(ctx, &pb.GetNextDownlinkFCntForDevEUIRequest{})
		log.Warnf("GetNextDownlinkFCntForDevEUI: ", err)
		// GetRandomDevAddr returns a random DevAddr taking the NwkID prefix into account.
		_, err = nsCli.GetRandomDevAddr(ctx, &empty.Empty{})
		log.Warnf("GetRandomDevAddr: ", err)
		// CreateMACCommandQueueItem adds the downlink mac-command to the queue.
		_, err = nsCli.CreateMACCommandQueueItem(ctx, &pb.CreateMACCommandQueueItemRequest{})
		log.Warnf("CreateMACCommandQueueItem: ", err)
		// SendProprietaryPayload send a payload using the 'Proprietary' LoRaWAN message-type.
		_, err = nsCli.SendProprietaryPayload(ctx, &pb.SendProprietaryPayloadRequest{})
		log.Warnf("SendProprietaryPayload: ", err)
		// CreateGateway creates the given gateway.
		_, err = nsCli.CreateGateway(ctx, &pb.CreateGatewayRequest{})
		log.Warnf("CreateGateway: ", err)
		// GetGateway returns data for a particular gateway.
		_, err = nsCli.GetGateway(ctx, &pb.GetGatewayRequest{})
		log.Warnf("GetGateway: ", err)
		// UpdateGateway updates an existing gateway.
		_, err = nsCli.UpdateGateway(ctx, &pb.UpdateGatewayRequest{})
		log.Warnf("UpdateGateway: ", err)
		// DeleteGateway deletes a gateway.
		_, err = nsCli.DeleteGateway(ctx, &pb.DeleteGatewayRequest{})
		log.Warnf("DeleteGateway: ", err)
		// GenerateGatewayClientCertificate returns TLS certificate gateway authentication / authorization.
		// This endpoint can ony be used when ChirpStack Network Server is configured with a gateway
		// CA certificate and key, which is used for signing the TLS certificate. The returned TLS
		// certificate will have the Gateway ID as Common Name.
		_, err = nsCli.GenerateGatewayClientCertificate(ctx, &pb.GenerateGatewayClientCertificateRequest{})
		log.Warnf("GenerateGatewayClientCertificate: ", err)
		// CreateGatewayProfile creates the given gateway-profile.
		_, err = nsCli.CreateGatewayProfile(ctx, &pb.CreateGatewayProfileRequest{})
		log.Warnf("CreateGatewayProfile: ", err)
		// GetGatewayProfile returns the gateway-profile given an id.
		_, err = nsCli.GetGatewayProfile(ctx, &pb.GetGatewayProfileRequest{})
		log.Warnf("GetGatewayProfile: ", err)
		// UpdateGatewayProfile updates the given gateway-profile.
		_, err = nsCli.UpdateGatewayProfile(ctx, &pb.UpdateGatewayProfileRequest{})
		log.Warnf("UpdateGatewayProfile: ", err)
		// DeleteGatewayProfile deletes the gateway-profile matching a given id.
		_, err = nsCli.DeleteGatewayProfile(ctx, &pb.DeleteGatewayProfileRequest{})
		log.Warnf("DeleteGatewayProfile: ", err)
		// GetGatewayStats returns stats of an existing gateway.
		// Deprecated (stats are forwarded to Application Server API{})
		_, err = nsCli.GetGatewayStats(ctx, &pb.GetGatewayStatsRequest{})
		log.Warnf("GetGatewayStats: ", err)
		// StreamFrameLogsForGateway returns a stream of frames seen by the given gateway.
		_, err = nsCli.StreamFrameLogsForGateway(ctx, &pb.StreamFrameLogsForGatewayRequest{})
		log.Warnf("StreamFrameLogsForGateway: ", err)
		// StreamFrameLogsForDevice returns a stream of frames seen by the given device.
		_, err = nsCli.StreamFrameLogsForDevice(ctx, &pb.StreamFrameLogsForDeviceRequest{})
		log.Warnf("StreamFrameLogsForDevice: ", err)
		// CreateMulticastGroup creates the given multicast-group.
		_, err = nsCli.CreateMulticastGroup(ctx, &pb.CreateMulticastGroupRequest{})
		log.Warnf("CreateMulticastGroup: ", err)
		// GetMulticastGroup returns the multicast-group given an id.
		_, err = nsCli.GetMulticastGroup(ctx, &pb.GetMulticastGroupRequest{})
		log.Warnf("GetMulticastGroup: ", err)
		// UpdateMulticastGroup updates the given multicast-group.
		_, err = nsCli.UpdateMulticastGroup(ctx, &pb.UpdateMulticastGroupRequest{})
		log.Warnf("UpdateMulticastGroup: ", err)
		// DeleteMulticastGroup deletes a multicast-group given an id.
		_, err = nsCli.DeleteMulticastGroup(ctx, &pb.DeleteMulticastGroupRequest{})
		log.Warnf("DeleteMulticastGroup: ", err)
		// AddDeviceToMulticastGroup adds the given device to the given multicast-group.
		_, err = nsCli.AddDeviceToMulticastGroup(ctx, &pb.AddDeviceToMulticastGroupRequest{})
		log.Warnf("AddDeviceToMulticastGroup: ", err)
		// RemoveDeviceFromMulticastGroup removes the given device from the given multicast-group.
		_, err = nsCli.RemoveDeviceFromMulticastGroup(ctx, &pb.RemoveDeviceFromMulticastGroupRequest{})
		log.Warnf("RemoveDeviceFromMulticastGroup: ", err)
		// EnqueueMulticastQueueItem enqueues the given multicast queue-item and
		// increments the frame-counter after enqueueing.
		_, err = nsCli.EnqueueMulticastQueueItem(ctx, &pb.EnqueueMulticastQueueItemRequest{})
		log.Warnf("EnqueueMulticastQueueItem: ", err)
		// FlushMulticastQueueForMulticastGroup flushes the multicast device-queue given a multicast-group id.
		_, err = nsCli.FlushMulticastQueueForMulticastGroup(ctx, &pb.FlushMulticastQueueForMulticastGroupRequest{})
		log.Warnf("FlushMulticastQueueForMulticastGroup: ", err)
		// GetMulticastQueueItemsForMulticastGroup returns the queue-items given a multicast-group id.
		_, err = nsCli.GetMulticastQueueItemsForMulticastGroup(ctx, &pb.GetMulticastQueueItemsForMulticastGroupRequest{})
		log.Warnf("GetMulticastQueueItemsForMulticastGroup: ", err)
		// GetVersion returns the ChirpStack Network Server version.
		_, err = nsCli.GetVersion(ctx, &empty.Empty{})
		log.Warnf("GetVersion: ", err)
	},
}
