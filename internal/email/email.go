package email

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/static"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
)

var cli *Client

type Client struct {
	senderID string
	username string
	password string
	authType string
	host     string
	smtpHost string
	smtpPort string
}

type operatorInfo struct {
	from,
	host,
	MXCLogo,
	operatorName,
	downloadAppStore,
	downloadGoogle,
	appStoreLogo,
	androidLogo,
	operatorAddress,
	operatorLegal,
	operatorLogo,
	operatorContact,
	operatorSupport string
}

var email struct {
	base32endocoding *base32.Encoding
	operator         operatorInfo
	mailTemplates    map[EmailOptions]*template.Template
}

// Setup configures the package.
func Setup(c config.Config) error {
	port := c.SMTP.Port
	if port == "" {
		port = "25"
	}
	cli = &Client{
		senderID: c.SMTP.Email,
		host:     os.Getenv("APPSERVER"),

		username: c.SMTP.Username,
		password: c.SMTP.Password,
		authType: c.SMTP.AuthType,
		smtpHost: c.SMTP.Host,
		smtpPort: port,
	}

	email.base32endocoding = base32.StdEncoding.WithPadding(base32.NoPadding)

	email.operator = operatorInfo{
		from:             cli.senderID,
		host:             cli.host,
		MXCLogo:          c.General.MXCLogo,
		operatorName:     c.Operator.Operator,
		downloadAppStore: c.Operator.DownloadAppStore,
		downloadGoogle:   c.Operator.DownloadGoogle,
		appStoreLogo:     c.Operator.AppStoreLogo,
		androidLogo:      c.Operator.AndroidLogo,
		operatorAddress:  c.Operator.OperatorAddress,
		operatorLegal:    c.Operator.OperatorLegal,
		operatorLogo:     c.Operator.OperatorLogo,
		operatorContact:  c.Operator.OperatorContact,
		operatorSupport:  c.Operator.OperatorSupport,
	}

	if err := loadEmailTemplates(); err != nil {
		return err
	}

	return nil
}

// SendInvite ...
func SendInvite(user string, param Param, language EmailLanguage, option EmailOptions) error {
	var err error

	if cli == nil {
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

	writeMsgToFile(fmt.Sprintf("%s-%s", option, language), msg)

	err = cli.send(user, msg)

	if err != nil {
		log.WithError(err).Error("Unable to send confirmation email")
		return err
	}

	return nil
}

func (c *Client) send(user string, msg bytes.Buffer) error {
	var auth smtp.Auth
	if c.authType == "PLAIN" {
		auth = smtp.PlainAuth("", c.username, c.password, c.smtpHost)
	} else if c.authType == "CRAM-MD5" {
		auth = smtp.CRAMMD5Auth(c.username, c.password)
	} else if c.authType != "" {
		return fmt.Errorf("unsupported authentication type: %s", c.authType)
	}

	err := smtp.SendMail(c.smtpHost+":"+c.smtpPort, auth, c.senderID, []string{user}, msg.Bytes())
	if err != nil {
		return fmt.Errorf("couldn't send an email: %v", err)
	}

	log.Infof("email sent to %s", user)

	return nil
}

// only for debugging purpose
func writeMsgToFile(filename string, msg bytes.Buffer) {
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

	return
}
