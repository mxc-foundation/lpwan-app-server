package email

import (
	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver_serves_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/static"
	"text/template"
)

type EmailOptions string
type EmailLanguage string

const (
	RegistrationConfirmation EmailOptions = "registration-confirm"
	TwoFALogin               EmailOptions = "2fa-login"
	TwoFAWithdraw            EmailOptions = "2fa-withdraw"
	StakingIncome            EmailOptions = "staking-income"
	TopupConfirmation        EmailOptions = "topup-confirm"
	WithdrawDenied           EmailOptions = "withdraw-denied"
	WithdrawSuccess          EmailOptions = "withdraw-success"
)

var emailOptionsList = []EmailOptions{
	RegistrationConfirmation,
	TwoFALogin,
	TwoFAWithdraw,
	StakingIncome,
	TopupConfirmation,
	WithdrawDenied,
	WithdrawSuccess,
}

func loadEmailTemplates() {
	for _, option := range emailOptionsList {

		for _, language := range pb.Language_name {
			mailTemplateNames[option][EmailLanguage(language)] = mailTemplateStruct{
				templatePath: "templates/email/" + string(option) + "/" + string(option) + "/" + language,
			}
		}

		if option == RegistrationConfirmation {
			for _, language := range pb.Language_name {
				mailTemplateNames[option][EmailLanguage(language)] = mailTemplateStruct{
					url: "/#/registration-confirm/",
				}
			}
		}

	}

	mailTemplates = make(mailTemplatesType, len(mailTemplateNames))

	for _, option := range emailOptionsList {

		for _, language := range pb.Language_name {
			mailTemplates[option][EmailLanguage(language)] = template.Must(
				template.New(mailTemplateNames[option][EmailLanguage(language)].templatePath).Parse(
					string(static.MustAsset(mailTemplateNames[option][EmailLanguage(language)].templatePath))))
		}
	}


}
