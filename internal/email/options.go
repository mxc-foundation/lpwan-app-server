package email

import (
	"bytes"
	"html/template"

	"github.com/mxc-foundation/lpwan-app-server/internal/static"
)

type EmailOptions string
type EmailLanguage string

const (
	EmailTemplatePath string = "email/templates"
	EmailTextPath     string = "email/text"

	EmailTemplateMain string = "mainTemplate"
	EmailTemplateHead string = "htmlBodyPartOne"
	BodyTemplateName  string = "bodyTemplate"
	EmailTemplateTail string = "htmlBodyPartTwo"

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

type Param struct {
	Token      string
	Amount     string
	messageID  string
	commonJSON struct {
		Str1 string `json:"str1"`
		Str2 string `json:"str2"`
		Str3 string `json:"str3"`
		Str4 string `json:"str4"`
		Str5 string `json:"str5"`
		Str6 string `json:"str6"`
	}
}

// define interfaces for each email option
type emailInterface interface {
	getEmailParam(user string, param Param, jsonData []byte) (interface{}, error)
}

var emailOptionsList = map[EmailOptions]emailInterface{
	RegistrationConfirmation: registrationInterface,
	TwoFALogin:               twofaLogin,
	TwoFAWithdraw:            twoFAWithdraw,
	PasswordReset:            passwordResetInterface,
	PasswordResetUnknown:     passwordResetUnknownInterface,
	/*		StakingIncome:            stakingIncome,
			TopupConfirmation:        topupConfirmation,
			WithdrawDenied:           withdrawDenied,
			WithdrawSuccess:          withdrawSuccess,*/
}

var (
	registrationInterface         = emailInterface(&registrationEmail)
	twofaLogin                    = emailInterface(&twoFALoginEmail)
	twoFAWithdraw                 = emailInterface(&twoFAWithdrawEmail)
	passwordResetInterface        = emailInterface(&passwordResetEmail)
	passwordResetUnknownInterface = emailInterface(&passwordResetUnknownEmail)
	/*	stakingIncome                 = emailInterface(&stakingIncomeEmail)
		topupConfirmation             = emailInterface(&topupConfirmEmail)
		withdrawDenied                = emailInterface(&withdrawDeniedEmail)
		withdrawSuccess               = emailInterface(&withdrawSuccessEmail)*/
)

func loadEmailTemplates() error {
	email.mailTemplates = make(map[EmailOptions]*template.Template)

	templatePath0 := EmailTemplatePath + "/" + "email_template"
	templatePath1 := EmailTemplatePath + "/" + "email_template_1"
	_ = static.MustAsset(templatePath0)
	_ = static.MustAsset(templatePath1)

	// get list of existing templates
	for option := range emailOptionsList {
		bodyTemplatePath := EmailTemplatePath + "/" + string(option)

		if _, err := static.Asset(bodyTemplatePath); err != nil {
			continue
		}

		tpl := template.New(EmailTemplateMain)

		// provide a func in the FuncMap which can access tpl to be able to look up templates
		tpl.Funcs(map[string]interface{}{
			"CallTemplate": func(name string, data interface{}) (ret template.HTML, err error) {
				buf := bytes.NewBuffer([]byte{})
				err = tpl.ExecuteTemplate(buf, name, data)
				// #nosec: this method will not auto-escape HTML. Verify data is well formed.
				ret = template.HTML(buf.String())
				return
			},
		})

		email.mailTemplates[option] = template.Must(tpl.Parse(string(static.MustAsset(templatePath0))))
		email.mailTemplates[option] = template.Must(tpl.New(EmailTemplateHead).Parse(string(static.MustAsset(templatePath1))))
		email.mailTemplates[option] = template.Must(tpl.New(BodyTemplateName).Parse(string(static.MustAsset(bodyTemplatePath))))
	}

	return nil
}
