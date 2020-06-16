package email

import (
	"bytes"
	"crypto/tls"
	"encoding/base32"
	"net"
	"net/smtp"
	"os"
	"text/template"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
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
	host = os.Getenv("APPSERVER")
	serverRegion := os.Getenv("SERVER_REGION")
	if serverRegion == pb.ServerRegion_name[int32(pb.ServerRegion_RESTRICTED)] {
		senderID = c.SMTP.Restricted.Email
		password = c.SMTP.Restricted.Password
		smtpServer = c.SMTP.Restricted.Host
		smtpPort = c.SMTP.Restricted.Port

	} else {
		senderID = c.SMTP.Average.Email
		password = c.SMTP.Average.Password
		smtpServer = c.SMTP.Average.Host
		smtpPort = c.SMTP.Average.Port
	}

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

func sendEmail(user string, msg bytes.Buffer) error {
	serverRegion := os.Getenv("SERVER_REGION")

	if serverRegion == pb.ServerRegion_name[int32(pb.ServerRegion_RESTRICTED)] {
		return sendEmailRestricted(user, msg)
	}

	// serverRegion == pb.ServerRegion_name[int32(pb.ServerRegion_NOT_DEFINED)]
	// serverRegion == pb.ServerRegion_name[int32(pb.ServerRegion_AVERAGE)]
	// serverRegion == ""
	return sendEmailAverage(user, msg)
}

func sendEmailAverage(user string, msg bytes.Buffer) error {
	return smtp.SendMail(smtpServer+":"+smtpPort,
		smtp.CRAMMD5Auth(senderID, password), senderID, []string{user}, msg.Bytes())
}

func sendEmailRestricted(user string, msg bytes.Buffer) error {
	return SendMailUsingTLS(
		smtpServer+":465",
		smtp.PlainAuth(
			"",
			senderID,
			password,
			smtpServer,
		),
		senderID,
		[]string{user},
		msg.Bytes(),
	)
}

//return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}
	//split host address and port
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
