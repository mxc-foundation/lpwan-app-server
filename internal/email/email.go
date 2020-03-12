package email

import (
	"encoding/base32"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"

	"text/template"
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
	mailTemplates     = mailTemplatesType{}
	base32endocoding  *base32.Encoding
	mailTemplateNames = mailTemplateNamesType{}
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

	err = emailOptionsList[option].sendEmail(user, token, language)

	return errors.Wrap(err, "")
}
