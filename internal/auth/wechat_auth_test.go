package auth

import (
	"context"
	"testing"

	"github.com/mxc-foundation/lpwan-app-server/internal/httpcli"
)

type testRespIncomplete struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type testRespComplete struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	ExtraField string `json:"extra_field"`
}

func TestGetHTTPResponse(t *testing.T) {
	tests := []struct {
		name                  string
		url                   string
		disallowUnknownFields bool
		respIncomplete        testRespIncomplete
		respComplete          testRespComplete
	}{
		{
			name:                  "allow unknown fields in response",
			url:                   "https://lora.build.cloud.mxc.org/test/",
			disallowUnknownFields: false,
			respIncomplete:        testRespIncomplete{},
			respComplete:          testRespComplete{},
		},
		{
			name:                  "disallow unknown fields in response",
			url:                   "https://lora.build.cloud.mxc.org/test/",
			disallowUnknownFields: true,
			respIncomplete:        testRespIncomplete{},
			respComplete:          testRespComplete{},
		},
	}

	for _, tc := range tests {
		t.Logf(tc.name)
		err := httpcli.GetResponse(tc.url, &tc.respComplete, tc.disallowUnknownFields)
		if err != nil {
			t.Errorf("expected: no error, got error %v", err)
		}

		err = httpcli.GetResponse(tc.url, &tc.respIncomplete, tc.disallowUnknownFields)
		if tc.disallowUnknownFields {
			if err == nil {
				t.Errorf("expected: error %v, got no error", err)
			}
		} else {
			if err != nil {
				t.Errorf("expected: no error, got error %v", err)
			}
		}
	}
}

func TestGetAccessTokenFromCode(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name   string
		ctx    context.Context
		urlStr string
		code   string
		appID  string
		secret string
		resp   GetAccessTokenResponse
		noErr  bool
	}{
		{
			name:   "input code is empty string",
			ctx:    ctx,
			urlStr: "https://lora.build.cloud.mxc.org/test/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
			code:   "",
			appID:  "123456",
			secret: "123456",
			resp:   GetAccessTokenResponse{},
			noErr:  false,
		},
		{
			name:   "input appID is empty string",
			ctx:    ctx,
			urlStr: "https://lora.build.cloud.mxc.org/test/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
			code:   "123456",
			appID:  "",
			secret: "123456",
			resp:   GetAccessTokenResponse{},
			noErr:  false,
		},
		{
			name:   "input secret is empty string",
			ctx:    ctx,
			urlStr: "https://lora.build.cloud.mxc.org/test/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
			code:   "123456",
			appID:  "123456",
			secret: "",
			resp:   GetAccessTokenResponse{},
			noErr:  false,
		},
		{
			name:   "all input valid",
			ctx:    ctx,
			urlStr: "https://lora.build.cloud.mxc.org/test/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
			code:   "123456",
			appID:  "123456",
			secret: "123456",
			resp:   GetAccessTokenResponse{},
			noErr:  true,
		},
	}

	for _, tc := range tests {
		t.Logf(tc.name)
		err := GetAccessTokenFromCode(tc.ctx, tc.code, tc.appID, tc.secret, &tc.resp)
		if tc.noErr != (err == nil) {
			t.Errorf("expected: noError = %v, got %v", tc.noErr, err)
		}
	}
}

func TestGetWeChatUserInfoFromAccessToken(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		ctx         context.Context
		urlStr      string
		accessToken string
		openID      string
		resp        GetWeChatUserInfoResponse
		noErr       bool
	}{
		{
			name:        "input accessToken is empty string",
			ctx:         ctx,
			urlStr:      "https://lora.build.cloud.mxc.org/test/userinfo?access_token=%s&openid=%s",
			accessToken: "",
			openID:      "asdfasdfasdfsadf",
			resp:        GetWeChatUserInfoResponse{},
			noErr:       false,
		},
		{
			name:        "input openID is empty string",
			ctx:         ctx,
			urlStr:      "https://lora.build.cloud.mxc.org/test/userinfo?access_token=%s&openid=%s",
			accessToken: "afasfdasdfas",
			openID:      "",
			resp:        GetWeChatUserInfoResponse{},
			noErr:       false,
		},
		{
			name:        "all input valid",
			ctx:         ctx,
			urlStr:      "https://lora.build.cloud.mxc.org/test/userinfo?access_token=%s&openid=%s",
			accessToken: "afasfdasdfas",
			openID:      "asdfasdfasfds",
			resp:        GetWeChatUserInfoResponse{},
			noErr:       true,
		},
	}

	for _, tc := range tests {
		t.Logf(tc.name)
		err := GetWeChatUserInfoFromAccessToken(tc.ctx, tc.accessToken, tc.openID, &tc.resp)
		if tc.noErr != (err == nil) {
			t.Errorf("expected: noError = %v, got %v", tc.noErr, err)
		}
	}

}
