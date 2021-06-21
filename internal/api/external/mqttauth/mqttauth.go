package mqttauth

import (
	"context"
	"fmt"
	"github.com/brocaar/lorawan"
	"regexp"
	"strconv"

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

	eventTopicRegexp     *regexp.Regexp
	commandTopicRegexp   *regexp.Regexp
	allEventsTopicRegexp *regexp.Regexp
}

// Store defines set of db APIs used in mqttauth service
type Store interface {
	GetDevice(ctx context.Context, devEUI lorawan.EUI64, forUpdate bool) (device.Device, error)
	GetApplicationWithIDAndOrganizationID(ctx context.Context, id, orgID int64) (app.Application, error)
}

// NewServer returns a new MosquittoAuth Service Server
func NewServer(st Store, auth auth.Authenticator, jwtv *jwt.Validator,
	eventTopic, commandTopic, allEventsTopic *regexp.Regexp) *Server {
	return &Server{
		auth:                 auth,
		st:                   st,
		jwtv:                 jwtv,
		eventTopicRegexp:     eventTopic,
		commandTopicRegexp:   commandTopic,
		allEventsTopicRegexp: allEventsTopic,
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
	EventTopicTemplate = "application/(?P<application_id>\\w+)/device/(?P<dev_eui>\\w+)/event/(?P<type>\\w)"
	// AllEventsTopicTemplate defines topic that can be subscribed by user
	AllEventsTopicTemplate = "application/(?P<application_id>\\w+)/device/(?P<dev_eui>\\w+)/event/#"

	// CommandTopicTemplate defines topic template which will be published by user, read by appserver
	CommandTopicTemplate = "application/(?P<application_id>\\w+)/device/(?P<dev_eui>\\w+)/command/(?P<type>\\w)"
)

// GetTopicVariables parses given topic string and extract variables from it
func GetTopicVariables(topicRegexp *regexp.Regexp, topic string) (int64, lorawan.EUI64, error) {
	var applicationID int64
	var devEUI lorawan.EUI64
	var err error

	match := topicRegexp.FindStringSubmatch(topic)
	if len(match) != len(topicRegexp.SubexpNames()) {
		return applicationID, devEUI, fmt.Errorf("topic regex match error")
	}

	result := make(map[string]string)
	for i, name := range topicRegexp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	if idStr, ok := result["application_id"]; ok {
		applicationID, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return applicationID, devEUI, fmt.Errorf("parse application id error: %v", err)
		}
	} else {
		return applicationID, devEUI, fmt.Errorf("topic regexp does not contain application id")
	}

	if devEUIStr, ok := result["dev_eui"]; ok {
		if err = devEUI.UnmarshalText([]byte(devEUIStr)); err != nil {
			return applicationID, devEUI, fmt.Errorf("parse deveui error: %v", err)
		}
	}

	return applicationID, devEUI, nil
}

func (s *Server) verifyTopicVariables(ctx context.Context, applicationID, organizationID int64, devEUI lorawan.EUI64) error {
	// get application with id
	_, err := s.st.GetApplicationWithIDAndOrganizationID(ctx, applicationID, organizationID)
	if err != nil {
		return fmt.Errorf("failed to get application with id %d: %v", applicationID, err)
	}

	device, err := s.st.GetDevice(ctx, devEUI, false)
	if err != nil {
		return fmt.Errorf("no such device (%s) : %v", devEUI.String(), err)
	}

	if device.ApplicationID != applicationID {
		return fmt.Errorf("device (%s) is not under application %d", devEUI.String(), applicationID)
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
		return &pb.CheckACLResponse{}, s.checkACLForRead(ctx, req.Topic, cred.OrgID, cred.Username)
	case 4:
		// subscribe topic
		return &pb.CheckACLResponse{}, s.checkACLForSubscribe(ctx, req.Topic, cred.OrgID, cred.Username)
	case 2:
		// publish message to given topic
		return &pb.CheckACLResponse{}, s.checkACLForWrite(ctx, req.Topic, cred.OrgID, cred.Username)
	default:
		return nil, status.Errorf(codes.Unimplemented, "req.Acc is not supported: %d", req.Acc)
	}
}

func (s *Server) checkACLForSubscribe(ctx context.Context, topic string, orgID int64, email string) error {
	var applicationID int64
	var devEUI lorawan.EUI64
	var err error

	// check topic application/{{ .ApplicationID }}/device/{{ .DevEUI }}/event/{{ .EventType }}
	applicationID, devEUI, err = GetTopicVariables(s.eventTopicRegexp, topic)
	if err == nil {
		return s.verifyTopicVariables(ctx, applicationID, orgID, devEUI)
	}
	// check topic application/{{ .ApplicationID }}/device/{{ .DevEUI }}/event/#
	applicationID, devEUI, err = GetTopicVariables(s.allEventsTopicRegexp, topic)
	if err == nil {
		return s.verifyTopicVariables(ctx, applicationID, orgID, devEUI)
	}

	return fmt.Errorf("user %s subscribing to topic %s rejected", email, topic)
}

func (s *Server) checkACLForRead(ctx context.Context, topic string, orgID int64, email string) error {
	var applicationID int64
	var devEUI lorawan.EUI64
	var err error

	applicationID, devEUI, err = GetTopicVariables(s.eventTopicRegexp, topic)
	if err != nil {
		return fmt.Errorf("user %s reading topic %s rejected", email, topic)
	}

	return s.verifyTopicVariables(ctx, applicationID, orgID, devEUI)
}

func (s *Server) checkACLForWrite(ctx context.Context, topic string, orgID int64, email string) error {
	var applicationID int64
	var devEUI lorawan.EUI64
	var err error

	applicationID, devEUI, err = GetTopicVariables(s.commandTopicRegexp, topic)
	if err != nil {
		return fmt.Errorf("user %s writing topic %s rejected", email, topic)
	}

	return s.verifyTopicVariables(ctx, applicationID, orgID, devEUI)
}
