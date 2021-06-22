package mqttauth

import (
	"context"
	"fmt"
	"regexp"
	"runtime"
	"strconv"

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

	eventTopicRegexp   *regexp.Regexp
	commandTopicRegexp *regexp.Regexp
}

// Store defines set of db APIs used in mqttauth service
type Store interface {
	GetDevice(ctx context.Context, devEUI lorawan.EUI64, forUpdate bool) (device.Device, error)
	GetApplicationWithIDAndOrganizationID(ctx context.Context, id, orgID int64) (app.Application, error)
}

// NewServer returns a new MosquittoAuth Service Server
func NewServer(st Store, auth auth.Authenticator, jwtv *jwt.Validator,
	eventTopic, commandTopic *regexp.Regexp) *Server {
	return &Server{
		auth:               auth,
		st:                 st,
		jwtv:               jwtv,
		eventTopicRegexp:   eventTopic,
		commandTopicRegexp: commandTopic,
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
	EventTopicTemplate = "application/(?P<application_id>\\w+)/device/(?P<dev_eui>\\w+|\\+)/event/(?P<type>\\w+|\\+)"
	// CommandTopicTemplate defines topic template which will be published by user, read by appserver
	CommandTopicTemplate = "application/(?P<application_id>\\w+)/device/(?P<dev_eui>\\w+)/command/(?P<type>\\w+)"
)

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
			return fmt.Errorf("dev (%s) is not under application %d", devEUI.String(), variables.ApplicationID)
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

	runtime.Breakpoint()
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
	// application/(?P<application_id>\w+)/device/(?P<dev_eui>\w+|\+)/event/(?P<type>\w+|\+)
	if tv.DevEUI != "" && tv.Type != "" {
		return s.verifyTopicVariables(ctx, orgID, &tv, subDeviceEvent)
	}
	// application/(?P<application_id>\w+)/device/(?P<dev_eui>\w+|\+)/event/+
	if tv.DevEUI != "" && tv.Type == "+" {
		return s.verifyTopicVariables(ctx, orgID, &tv, subAllEvents)
	}
	// application/(?P<application_id>\w+)/device/+/event/+
	if tv.DevEUI == "+" && tv.Type == "+" {
		return s.verifyTopicVariables(ctx, orgID, &tv, subAllDevEvents)
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
