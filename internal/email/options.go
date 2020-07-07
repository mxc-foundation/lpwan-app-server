package email

import (
	"text/template"

	"github.com/mxc-foundation/lpwan-app-server/internal/static"
)

type EmailOptions string
type EmailLanguage string

const (
	RegistrationConfirmation EmailOptions = "registration-confirm"
	PasswordReset            EmailOptions = "password-reset"
	PasswordResetUnknown     EmailOptions = "password-reset-unknown"
	TwoFALogin               EmailOptions = "2fa-login"
	TwoFAWithdraw            EmailOptions = "2fa-withdraw"
	StakingIncome            EmailOptions = "staking-income"
	TopupConfirmation        EmailOptions = "topup-confirm"
	WithdrawDenied           EmailOptions = "withdraw-denied"
	WithdrawSuccess          EmailOptions = "withdraw-success"
)

var emailOptionsList = map[EmailOptions]emailInterface{
	RegistrationConfirmation: registrationInterface,
	PasswordReset:            passwordReset,
	PasswordResetUnknown:     passwordResetUnknown,
	TwoFALogin:               twofaLogin,
	TwoFAWithdraw:            twoFAWithdraw,
	StakingIncome:            stakingIncome,
	TopupConfirmation:        topupConfirmation,
	WithdrawDenied:           withdrawDenied,
	WithdrawSuccess:          withdrawSuccess,
}

func loadEmailTemplates() {
	// this is not a nice solution and we should move away from it to use
	// whatever languages templates are available in, instead of hardcoding the
	// list of supported languages
	supportedLanguages := []string{"de", "en", "es", "fr", "ja", "ko", "nl", "ru", "zhCN", "zhTW"}
	for option := range emailOptionsList {
		mailTemplateNames[option] = make(map[EmailLanguage]mailTemplateStruct)

		for _, language := range supportedLanguages {
			var url string
			if option == RegistrationConfirmation {
				url = "/#/registration-confirm/"
			} else if option == PasswordReset {
				url = "/#/reset-password-confirm"
			}
			mailTemplateNames[option][EmailLanguage(language)] = mailTemplateStruct{
				templatePath: "templates/email/" + string(option) + "/" + string(option) + "-" + language,
				url:          url,
			}
		}
	}

	for option := range emailOptionsList {
		mailTemplates[option] = make(map[EmailLanguage]*template.Template)

		for _, language := range supportedLanguages {
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
