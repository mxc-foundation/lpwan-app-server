package tlssmtp

import (
	"crypto/tls"
	"errors"
	"net"
	"net/smtp"
)

// SendMail connects to the server at addr using TLS, authenticates with the
// optional mechanism a if possible, and then sends email from address from to
// addresses to with message msg.
func SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	tlsConn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host})
	if err != nil {
		return err
	}
	defer tlsConn.Close()
	c, err := smtp.NewClient(tlsConn, host)
	if err != nil {
		return err
	}
	defer c.Close()
	if a != nil {
		if ok, _ := c.Extension("AUTH"); !ok {
			return errors.New("smtp server doesn't support AUTH")
		}
		if err := c.Auth(a); err != nil {
			return err
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
	if _, err := w.Write(msg); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.Quit()
}
