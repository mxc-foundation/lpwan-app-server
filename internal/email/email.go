package email

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"net/smtp"
	"os"
	"text/template"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	senderID   string
	password   string
	host       string
	smtpServer string
	smtpPort   string
)

type mailTemplateStruct struct {
	templatePath string
	url          string
}

type mailTemplatesType map[EmailOptions]map[EmailLanguage]*template.Template
type mailTemplateNamesType map[EmailOptions]map[EmailLanguage]mailTemplateStruct

var (
	mailTemplates     mailTemplatesType
	base32endocoding  *base32.Encoding
	mailTemplateNames mailTemplateNamesType
)

// Setup configures the package.
func Setup(c config.Config) error {
	senderID = c.SMTP.Email
	password = c.SMTP.Password
	smtpServer = c.SMTP.Host
	smtpPort = c.SMTP.Port
	host = os.Getenv("APPSERVER")

	base32endocoding = base32.StdEncoding.WithPadding(base32.NoPadding)

	loadEmailTemplates()

	return nil
}

// SendInvite ...
func SendInvite(user, token string, language EmailLanguage, option EmailOptions) error {
	var err error

	if smtpServer == "" {
		log.Error("Tried to send registration email, but SMTP is not configured")
		return errors.New("Unable to send confirmation email")
	}

	link := host + mailTemplateNames[option][language].url + token

	logo := host + "/branding.png"

	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	messageID := time.Now().Format("20060102150405.") + base32endocoding.EncodeToString(b)

	var msg bytes.Buffer
	if err := mailTemplates[option][language].Execute(&msg, struct {
		From, To, Host, MsgID, Boundary, Link, Logo string
	}{
		From:     senderID,
		To:       user,
		Host:     host,
		MsgID:    messageID + "@" + host,
		Boundary: "----=_Part_" + messageID,
		Link:     link,
		Logo:     logo,
	}); err != nil {
		log.Error(err)
		return err
	}

	err = smtp.SendMail(smtpServer+":"+smtpPort,
		smtp.CRAMMD5Auth(senderID, password), senderID, []string{user}, msg.Bytes())

	return errors.Wrap(err, "")
}
