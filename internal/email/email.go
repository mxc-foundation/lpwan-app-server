package email

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"net/smtp"
	"os"
	"runtime"
	"text/template"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/static"
	log "github.com/sirupsen/logrus"
)

var (
	senderID   string
	password   string
	host       string
	smtpServer string
	smtpPort   string
	disable    bool
)

const (
	sendInvite = iota
)

// Setup configures the package.
func Setup(c config.Config) error {
	senderID = c.SMTP.Email
	password = c.SMTP.Password
	smtpServer = c.SMTP.Host
	smtpPort = c.SMTP.Port
	host = os.Getenv("APPSERVER")
	disable = false

	base32endocoding = base32.StdEncoding.WithPadding(base32.NoPadding)
	mailTemplates = make([]*template.Template, len(mailTemplateNames))
	for k, v := range mailTemplateNames {
		mailTemplates[k] = template.Must(
			template.New(v.templatePath).Parse(string(static.MustAsset(v.templatePath))))
	}

	return nil
}

// Disable stops emails from being sent.
func Disable() {
	disable = true
}

var (
	mailTemplates    []*template.Template
	base32endocoding *base32.Encoding

	mailTemplateNames = []struct {
		templatePath string
		url          string
	}{
		sendInvite: {
			templatePath: "templates/registration-confirm",
			url:          "/#/registration-confirm/",
		},
	}
)

// SendInvite ...
func SendInvite(user string, token string) error {
	var err error

	if disable == true {
		return nil
	}

	if smtpServer == "" {
		log.Error("Tried to send registration email, but SMTP is not configured")
		return errors.New("Unable to send confirmation email")
	}

	link := host + mailTemplateNames[sendInvite].url + token

	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	messageID := time.Now().Format("20060102150405.") + base32endocoding.EncodeToString(b)

	var msg bytes.Buffer
	if err := mailTemplates[sendInvite].Execute(&msg, struct {
		From, To, Host, MsgId, Boundary, Link string
	}{
		From:     senderID,
		To:       user,
		Host:     host,
		MsgId:    messageID + "@" + host,
		Boundary: "----=_Part_" + messageID,
		Link:     link,
	}); err != nil {
		log.Error(err)
		return err
	}

	runtime.Breakpoint()
	err = smtp.SendMail(smtpServer+":"+smtpPort,
		smtp.CRAMMD5Auth(senderID, password), senderID, []string{user}, msg.Bytes())

	if err != nil {
		return err
	}

	return nil
}
