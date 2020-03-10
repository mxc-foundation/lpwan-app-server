package email

import (
	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver_serves_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/static"
	"os"
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
		tmpNames := make(map[EmailLanguage]mailTemplateStruct)

		for _, language := range pb.Language_name {
			tmpNames[EmailLanguage(language)] = mailTemplateStruct{
				templatePath: "templates/email/" + string(option) + "/" + string(option) + "-" + language,
			}

			/*mailTemplateNames[option][EmailLanguage(language)] = mailTemplateStruct{
				templatePath: "templates/email/" + string(option) + "/" + string(option) + "-" + language,
			}*/
		}

		mailTemplateNames[option] = tmpNames

		if option == RegistrationConfirmation {
			for _, language := range pb.Language_name {
				mailTemplateNames[option][EmailLanguage(language)] = mailTemplateStruct{
					url: "/#/registration-confirm/",
				}
			}
		}

	}

	tmpTemplates := make(map[EmailLanguage]*template.Template)

	for option, _ := range emailOptionsList {
		for _, language := range pb.Language_name {
			filePath := mailTemplateNames[option][EmailLanguage(language)].templatePath
			_, err := os.Stat(filePath)
			if os.IsNotExist(err) {
				continue
			}

			tmpTemplates[EmailLanguage(language)] = template.Must(
				template.New(mailTemplateNames[option][EmailLanguage(language)].templatePath).Parse(
					string(static.MustAsset(mailTemplateNames[option][EmailLanguage(language)].templatePath))))

			/*			mailTemplates[option][EmailLanguage(language)] = template.Must(
						template.New(mailTemplateNames[option][EmailLanguage(language)].templatePath).Parse(
							string(static.MustAsset(mailTemplateNames[option][EmailLanguage(language)].templatePath))))*/
		}

		mailTemplates[option] = tmpTemplates
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
