package email

import (
	"bytes"
	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver_serves_ui"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type registrationEmailInterface struct {
	option EmailOptions
}

var registrationEmail = registrationEmailInterface{option: RegistrationConfirmation}

func (s *registrationEmailInterface) sendEmail(user, token string, language EmailLanguage) error {
	mailTemplate := mailTemplates[s.option][language]
	if mailTemplate == nil {
		mailTemplate = mailTemplates[s.option][EmailLanguage(pb.Language_name[int32(pb.Language_en)])]
	}

	link := host + mailTemplateNames[s.option][language].url + token

	logo := host + "/branding.png"

	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	messageID := time.Now().Format("20060102150405.") + base32endocoding.EncodeToString(b)

	var msg bytes.Buffer
	if err := mailTemplate.Execute(&msg, struct {
		From, To, Host, MsgID, Boundary, Link, Logo, Operator, PrimaryColor, SecondaryColor string
	}{
		From:     senderID,
		To:       user,
		Host:     host,
		MsgID:    messageID + "@" + host,
		Boundary: "----=_Part_" + messageID,
		Link:     link,
		Logo:     logo,
		Operator: "MXC",
		PrimaryColor: "#71B6F9",
		SecondaryColor: "#10c469",
	}); err != nil {
		log.Error(err)
		return err
	}

	err := sendEmail(user, msg)

	return err
}
