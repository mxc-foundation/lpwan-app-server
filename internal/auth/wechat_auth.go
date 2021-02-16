package auth

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/httpcli"
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

// WeChatAuth defines necessary variables for authenticating a wechat user
type WeChatAuth struct {
	AccessToken string `json:"access_token"`
	OpenID      string `json:"openid"`
}

// GetWeChatUserInfoResponse represents struct of get wechat user info response
type GetWeChatUserInfoResponse struct {
	OpenID     string `json:"openid"`
	HeadImgURL string `json:"headimgurl"`
	UnionID    string `json:"unionid"`
	NickName   string `json:"nickname"`
}

// GetAccessTokenFromCode sends http request and return response for getting access token with appid, secret and code
func GetAccessTokenFromCode(ctx context.Context, code, appID, secret string, response *GetAccessTokenResponse) error {
	if code == "" || appID == "" || secret == "" {
		return fmt.Errorf("cannot get access_token: invalid argument")
	}

	// get access_token
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		appID, secret, code)
	if err := httpcli.GetResponse(url, response, true); err != nil {
		return err
	}

	log.Debugf("GetAccessTokenFromCode response: {'access_token': %s, 'openid': %s, 'unionid': %s, 'scope': %s}",
		response.AccessToken, response.OpenID, response.UnionID, response.Scope)

	return nil
}

// GetWeChatUserInfoFromAccessToken sends http request and return response of getting wechat user info with access token and openid
func GetWeChatUserInfoFromAccessToken(ctx context.Context, accessToken, openID string, response *GetWeChatUserInfoResponse) error {
	// get user info
	if accessToken == "" || openID == "" {
		return fmt.Errorf("cannot get wechat user info: invalid argument")
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", accessToken, openID)
	if err := httpcli.GetResponse(url, response, false); err != nil {
		return err
	}

	if response.UnionID == "" {
		return fmt.Errorf("unionid is required, cannot be empty string")
	}

	log.Debugf("GetWeChatUserInfoFromAccessToken response: {'unionid': %s, 'nickname': %s, 'openid': %s}",
		response.UnionID, response.NickName, response.OpenID)

	return nil
}
