package external

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/http"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/influxdb"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/loracloud"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/mydevices"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/thingsboard"
	spmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"

	app "github.com/mxc-foundation/lpwan-app-server/internal/modules/application"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// ApplicationAPI exports the Application related functions.
type ApplicationAPI struct {
	st    *store.Handler
	nsCli *nscli.Client
}

// NewApplicationAPI creates a new ApplicationAPI.
func NewApplicationAPI(h *store.Handler, nsCli *nscli.Client) *ApplicationAPI {
	return &ApplicationAPI{
		st:    h,
		nsCli: nsCli,
	}
}

// Create creates the given application.
func (a *ApplicationAPI) Create(ctx context.Context, req *pb.CreateApplicationRequest) (*pb.CreateApplicationResponse, error) {
	if req.Application == nil {
		return nil, status.Errorf(codes.InvalidArgument, "application must not be nil")
	}

	if valid, err := app.NewValidator(a.st).ValidateGlobalApplicationsAccess(ctx, auth.Create, req.Application.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	spID, err := uuid.FromString(req.Application.ServiceProfileId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	sp, err := spmod.GetServiceProfile(ctx, a.st, spID, a.nsCli, true)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if sp.OrganizationID != req.Application.OrganizationId {
		return nil, status.Errorf(codes.InvalidArgument, "application and service-profile must be under the same organization")
	}

	app := Application{
		Name:                 req.Application.Name,
		Description:          req.Application.Description,
		OrganizationID:       req.Application.OrganizationId,
		ServiceProfileID:     spID,
		PayloadCodec:         req.Application.PayloadCodec,
		PayloadEncoderScript: req.Application.PayloadEncoderScript,
		PayloadDecoderScript: req.Application.PayloadDecoderScript,
	}

	err = a.st.CreateApplication(ctx, &app)
	if err != nil {
		return nil, err
	}

	return &pb.CreateApplicationResponse{
		Id: app.ID,
	}, nil
}

// Get returns the requested application.
func (a *ApplicationAPI) Get(ctx context.Context, req *pb.GetApplicationRequest) (*pb.GetApplicationResponse, error) {
	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Read, req.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	app, err := a.st.GetApplication(ctx, req.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	resp := pb.GetApplicationResponse{
		Application: &pb.Application{
			Id:                   app.ID,
			Name:                 app.Name,
			Description:          app.Description,
			OrganizationId:       app.OrganizationID,
			ServiceProfileId:     app.ServiceProfileID.String(),
			PayloadCodec:         app.PayloadCodec,
			PayloadEncoderScript: app.PayloadEncoderScript,
			PayloadDecoderScript: app.PayloadDecoderScript,
		},
	}

	return &resp, nil
}

// Update updates the given application.
func (a *ApplicationAPI) Update(ctx context.Context, req *pb.UpdateApplicationRequest) (*empty.Empty, error) {
	if req.Application == nil {
		return nil, status.Errorf(codes.InvalidArgument, "application must not be nil")
	}

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, req.Application.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	app, err := a.st.GetApplication(ctx, req.Application.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	spID, err := uuid.FromString(req.Application.ServiceProfileId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	sp, err := spmod.GetServiceProfile(ctx, a.st, spID, a.nsCli, true)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if sp.OrganizationID != app.OrganizationID {
		return nil, status.Errorf(codes.InvalidArgument, "application and service-profile must be under the same organization")
	}

	// update the fields
	app.Name = req.Application.Name
	app.Description = req.Application.Description
	app.ServiceProfileID = spID
	app.PayloadCodec = req.Application.PayloadCodec
	app.PayloadEncoderScript = req.Application.PayloadEncoderScript
	app.PayloadDecoderScript = req.Application.PayloadDecoderScript

	err = a.st.UpdateApplication(ctx, app)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Delete deletes the given application.
func (a *ApplicationAPI) Delete(ctx context.Context, req *pb.DeleteApplicationRequest) (*empty.Empty, error) {
	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Delete, req.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if err := a.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		err := handler.DeleteAllDevicesForApplicationID(ctx, req.Id)
		if err != nil {
			return err
		}

		err = handler.DeleteApplication(ctx, req.Id)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return &empty.Empty{}, nil
}

// List lists the available applications.
func (a *ApplicationAPI) List(ctx context.Context, req *pb.ListApplicationRequest) (*pb.ListApplicationResponse, error) {
	if valid, err := app.NewValidator(a.st).ValidateGlobalApplicationsAccess(ctx, auth.List, req.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	filters := ApplicationFilters{
		Search:         req.Search,
		Limit:          int(req.Limit),
		Offset:         int(req.Offset),
		OrganizationID: req.OrganizationId,
	}

	u, err := app.NewValidator(a.st).GetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	// Filter on u ID when OrganizationID is not set and the u is
	// not a global admin.
	if !u.IsGlobalAdmin && filters.OrganizationID == 0 {
		filters.UserID = u.ID
	}

	count, err := a.st.GetApplicationCount(ctx, filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	apps, err := a.st.GetApplications(ctx, filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.ListApplicationResponse{
		TotalCount: int64(count),
	}
	for _, app := range apps {
		item := pb.ApplicationListItem{
			Id:                 app.ID,
			Name:               app.Name,
			Description:        app.Description,
			OrganizationId:     app.OrganizationID,
			ServiceProfileId:   app.ServiceProfileID.String(),
			ServiceProfileName: app.ServiceProfileName,
		}

		resp.Result = append(resp.Result, &item)
	}

	return &resp, nil
}

// CreateHTTPIntegration creates an HTTP application-integration.
func (a *ApplicationAPI) CreateHTTPIntegration(ctx context.Context, in *pb.CreateHTTPIntegrationRequest) (*empty.Empty, error) {
	if in.Integration == nil {
		return nil, status.Errorf(codes.InvalidArgument, "integration must not be nil")
	}

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.Integration.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	headers := make(map[string]string)
	for _, h := range in.Integration.Headers {
		headers[h.Key] = h.Value
	}

	conf := http.Config{
		Headers:                 headers,
		DataUpURL:               in.Integration.UplinkDataUrl,
		JoinNotificationURL:     in.Integration.JoinNotificationUrl,
		ACKNotificationURL:      in.Integration.AckNotificationUrl,
		ErrorNotificationURL:    in.Integration.ErrorNotificationUrl,
		StatusNotificationURL:   in.Integration.StatusNotificationUrl,
		LocationNotificationURL: in.Integration.LocationNotificationUrl,
		/*		TxAckNotificationURL:       in.Integration.TxAckNotificationUrl,
				IntegrationNotificationURL: in.Integration.IntegrationNotificationUrl,*/
	}
	if err := conf.Validate(); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	confJSON, err := json.Marshal(conf)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	integration := Integration{
		ApplicationID: in.Integration.ApplicationId,
		Kind:          integration.HTTP,
		Settings:      confJSON,
	}
	if err = a.st.CreateIntegration(ctx, &integration); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// GetHTTPIntegration returns the HTTP application-itegration.
func (a *ApplicationAPI) GetHTTPIntegration(ctx context.Context, in *pb.GetHTTPIntegrationRequest) (*pb.GetHTTPIntegrationResponse, error) {
	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integration, err := a.st.GetIntegrationByApplicationID(ctx, in.ApplicationId, integration.HTTP)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	var conf http.Config
	if err = json.Unmarshal(integration.Settings, &conf); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	var headers []*pb.HTTPIntegrationHeader
	for k, v := range conf.Headers {
		headers = append(headers, &pb.HTTPIntegrationHeader{
			Key:   k,
			Value: v,
		})

	}

	return &pb.GetHTTPIntegrationResponse{
		Integration: &pb.HTTPIntegration{
			ApplicationId:           integration.ApplicationID,
			Headers:                 headers,
			UplinkDataUrl:           conf.DataUpURL,
			JoinNotificationUrl:     conf.JoinNotificationURL,
			AckNotificationUrl:      conf.ACKNotificationURL,
			ErrorNotificationUrl:    conf.ErrorNotificationURL,
			StatusNotificationUrl:   conf.StatusNotificationURL,
			LocationNotificationUrl: conf.LocationNotificationURL,
			/*			TxAckNotificationUrl:       conf.TxAckNotificationURL,
						IntegrationNotificationUrl: conf.IntegrationNotificationURL,*/
		},
	}, nil
}

// UpdateHTTPIntegration updates the HTTP application-integration.
func (a *ApplicationAPI) UpdateHTTPIntegration(ctx context.Context, in *pb.UpdateHTTPIntegrationRequest) (*empty.Empty, error) {
	if in.Integration == nil {
		return nil, status.Errorf(codes.InvalidArgument, "integration must not be nil")
	}

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.Integration.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integration, err := a.st.GetIntegrationByApplicationID(ctx, in.Integration.ApplicationId, integration.HTTP)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	headers := make(map[string]string)
	for _, h := range in.Integration.Headers {
		headers[h.Key] = h.Value
	}

	conf := http.Config{
		Headers:                 headers,
		DataUpURL:               in.Integration.UplinkDataUrl,
		JoinNotificationURL:     in.Integration.JoinNotificationUrl,
		ACKNotificationURL:      in.Integration.AckNotificationUrl,
		ErrorNotificationURL:    in.Integration.ErrorNotificationUrl,
		StatusNotificationURL:   in.Integration.StatusNotificationUrl,
		LocationNotificationURL: in.Integration.LocationNotificationUrl,
		/*		TxAckNotificationURL:       in.Integration.TxAckNotificationUrl,
				IntegrationNotificationURL: in.Integration.IntegrationNotificationUrl,*/
	}
	if err := conf.Validate(); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	confJSON, err := json.Marshal(conf)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	integration.Settings = confJSON

	if err = a.st.UpdateIntegration(ctx, &integration); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// DeleteHTTPIntegration deletes the application-integration of the given type.
func (a *ApplicationAPI) DeleteHTTPIntegration(ctx context.Context, in *pb.DeleteHTTPIntegrationRequest) (*empty.Empty, error) {
	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integration, err := a.st.GetIntegrationByApplicationID(ctx, in.ApplicationId, integration.HTTP)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if err = a.st.DeleteIntegration(ctx, integration.ID); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// CreateInfluxDBIntegration create an InfluxDB application-integration.
func (a *ApplicationAPI) CreateInfluxDBIntegration(ctx context.Context, in *pb.CreateInfluxDBIntegrationRequest) (*empty.Empty, error) {
	if in.Integration == nil {
		return nil, status.Errorf(codes.InvalidArgument, "integration must not be nil")
	}

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.Integration.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	conf := influxdb.Config{
		Endpoint:            in.Integration.Endpoint,
		DB:                  in.Integration.Db,
		Username:            in.Integration.Username,
		Password:            in.Integration.Password,
		RetentionPolicyName: in.Integration.RetentionPolicyName,
		Precision:           strings.ToLower(in.Integration.Precision.String()),
	}
	if err := conf.Validate(); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	confJSON, err := json.Marshal(conf)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	integration := Integration{
		ApplicationID: in.Integration.ApplicationId,
		Kind:          integration.InfluxDB,
		Settings:      confJSON,
	}
	if err := a.st.CreateIntegration(ctx, &integration); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// GetInfluxDBIntegration returns the InfluxDB application-integration.
func (a *ApplicationAPI) GetInfluxDBIntegration(ctx context.Context, in *pb.GetInfluxDBIntegrationRequest) (*pb.GetInfluxDBIntegrationResponse, error) {
	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integration, err := a.st.GetIntegrationByApplicationID(ctx, in.ApplicationId, integration.InfluxDB)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	var conf influxdb.Config
	if err = json.Unmarshal(integration.Settings, &conf); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	prec, _ := pb.InfluxDBPrecision_value[strings.ToUpper(conf.Precision)]

	return &pb.GetInfluxDBIntegrationResponse{
		Integration: &pb.InfluxDBIntegration{
			ApplicationId:       in.ApplicationId,
			Endpoint:            conf.Endpoint,
			Db:                  conf.DB,
			Username:            conf.Username,
			Password:            conf.Password,
			RetentionPolicyName: conf.RetentionPolicyName,
			Precision:           pb.InfluxDBPrecision(prec),
		},
	}, nil
}

// UpdateInfluxDBIntegration updates the InfluxDB application-integration.
func (a *ApplicationAPI) UpdateInfluxDBIntegration(ctx context.Context, in *pb.UpdateInfluxDBIntegrationRequest) (*empty.Empty, error) {
	if in.Integration == nil {
		return nil, status.Errorf(codes.InvalidArgument, "integration must not be nil")
	}

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.Integration.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integration, err := a.st.GetIntegrationByApplicationID(ctx, in.Integration.ApplicationId, integration.InfluxDB)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	conf := influxdb.Config{
		Endpoint:            in.Integration.Endpoint,
		DB:                  in.Integration.Db,
		Username:            in.Integration.Username,
		Password:            in.Integration.Password,
		RetentionPolicyName: in.Integration.RetentionPolicyName,
		Precision:           strings.ToLower(in.Integration.Precision.String()),
	}
	if err := conf.Validate(); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	confJSON, err := json.Marshal(conf)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	integration.Settings = confJSON
	if err = a.st.UpdateIntegration(ctx, &integration); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// DeleteInfluxDBIntegration deletes the InfluxDB application-integration.
func (a *ApplicationAPI) DeleteInfluxDBIntegration(ctx context.Context, in *pb.DeleteInfluxDBIntegrationRequest) (*empty.Empty, error) {
	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integration, err := a.st.GetIntegrationByApplicationID(ctx, in.ApplicationId, integration.InfluxDB)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if err = a.st.DeleteIntegration(ctx, integration.ID); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// CreateThingsBoardIntegration creates a ThingsBoard application-integration.
func (a *ApplicationAPI) CreateThingsBoardIntegration(ctx context.Context, in *pb.CreateThingsBoardIntegrationRequest) (*empty.Empty, error) {
	if in.Integration == nil {
		return nil, status.Errorf(codes.InvalidArgument, "integration must not be nil")
	}

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.Integration.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	conf := thingsboard.Config{
		Server: in.Integration.Server,
	}
	if err := conf.Validate(); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	confJSON, err := json.Marshal(conf)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	integration := Integration{
		ApplicationID: in.Integration.ApplicationId,
		Kind:          integration.ThingsBoard,
		Settings:      confJSON,
	}
	if err := a.st.CreateIntegration(ctx, &integration); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// GetThingsBoardIntegration returns the ThingsBoard application-integration.
func (a *ApplicationAPI) GetThingsBoardIntegration(ctx context.Context, in *pb.GetThingsBoardIntegrationRequest) (*pb.GetThingsBoardIntegrationResponse, error) {
	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integration, err := a.st.GetIntegrationByApplicationID(ctx, in.ApplicationId, integration.ThingsBoard)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	var conf thingsboard.Config
	if err = json.Unmarshal(integration.Settings, &conf); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.GetThingsBoardIntegrationResponse{
		Integration: &pb.ThingsBoardIntegration{
			ApplicationId: in.ApplicationId,
			Server:        conf.Server,
		},
	}, nil
}

// UpdateThingsBoardIntegration updates the ThingsBoard application-integration.
func (a *ApplicationAPI) UpdateThingsBoardIntegration(ctx context.Context, in *pb.UpdateThingsBoardIntegrationRequest) (*empty.Empty, error) {
	if in.Integration == nil {
		return nil, status.Errorf(codes.InvalidArgument, "integration must not be nil")
	}

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.Integration.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integration, err := a.st.GetIntegrationByApplicationID(ctx, in.Integration.ApplicationId, integration.ThingsBoard)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	conf := thingsboard.Config{
		Server: in.Integration.Server,
	}
	if err := conf.Validate(); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	confJSON, err := json.Marshal(conf)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	integration.Settings = confJSON
	if err = a.st.UpdateIntegration(ctx, &integration); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// DeleteThingsBoardIntegration deletes the ThingsBoard application-integration.
func (a *ApplicationAPI) DeleteThingsBoardIntegration(ctx context.Context, in *pb.DeleteThingsBoardIntegrationRequest) (*empty.Empty, error) {

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integration, err := a.st.GetIntegrationByApplicationID(ctx, in.ApplicationId, integration.ThingsBoard)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if err = a.st.DeleteIntegration(ctx, integration.ID); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// CreateMyDevicesIntegration creates a MyDevices application-integration.
func (a *ApplicationAPI) CreateMyDevicesIntegration(ctx context.Context, in *pb.CreateMyDevicesIntegrationRequest) (*empty.Empty, error) {
	if in.Integration == nil {
		return nil, status.Errorf(codes.InvalidArgument, "integration must not be nil")
	}

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.Integration.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	config := mydevices.Config{
		Endpoint: in.Integration.Endpoint,
	}
	confJSON, err := json.Marshal(config)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	integration := Integration{
		ApplicationID: in.Integration.ApplicationId,
		Kind:          integration.MyDevices,
		Settings:      confJSON,
	}
	if err := a.st.CreateIntegration(ctx, &integration); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// GetMyDevicesIntegration returns the MyDevices application-integration.
func (a *ApplicationAPI) GetMyDevicesIntegration(ctx context.Context, in *pb.GetMyDevicesIntegrationRequest) (*pb.GetMyDevicesIntegrationResponse, error) {

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integration, err := a.st.GetIntegrationByApplicationID(ctx, in.ApplicationId, integration.MyDevices)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	var conf mydevices.Config
	if err := json.Unmarshal(integration.Settings, &conf); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.GetMyDevicesIntegrationResponse{
		Integration: &pb.MyDevicesIntegration{
			ApplicationId: in.ApplicationId,
			Endpoint:      conf.Endpoint,
		},
	}, nil
}

// UpdateMyDevicesIntegration updates the MyDevices application-integration.
func (a *ApplicationAPI) UpdateMyDevicesIntegration(ctx context.Context, in *pb.UpdateMyDevicesIntegrationRequest) (*empty.Empty, error) {
	if in.Integration == nil {
		return nil, status.Errorf(codes.InvalidArgument, "integration must not be nil")
	}

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.Integration.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integration, err := a.st.GetIntegrationByApplicationID(ctx, in.Integration.ApplicationId, integration.MyDevices)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	conf := mydevices.Config{
		Endpoint: in.Integration.Endpoint,
	}

	confJSON, err := json.Marshal(conf)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	integration.Settings = confJSON
	if err = a.st.UpdateIntegration(ctx, &integration); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// DeleteMyDevicesIntegration deletes the MyDevices application-integration.
func (a *ApplicationAPI) DeleteMyDevicesIntegration(ctx context.Context, in *pb.DeleteMyDevicesIntegrationRequest) (*empty.Empty, error) {

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integration, err := a.st.GetIntegrationByApplicationID(ctx, in.ApplicationId, integration.MyDevices)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if err = a.st.DeleteIntegration(ctx, integration.ID); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// CreateLoRaCloudIntegration creates a LoRaCloud application-integration.
func (a *ApplicationAPI) CreateLoRaCloudIntegration(ctx context.Context, in *pb.CreateLoRaCloudIntegrationRequest) (*empty.Empty, error) {
	if in.Integration == nil {
		return nil, status.Errorf(codes.InvalidArgument, "integration must not be nil")
	}

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.Integration.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	config := loracloud.Config{
		Geolocation:                 in.GetIntegration().Geolocation,
		GeolocationToken:            in.GetIntegration().GeolocationToken,
		GeolocationBufferTTL:        int(in.GetIntegration().GeolocationBufferTtl),
		GeolocationMinBufferSize:    int(in.GetIntegration().GeolocationMinBufferSize),
		GeolocationTDOA:             in.GetIntegration().GeolocationTdoa,
		GeolocationRSSI:             in.GetIntegration().GeolocationRssi,
		GeolocationGNSS:             in.GetIntegration().GeolocationGnss,
		GeolocationGNSSPayloadField: in.GetIntegration().GeolocationGnssPayloadField,
		GeolocationGNSSUseRxTime:    in.GetIntegration().GeolocationGnssUseRxTime,
		GeolocationWifi:             in.GetIntegration().GeolocationWifi,
		GeolocationWifiPayloadField: in.GetIntegration().GeolocationWifiPayloadField,
		DAS:                         in.GetIntegration().Das,
		DASToken:                    in.GetIntegration().DasToken,
		DASModemPort:                uint8(in.GetIntegration().DasModemPort),
	}
	confJSON, err := json.Marshal(config)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	integration := Integration{
		ApplicationID: in.GetIntegration().ApplicationId,
		Kind:          integration.LoRaCloud,
		Settings:      confJSON,
	}
	if err := a.st.CreateIntegration(ctx, &integration); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// GetLoRaCloudIntegration returns the LoRaCloud application-integration.
func (a *ApplicationAPI) GetLoRaCloudIntegration(ctx context.Context, in *pb.GetLoRaCloudIntegrationRequest) (*pb.GetLoRaCloudIntegrationResponse, error) {

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integration, err := a.st.GetIntegrationByApplicationID(ctx, in.ApplicationId, integration.LoRaCloud)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	var conf loracloud.Config
	if err := json.Unmarshal(integration.Settings, &conf); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.GetLoRaCloudIntegrationResponse{
		Integration: &pb.LoRaCloudIntegration{
			ApplicationId:               in.ApplicationId,
			Geolocation:                 conf.Geolocation,
			GeolocationToken:            conf.GeolocationToken,
			GeolocationBufferTtl:        uint32(conf.GeolocationBufferTTL),
			GeolocationMinBufferSize:    uint32(conf.GeolocationMinBufferSize),
			GeolocationTdoa:             conf.GeolocationTDOA,
			GeolocationRssi:             conf.GeolocationRSSI,
			GeolocationGnss:             conf.GeolocationGNSS,
			GeolocationGnssPayloadField: conf.GeolocationGNSSPayloadField,
			GeolocationGnssUseRxTime:    conf.GeolocationGNSSUseRxTime,
			GeolocationWifi:             conf.GeolocationWifi,
			GeolocationWifiPayloadField: conf.GeolocationWifiPayloadField,
			Das:                         conf.DAS,
			DasToken:                    conf.DASToken,
			DasModemPort:                uint32(conf.DASModemPort),
		},
	}, nil
}

// UpdateLoRaCloudIntegration updates the LoRaCloud application-integration.
func (a *ApplicationAPI) UpdateLoRaCloudIntegration(ctx context.Context, in *pb.UpdateLoRaCloudIntegrationRequest) (*empty.Empty, error) {
	if in.Integration == nil {
		return nil, status.Errorf(codes.InvalidArgument, "integration must not be nil")
	}

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.Integration.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integration, err := a.st.GetIntegrationByApplicationID(ctx, in.GetIntegration().ApplicationId, integration.LoRaCloud)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	conf := loracloud.Config{
		Geolocation:                 in.GetIntegration().Geolocation,
		GeolocationToken:            in.GetIntegration().GeolocationToken,
		GeolocationBufferTTL:        int(in.GetIntegration().GeolocationBufferTtl),
		GeolocationMinBufferSize:    int(in.GetIntegration().GeolocationMinBufferSize),
		GeolocationTDOA:             in.GetIntegration().GeolocationTdoa,
		GeolocationRSSI:             in.GetIntegration().GeolocationRssi,
		GeolocationGNSS:             in.GetIntegration().GeolocationGnss,
		GeolocationGNSSPayloadField: in.GetIntegration().GeolocationGnssPayloadField,
		GeolocationGNSSUseRxTime:    in.GetIntegration().GeolocationGnssUseRxTime,
		GeolocationWifi:             in.GetIntegration().GeolocationWifi,
		GeolocationWifiPayloadField: in.GetIntegration().GeolocationWifiPayloadField,
		DAS:                         in.GetIntegration().Das,
		DASToken:                    in.GetIntegration().DasToken,
		DASModemPort:                uint8(in.GetIntegration().DasModemPort),
	}
	confJSON, err := json.Marshal(conf)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	integration.Settings = confJSON
	if err = a.st.UpdateIntegration(ctx, &integration); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// DeleteLoRaCloudIntegration deletes the LoRaCloud application-integration.
func (a *ApplicationAPI) DeleteLoRaCloudIntegration(ctx context.Context, in *pb.DeleteLoRaCloudIntegrationRequest) (*empty.Empty, error) {

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integration, err := a.st.GetIntegrationByApplicationID(ctx, in.ApplicationId, integration.LoRaCloud)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if err = a.st.DeleteIntegration(ctx, integration.ID); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// ListIntegrations lists all configured integrations.
func (a *ApplicationAPI) ListIntegrations(ctx context.Context, in *pb.ListIntegrationRequest) (*pb.ListIntegrationResponse, error) {

	if valid, err := app.NewValidator(a.st).ValidateApplicationAccess(ctx, auth.Update, in.ApplicationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	integrations, err := a.st.GetIntegrationsForApplicationID(ctx, in.ApplicationId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	out := pb.ListIntegrationResponse{
		TotalCount: int64(len(integrations)),
	}

	for _, intgr := range integrations {
		switch intgr.Kind {
		case integration.HTTP:
			out.Result = append(out.Result, &pb.IntegrationListItem{Kind: pb.IntegrationKind_HTTP})
		case integration.InfluxDB:
			out.Result = append(out.Result, &pb.IntegrationListItem{Kind: pb.IntegrationKind_INFLUXDB})
		case integration.ThingsBoard:
			out.Result = append(out.Result, &pb.IntegrationListItem{Kind: pb.IntegrationKind_THINGSBOARD})
		case integration.MyDevices:
			out.Result = append(out.Result, &pb.IntegrationListItem{Kind: pb.IntegrationKind_MYDEVICES})
		case integration.LoRaCloud:
			out.Result = append(out.Result, &pb.IntegrationListItem{Kind: pb.IntegrationKind_LORACLOUD})
		default:
			return nil, status.Errorf(codes.Internal, "unknown integration kind: %s", intgr.Kind)
		}
	}

	return &out, nil
}
