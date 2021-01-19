package email

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	. "github.com/mxc-foundation/lpwan-app-server/internal/email/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/email/tlssmtp"
	"github.com/mxc-foundation/lpwan-app-server/internal/static"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "email"

type client struct {
	senderID    string
	username    string
	password    string
	authType    string
	smtpHost    string
	smtpPort    string
	tlsRequired bool
}

type operatorInfo struct {
	operatorName,
	downloadAppStore,
	downloadGoogle,
	downloadTestFlight,
	downloadAPK,
	operatorAddress,
	operatorLegal,
	operatorLogo,
	operatorContact,
	operatorSupport string
}

var email struct {
	base32endocoding *base32.Encoding
	from             string
	host             string
	operator         operatorInfo
	mailTemplates    map[EmailOptions]*template.Template
}

type controller struct {
	s        ServerInfoStruct
	operator OperatorStruct
	smtp     map[string]SMTPStruct
	cli      map[string]*client

	moduleUp bool
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) error {

	ctrl = &controller{
		operator: conf.Operator,
		smtp:     conf.SMTP,
		s: ServerInfoStruct{
			ServerAddr:      conf.General.ServerAddr,
			DefaultLanguage: conf.General.DefaultLanguage,
		},
		cli: make(map[string]*client),
	}

	return nil
}

// GetSettings returns ServerInfoStruct
func GetSettings() ServerInfoStruct {
	return ctrl.s
}

// GetOperatorInfo returns OperatorStruct
func GetOperatorInfo() OperatorStruct {
	return ctrl.operator
}

// Setup configures the package.
func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp == true {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	for key, value := range ctrl.smtp {
		port := value.Port
		if port == "" {
			port = "25"
		}

		ctrl.cli[key] = &client{
			senderID:    value.Email,
			username:    value.Username,
			password:    value.Password,
			authType:    value.AuthType,
			smtpHost:    value.Host,
			smtpPort:    port,
			tlsRequired: value.TLSRequired,
		}

		email.from = value.Email
	}

	email.base32endocoding = base32.StdEncoding.WithPadding(base32.NoPadding)
	email.host = "https://" + ctrl.s.ServerAddr
	email.operator = operatorInfo{
		operatorName:       ctrl.operator.Operator,
		downloadAppStore:   ctrl.operator.DownloadAppStore,
		downloadGoogle:     ctrl.operator.DownloadGoogle,
		downloadTestFlight: ctrl.operator.DownloadTestFlight,
		downloadAPK:        ctrl.operator.DownloadAPK,
		operatorAddress:    ctrl.operator.OperatorAddress,
		operatorLegal:      ctrl.operator.OperatorLegal,
		operatorLogo:       ctrl.operator.OperatorLogo,
		operatorContact:    ctrl.operator.OperatorContact,
		operatorSupport:    ctrl.operator.OperatorSupport,
	}

	if err := loadEmailTemplates(); err != nil {
		return err
	}

	return nil
}

type Mailer struct{}

func (m *Mailer) SendVerifyEmailConfirmation(email, lang, securityToken string) error {
	return SendInvite(email, Param{Token: securityToken}, EmailLanguage(lang), VerifyEmail)
}

func (m *Mailer) SendRegistrationConfirmation(email, lang, securityToken string) error {
	return SendInvite(email, Param{Token: securityToken}, EmailLanguage(lang), RegistrationConfirmation)
}

func (m *Mailer) SendPasswordResetUnknown(email, lang string) error {
	return SendInvite(email, Param{}, EmailLanguage(lang), PasswordResetUnknown)
}

func (m *Mailer) SendPasswordReset(email, lang, otp string) error {
	return SendInvite(email, Param{Token: otp}, EmailLanguage(lang), PasswordReset)
}

