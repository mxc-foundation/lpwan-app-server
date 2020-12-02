package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
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

type WeChatAuth struct {
	AccessToken string `json:"access_token"`
	OpenID      string `json:"openid"`
}

type GetWeChatUserInfoResponse struct {
	OpenID     string `json:"openid"`
	HeadImgURL string `json:"headimgurl"`
	UnionID    string `json:"unionid"`
	NickName   string `json:"nickname"`
}

func GetHTTPResponse(url string, dest interface{}, disallowUnknowFields bool) error {
	log.WithFields(log.Fields{
		"url":                  url,
		"disallowUnkownFields": disallowUnknowFields,
	}).Debug("auth.GetHTTPResponse is called")

	// #nosec
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("invalid url %s", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(resp.Status)
	}

	// disallow unknow fileds to filter out error messages from wechat server when no err is returned
	decoder := json.NewDecoder(resp.Body)
	if disallowUnknowFields {
		decoder.DisallowUnknownFields()
	}

	if err := decoder.Decode(dest); err != nil {
		return err
	}

	return nil
}

func GetAccessTokenFromCode(ctx context.Context, code, appID, secret string, response *GetAccessTokenResponse) error {
	if code == "" || appID == "" || secret == "" {
		return fmt.Errorf("cannot get access_token: invalid argument")
	}

	// get access_token
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		appID, secret, code)
	if err := GetHTTPResponse(url, response, true); err != nil {
		return err
	}

	log.Debugf("GetAccessTokenFromCode response: %v", response)

	return nil
}

func GetWeChatUserInfoFromAccessToken(ctx context.Context, accessToken, openID string, response *GetWeChatUserInfoResponse) error {
	// get user info
	if accessToken == "" || openID == "" {
		return fmt.Errorf("cannot get wechat user info: invalid argument")
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s",
		accessToken, openID)
	if err := GetHTTPResponse(url, response, false); err != nil {
		return err
	}

	if response.UnionID == "" {
		return fmt.Errorf("unionid is required, cannot be empty string")
	}

	log.Debugf("GetAccessTokenFromCode response: %v", response)

	return nil
}
