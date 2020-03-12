package email

import (
	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver_serves_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/static"
	"runtime"
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

var emailOptionsList = map[EmailOptions]emailInterface{
	RegistrationConfirmation: registrationInterface,
	TwoFALogin:               twofaLogin,
	TwoFAWithdraw:            twoFAWithdraw,
	StakingIncome:            stakingIncome,
	TopupConfirmation:        topupConfirmation,
	WithdrawDenied:           withdrawDenied,
	WithdrawSuccess:          withdrawSuccess,
}

func loadEmailTemplates() {
	for option, _ := range emailOptionsList {
		mailTemplateNames[option] = make(map[EmailLanguage]mailTemplateStruct)

		for _, language := range pb.Language_name {
			mailTemplateNames[option][EmailLanguage(language)] = mailTemplateStruct{
				templatePath: "templates/email/" + string(option) + "/" + string(option) + "-" + language,
			}
		}

		if option == RegistrationConfirmation {
			for _, language := range pb.Language_name {
				mailTemplateNames[option][EmailLanguage(language)] = mailTemplateStruct{
					templatePath: mailTemplateNames[option][EmailLanguage(language)].templatePath,
					url: "/#/registration-confirm/",
				}
			}
		}
	}

	runtime.Breakpoint()
	for option, _ := range emailOptionsList {
		mailTemplates[option] = make(map[EmailLanguage]*template.Template)

		for _, language := range pb.Language_name {
			_, err := static.AssetInfo(mailTemplateNames[option][EmailLanguage(language)].templatePath)
			if err != nil {
				continue
			}

			mailTemplates[option][EmailLanguage(language)] = template.Must(
				template.New(mailTemplateNames[option][EmailLanguage(language)].templatePath).Parse(
					string(static.MustAsset(mailTemplateNames[option][EmailLanguage(language)].templatePath))))
		}
	}
}

// define interfaces for each email option
type emailInterface interface {
	sendEmail(user, token string, language EmailLanguage) error
}

var (
	registrationInterface = emailInterface(&registrationEmail)
	twofaLogin            = emailInterface(&twoFALoginEmail)
	twoFAWithdraw         = emailInterface(&twoFAWithdrawEmail)
	stakingIncome         = emailInterface(&stakingIncomeEmail)
	topupConfirmation     = emailInterface(&topupConfirmEmail)
	withdrawDenied        = emailInterface(&withdrawDeniedEmail)
	withdrawSuccess       = emailInterface(&withdrawSuccessEmail)
)