// SendInvite ...
func SendInvite(user string, param Param, language EmailLanguage, option EmailOptions) error {
	var err error

	if ctrl.cli == nil {
		log.Error("Tried to send registration email, but SMTP is not configured")
		return errors.New("Unable to send confirmation email")
	}

	if email.mailTemplates[option] == nil {
		log.Errorf("Email template for %s does not exist", option)
		return errors.New("Unable to send confirmation email")
	}

	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	param.messageID = time.Now().Format("20060102150405.") + email.base32endocoding.EncodeToString(b)

	// get json object for this language, this context only exists in main template
	mainContentName := fmt.Sprintf("%s/main-template/main-template-%s.json",
		EmailTextPath, string(language))
	defaultMainContentName := fmt.Sprintf("%s/main-template/main-template-en.json", EmailTextPath)

	jsonDataMain, err := static.Asset(mainContentName)
	if err != nil {
		log.WithError(err).Warnf("%s does not exist", mainContentName)

		jsonDataMain, err = static.Asset(defaultMainContentName)
		if err != nil {
			log.WithError(err).Errorf("%s does not exist", defaultMainContentName)
			return err
		}
	}
	err = json.Unmarshal(jsonDataMain, &param.commonJSON)
	if err != nil {
		log.WithError(err).Errorf("Parse json data error")
		return err
	}

	// get json object for this language and this email option
	assetName := fmt.Sprintf("%s/%s/%s-%s.json",
		EmailTextPath, string(option), string(option), string(language))

	// always use English as default language
	defaultAssetName := fmt.Sprintf("%s/%s/%s-en.json", EmailTextPath, string(option), string(option))

	jsonData, err := static.Asset(assetName)
	if err != nil {
		log.WithError(err).Warnf("Email template does not support %s with %s", string(option), string(language))

		jsonData, err = static.Asset(defaultAssetName)
		if err != nil {
			log.WithError(err).Errorf("Email template does not support %s ", string(option))
			return err
		}
	}

	data, err := emailOptionsList[option].getEmailParam(user, param, jsonData)
	if err != nil {
		log.Error(err)
		return err
	}

	var msg bytes.Buffer
	if err := email.mailTemplates[option].ExecuteTemplate(&msg, EmailTemplateMain, &data); err != nil {
		log.Error(err)
		return err
	}

	str := strings.Replace(msg.String(), "=\"", "=3D\"", -1)
	out := bytes.NewBufferString(str)

	for k, v := range ctrl.cli {
		if v != nil {
			err = v.send(user, *out)
			if err == nil {
				return nil
			}
			log.WithError(err).Warnf("Failed to send email with %s, try with other provider", k)
		}
	}

	log.Error("Unable to send confirmation email")
	return errors.New("SMTP server failed")

}

func (c *client) send(user string, msg bytes.Buffer) error {
	var auth smtp.Auth
	if c.authType == "PLAIN" {
		auth = smtp.PlainAuth("", c.username, c.password, c.smtpHost)
	} else if c.authType == "CRAM-MD5" {
		auth = smtp.CRAMMD5Auth(c.username, c.password)
	} else if c.authType != "" {
		return fmt.Errorf("unsupported authentication type: %s", c.authType)
	}

	var err error
	if c.tlsRequired {
		err = tlssmtp.SendMail(c.smtpHost+":"+c.smtpPort, auth, c.senderID, []string{user}, msg.Bytes())
	} else {
		err = smtp.SendMail(c.smtpHost+":"+c.smtpPort, auth, c.senderID, []string{user}, msg.Bytes())
	}
	if err != nil {
		return fmt.Errorf("couldn't send an email: %v", err)
	}

	log.Infof("email sent to %s", user)

	return nil
}

// only for debugging purpose
/*func writeMsgToFile(filename string, msg bytes.Buffer) {
	log.Infof("write msg to file %s", filename)

	f, err := os.Create(filename)
	if err != nil {
		log.WithError(err).Error("deubg: writeMsgToFile")
		return
	}
	defer f.Close()

	_, err = f.Write(msg.Bytes())
	if err != nil {
		log.WithError(err).Error("deubg: writeMsgToFile")
		return
	}

	err = f.Sync()
	if err != nil {
		log.WithError(err).Error("deubg: writeMsgToFile")
		return
	}

}*/
