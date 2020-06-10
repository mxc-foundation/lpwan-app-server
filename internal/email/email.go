package email

import (
	"bytes"
	"encoding/base32"
	"fmt"
	"net/smtp"
	"os"
	"text/template"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
)

var cli *Client

var (
	senderID string
	host     string
)

type Client struct {
	senderID string
	username string
	password string
	authType string
	host     string
	smtpHost string
	smtpPort string
}

type mailTemplateStruct struct {
	templatePath string
	url          string
}

type mailTemplatesType map[EmailOptions]map[EmailLanguage]*template.Template
type mailTemplateNamesType map[EmailOptions]map[EmailLanguage]mailTemplateStruct

var (
	mailTemplates     = mailTemplatesType{}
	base32endocoding  *base32.Encoding
	mailTemplateNames = mailTemplateNamesType{}
)

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

	senderID = cli.senderID
	host = cli.host

	base32endocoding = base32.StdEncoding.WithPadding(base32.NoPadding)

	loadEmailTemplates()

	return nil
}

// SendInvite ...
func SendInvite(user, token string, language EmailLanguage, option EmailOptions) error {
	var err error

	if cli == nil {
		log.Error("Tried to send registration email, but SMTP is not configured")
		return errors.New("Unable to send confirmation email")
	}

	err = emailOptionsList[option].sendEmail(user, token, language)

	return errors.Wrap(err, "")
}

func (c *Client) sendEmail(user string, msg bytes.Buffer) error {
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

	return nil
}
