package email

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/smtp"
	"time"
)

type registrationEmailInterface struct {
	option EmailOptions
}

var registrationEmail = registrationEmailInterface{option: RegistrationConfirmation}

func (s *registrationEmailInterface) sendEmail(user, token string, language EmailLanguage) error {
	link := host + mailTemplateNames[s.option][language].url + token

	logo := host + "/branding.png"

	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	messageID := time.Now().Format("20060102150405.") + base32endocoding.EncodeToString(b)

	var msg bytes.Buffer
	if err := mailTemplates[s.option][language].Execute(&msg, struct {
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

	err := smtp.SendMail(smtpServer+":"+smtpPort,
		smtp.CRAMMD5Auth(senderID, password), senderID, []string{user}, msg.Bytes())

	return err
}
