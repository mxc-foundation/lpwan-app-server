package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
)

// WeChatAuthentication defines configuration for authorizing wechat users
type WeChatAuthentication struct {
	AppID  string `mapstructure:"app_id"`
	Secret string `mapstructure:"secret"`
}

// GetAccessTokenResponse represents struct of get wechat access_token response
type GetAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
}

type GetWeChatUserInfoResponse struct {
	OpenID     string `json:"openid"`
	HeadImgURL string `json:"headimgurl"`
	UnionID    string `json:"unionid"`
}

func getHTTPResponse(url string, dest interface{}) error {
	// #nosec
	resp, err := http.Get(url)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return status.Errorf(codes.FailedPrecondition, resp.Status)
	}

	// disallow unknow fileds to filter out error messages from wechat server when no err is returned
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dest); err != nil {
		return status.Errorf(codes.FailedPrecondition, err.Error())
	}

	return nil
}

// AuthenticateWeChatUser interacts with wechat open platform to authenticate wechat user
// then check binding status of this wechat user
func (a *Server) AuthenticateWeChatUser(ctx context.Context, req *pb.AuthenticateWeChatUserRequest) (*pb.AuthenticateWeChatUserResponse, error) {
	// get access_token
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		a.config.WeChatLogin.AppID, a.config.WeChatLogin.Secret, req.Code)
	body := GetAccessTokenResponse{}

	if err := getHTTPResponse(url, &body); err != nil {
		return nil, err
	}

	if body.UnionID == "" {
		// get user info
		if body.AccessToken == "" || body.OpenID == "" {
			return nil, status.Errorf(codes.DataLoss, "")
		}

		url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", body.AccessToken, body.OpenID)
		user := GetWeChatUserInfoResponse{}

		if err := getHTTPResponse(url, &user); err != nil {
			return nil, err
		}

		body.UnionID = user.UnionID
	}

	// check whether wechat user has already bound to existing account

	return &pb.AuthenticateWeChatUserResponse{Jwt: "", BindingIsRequired: false}, nil
}

// BindExternalUser binds external user id to supernode user
func (a *Server) BindExternalUser(ctx context.Context, req *pb.BindExternalUserRequest) (*pb.BindExternalUserResponse, error) {

	return &pb.BindExternalUserResponse{}, nil
}

// RegisterExternalUser creates new supernode account then bind it with external user id
func (a *Server) RegisterExternalUser(ctx context.Context, req *pb.RegisterExternalUserRequest) (*pb.RegisterExternalUserResponse, error) {

	return &pb.RegisterExternalUserResponse{}, nil
}

// GetExternalUserIDFromUserID gets external user id by supernode user id
func (a *Server) GetExternalUserIDFromUserID(ctx context.Context, req *pb.GetExternalUserIDFromUserIDRequest) (*pb.GetExternalUserIDFromUserIDResponse, error) {

	return &pb.GetExternalUserIDFromUserIDResponse{}, nil
}

// UnbindExternalUser unbinds external user and supernode user account
func (a *Server) UnbindExternalUser(ctx context.Context, req *pb.UnbindExternalUserRequest) (*pb.UnbindExternalUserResponse, error) {

	return &pb.UnbindExternalUserResponse{}, nil
}
