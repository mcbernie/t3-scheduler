package schedule

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/aerth/mbox"
)

func (w *watch) notifiyUsers(s *Schedule) {
	s.sendMail(w)
}

type mailing struct {
	email   string
	subject string
	body    string
}

func createMail(u *feuser, p *watch, s *Schedule) mailing {
	mailSection, err := s.cfgIni.GetSection("mail")
	if err != nil {
		panic("no mail configuration found!")
	}

	headerTemplate := mailSection.Key("header").String()
	bodyTemplate := mailSection.Key("body").String()

	header := headerTemplate
	header = strings.Replace(header, "%user%", u.username, -1)
	header = strings.Replace(header, "%pagename%", p.title, -1)

	body := bodyTemplate
	body = strings.Replace(body, "%user%", u.username, -1)
	body = strings.Replace(body, "%pagename%", p.title, -1)
	body = strings.Replace(body, "\\n", "\n", -1)

	return mailing{
		email:   u.mail,
		subject: header,
		body:    body,
	}
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{
		Name:    String,
		Address: "",
	}
	return strings.Trim(addr.String(), " <@>")
}

func (s *Schedule) sendMail(p *watch) {
	mailConfig := s.cfgTypo3.MailConfiguration

	if mailConfig.Transport == "mbox" {

		for _, u := range p.users {

			mailing := createMail(&u, p, s)

			mbox.Destination = mailing.email
			mbox.Open(mailConfig.TransportMboxFile)

			var form mbox.Form
			form.From = mailConfig.FromMailAddress
			form.Subject = mailing.subject
			form.Body = []byte(mailing.body)

			log.Printf("MBOX from address: %s, username: %s, mail: %s\n", mailConfig.FromMailAddress, u.username, u.mail)

			mbox.Save(&form)
		}
	} else if mailConfig.Transport == "smtp" {
		host, _, _ := net.SplitHostPort(mailConfig.TransportSMTPServer)
		auth := smtp.PlainAuth("", mailConfig.TransportSMTPUsername, mailConfig.TransportSMTPPassword, host)

		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}

		fmt.Println(host)

		for _, u := range p.users {

			mailing := createMail(&u, p, s)
			to := []string{mailing.email}
			msg := []byte("To: " + mailing.email + "\r\n" +
				"From: " + mailConfig.FromMailName + "<" + mailConfig.FromMailAddress + ">\r\n" +
				"Subject: " + mailing.subject + "\r\n" +
				"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
				"\r\n" +
				mailing.body + "\r\n")
			fmt.Println(mailing.body)
			log.Printf("SMTP from address: %s, username: %s, mail: %s\n", mailConfig.FromMailAddress, u.username, u.mail)

			if mailConfig.TransportSMTPEncrypt == "ssl" {
				conn, err := tls.Dial("tcp", mailConfig.TransportSMTPServer, tlsconfig)
				if err != nil {
					log.Printf("Error on Sending mail: %s\n", err.Error())
					continue
				}
				c, err := smtp.NewClient(conn, host)
				if err != nil {
					log.Println(err)
					continue
				}

				// Auth
				if err = c.Auth(auth); err != nil {
					log.Println(err)
					continue
				}

				// To && From
				if err = c.Mail(mailConfig.FromMailAddress); err != nil {
					log.Println(err)
					continue
				}

				if err = c.Rcpt(mailing.email); err != nil {
					log.Println(err)
					continue
				}

				// Data
				w, err := c.Data()
				if err != nil {
					log.Println(err)
					continue
				}

				_, err = w.Write(msg)
				if err != nil {
					log.Println(err)
					continue
				}

				err = w.Close()
				if err != nil {
					log.Println(err)
					continue
				}

				c.Quit()
			} else {
				err := smtp.SendMail(mailConfig.TransportSMTPServer, auth, host, to, msg)
				if err != nil {
					log.Fatalf("Error on Sending mail: %s", err.Error())
				}
			}

		}

	}

}
