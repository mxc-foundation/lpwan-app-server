package email

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"testing"
)

type testStore struct {
	smtp     config.SMTPStruct
	operator config.OperatorStruct
}

func newTestStore() *testStore {
	return &testStore{
		smtp: config.SMTPStruct{
			Email:    "do-not-reply@mxc.org",
			Username: "AKIAYLLLRKLATNOTY3F5",
			Password: "BN3u+u/u2JLwF2rgY2zonMYfPQHz/D8ycAExRzTEJbDd",
			AuthType: "PLAIN",
			Host:     "email-smtp.eu-central-1.amazonaws.com",
			Port:     "587",
			/*			Email:    "do-not-reply@mxc.org",
						Username: "apikey",
						Password: "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",
						AuthType: "PLAIN",
						Host:     "smtp.sendgrid.net",
						Port:     "587",*/
		},
		operator: config.OperatorStruct{
			Operator:         "MatchX Test",
			DownloadAppStore: "https://apps.apple.com/app/mxc-datadash/id1509218470",
			DownloadGoogle:   "https://play.google.com/store/apps/details?id=com.mxc.smartcity",
			OperatorAddress:  "Brückenstraße 4, 10319 Berlin awesome@matchx.io",
			OperatorLegal:    "MatchX GmbH",
			OperatorLogo:     "https://lora.supernode.matchx.io/branding.png",
		},
	}
}

func TestSendInvite(t *testing.T) {
	ts := newTestStore()

	// setup settings and load templates
	if err := Setup(config.Config{
		SMTP:     ts.smtp,
		Operator: ts.operator,
	}); err != nil {
		t.Fatalf("%v", err)
	}

	/*	languageList := []string{"zhCN"}
		for _, language := range languageList {

			for option := range emailOptionsList {
				time.Sleep(5 * time.Second)
				if err := SendInvite("lixuan@mxc.org", Param{Token: "1234567890"}, EmailLanguage(language), option); err != nil {
					continue
				}
			}
		}*/
	SendInvite("lixuan@mxc.org", Param{Token: "1234567890"}, "zhCN", RegistrationConfirmation)
}
