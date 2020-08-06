package m2m

import (
	"context"
	"fmt"
	"net"

	"github.com/mxc-foundation/lpwan-app-server/internal/email"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/brocaar/lorawan"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/tls"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

var serviceName = "m2m server"

// Setup :
func Setup(conf config.Config) error {
	log.Info("Set up API for m2m server")

	if err := listenWithCredentials(conf.ApplicationServer.APIForM2M.Bind,
		conf.ApplicationServer.APIForM2M.CACert,
		conf.ApplicationServer.APIForM2M.TLSCert,
		conf.ApplicationServer.APIForM2M.TLSKey); err != nil {
		return err
	}

	return nil
}

func listenWithCredentials(bind, caCert, tlsCert, tlsKey string) error {
	log.WithFields(log.Fields{
		"bind":     bind,
		"ca-cert":  caCert,
		"tls-cert": tlsCert,
		"tls-key":  tlsKey,
	}).Info("listen With Credentials")

	gs, err := tls.NewServerWithTLSCredentials(serviceName, caCert, tlsCert, tlsKey)
	if err != nil {
		return errors.Wrap(err, "listenWithCredentials: get new server error")
	}

	m2mAPI := NewM2MAPI()
	pb.RegisterAppServerServiceServer(gs, m2mAPI)

	ln, err := net.Listen("tcp", bind)
	if err != nil {
		return errors.Wrap(err, "listenWithCredentials: start api listener error")
	}
	go gs.Serve(ln)

	return nil
}

// API exports the API related functions.
type API struct {
}

// NewM2MAPI creates new API
func NewM2MAPI() *API {
	return &API{}
}

// GetDeviceDevEuiList defines the response of the Device DevEui list
func (a *API) GetDeviceDevEuiList(ctx context.Context, req *empty.Empty) (*pb.GetDeviceDevEuiListResponse, error) {
	devEuiList, err := storage.GetAllDeviceEuis(ctx, storage.DB())
	if err != nil {
		return &pb.GetDeviceDevEuiListResponse{}, status.Errorf(codes.DataLoss, err.Error())
	}

	return &pb.GetDeviceDevEuiListResponse{DevEui: devEuiList}, nil
}

// GetGatewayMacList defines the response of the Gateway MAC list
func (a *API) GetGatewayMacList(ctx context.Context, req *empty.Empty) (*pb.GetGatewayMacListResponse, error) {
	gwMacList, err := storage.GetAllGatewayMacList(ctx, storage.DB())
	if err != nil {
		return &pb.GetGatewayMacListResponse{}, status.Errorf(codes.DataLoss, err.Error())
	}

	return &pb.GetGatewayMacListResponse{GatewayMac: gwMacList}, nil
}

// GetDeviceByDevEui defines the request and response of the Device DevEui
func (a *API) GetDeviceByDevEui(ctx context.Context, req *pb.GetDeviceByDevEuiRequest) (*pb.GetDeviceByDevEuiResponse, error) {
	var devEui lorawan.EUI64
	resp := pb.GetDeviceByDevEuiResponse{DevProfile: &pb.AppServerDeviceProfile{}}

	if err := devEui.UnmarshalText([]byte(req.DevEui)); err != nil {
		return &resp, status.Errorf(codes.InvalidArgument, err.Error())
	}

	device, err := storage.GetDevice(ctx, storage.DB(), devEui, false, true)
	if err == storage.ErrDoesNotExist {
		return &resp, nil
	} else if err != nil {
		return &resp, status.Errorf(codes.Unknown, err.Error())
	}

	application, err := storage.GetApplication(ctx, storage.DB(), device.ApplicationID)
	if err != nil {
		return &resp, status.Errorf(codes.Unknown, err.Error())
	}

	resp.OrgId = application.OrganizationID
	resp.DevProfile.DevEui = req.DevEui
	resp.DevProfile.Name = device.Name
	resp.DevProfile.ApplicationId = device.ApplicationID
	resp.DevProfile.CreatedAt, _ = ptypes.TimestampProto(device.CreatedAt)

	return &resp, nil
}

// GetGatewayByMac defines the request and response to the the gateway by MAC
func (a *API) GetGatewayByMac(ctx context.Context, req *pb.GetGatewayByMacRequest) (*pb.GetGatewayByMacResponse, error) {
	var mac lorawan.EUI64
	resp := pb.GetGatewayByMacResponse{GwProfile: &pb.AppServerGatewayProfile{}}

	if err := mac.UnmarshalText([]byte(req.Mac)); err != nil {
		return &resp, status.Errorf(codes.InvalidArgument, err.Error())
	}

	gateway, err := storage.GetGateway(ctx, storage.DB(), mac, false)
	if err == storage.ErrDoesNotExist {
		return &resp, nil
	} else if err != nil {
		return &resp, status.Errorf(codes.InvalidArgument, err.Error())
	}

	resp.OrgId = gateway.OrganizationID
	resp.GwProfile.OrgId = gateway.OrganizationID
	resp.GwProfile.Mac = req.Mac
	resp.GwProfile.Name = gateway.Name
	resp.GwProfile.Description = gateway.Description
	resp.GwProfile.CreatedAt, _ = ptypes.TimestampProto(gateway.CreatedAt)

	return &resp, nil
}

// SendStakeIncomeNotification is called to send email to user when stake revenue is applied
func (a *API) SendStakeIncomeNotification(ctx context.Context, req *pb.SendStakeIncomeNotificationRequest) (*pb.SendStakeIncomeNotificationResponse, error) {
	resp := pb.SendStakeIncomeNotificationResponse{}
	// get user id from organization id
	users, err := storage.GetOrganizationUsers(ctx, storage.DB(), req.OrganizationId, 999, 0)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to get users for organization")
	}

	for _, v := range users {
		amountMap := make(map[string]string)
		amountMap[email.StakeIncomeAmount] = req.StakeIncomeAmount
		amountMap[email.StakeAmount] = req.StakeAmount
		amountMap[email.StakeIncomeInterest] = req.StakeIncomeInterest

		itemIDMap := make(map[string]string)
		itemIDMap[email.UserID] = fmt.Sprintf("%d", v.UserID)
		itemIDMap[email.StakeID] = req.StakeId
		itemIDMap[email.StakeRevenueID] = req.StakeRevenueId

		dateMap := make(map[string]string)
		dateMap[email.StakeStartDate] = req.StakeStartDate
		dateMap[email.StakeRevenueDate] = req.StakeRevenueDate

		_ = email.SendInvite(v.Username, email.Param{
			Amount: amountMap,
			ItemID: itemIDMap,
			Date:   dateMap,
		}, email.EmailLanguage(config.C.General.DefaultLanguage), email.StakingIncome)

	}

	return &resp, nil
}
