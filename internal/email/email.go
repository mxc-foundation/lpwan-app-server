package email

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"net/smtp"
	"os"
	"strings"
	"text/template"

	"github.com/brocaar/lora-app-server/internal/config"
	"github.com/brocaar/lora-app-server/internal/static"
	log "github.com/sirupsen/logrus"
)

var conf Conf

// Setup configures the package.
func Setup(c config.Config) error {
	conf = Conf{
		senderID: c.SMTP.Email,
		password: c.SMTP.Password,
		host:     c.SMTP.Host,
		port:     c.SMTP.Port,
	}

	return nil
}

// Disable stops emails from being sent.
func Disable() {
	conf = Conf{
		disable: true,
	}
}

// Conf ...
type Conf struct {
	senderID string
	password string
	host     string
	port     string
	disable  bool
}

// Mail ...
type Mail struct {
	Sender    string
	Recepient string
	Host      string
	MessageID string
	Link      string
}

// SendInvite ...
func SendInvite(user string, token string) error {
	var err error

	if conf.disable == true {
		return nil
	}
	if conf.host == "" {
		log.Error("Tried to send registration email, but SMTP is not configured")
		return errors.New("Unable to send confirmation email")
	}

	var mail = Mail{
		Sender:    conf.senderID,
		Recepient: user,
	}

	if mail.Host, err = os.Hostname(); err != nil {
		return err
	}

	if !strings.ContainsRune(mail.Host, '.') {
		mail.Host = ".matchx.io"
	}

	localHostAddr := os.Getenv("LOCAL_HOST_ADDRESS")
	if localHostAddr != "" {
		mail.Link = localHostAddr + "/#/registration-confirm/" + token
	} else {
		mail.Link = "https://" + mail.Host + "/#/registration-confirm/" + token
	}
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	mail.MessageID = base32.StdEncoding.EncodeToString(b)
	content, _ := static.Asset("templates/registration-confirm.txt")
	t := template.New("Invite")
	tpl, _ := t.Parse(string(content))
	if err != nil {
		return err
	}

	var msg bytes.Buffer
	if err = tpl.Execute(&msg, mail); err != nil {
		return err
	}

	err = smtp.SendMail(conf.host+":"+conf.port,
		smtp.CRAMMD5Auth(mail.Sender, conf.password),
		mail.Sender, []string{mail.Recepient}, msg.Bytes())

	if err != nil {
		return err
	}

	return nil
}
