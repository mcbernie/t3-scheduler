package schedule

import (
	"bytes"
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
	header := headerTemplate
	header = strings.Replace(header, "%user%", u.username, -1)
	header = strings.Replace(header, "%pagename%", p.title, -1)

	// Will Replaced by template engine
	bodyTemplate := mailSection.Key("body").String()
	body := bodyTemplate
	body = strings.Replace(body, "%user%", u.username, -1)
	body = strings.Replace(body, "%pagename%", p.title, -1)
	body = strings.Replace(body, "\\n", "\n", -1)

	/*templateData := struct {
		User      string
		Pagelink  string
		Pagetitle string
	}{
		User:      u.username,
		Pagelink:  fmt.Sprintf("http://intranet/index.php?id=%d", p.pageID),
		Pagetitle: p.title,
	}

	t, tErr := template.ParseFiles(filepath.Join(s.fileadminPath, "scheduler-template.html"))
	if tErr != nil {
		log.Fatalf("Error on Template Parsing: %s", tErr.Error())
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		log.Fatalf("Error on Template Executing: %s", err.Error())
	}*/

	//body := buf.String()

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
				"MIME-version: 1.0;\r\n" +
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
				if mailConfig.TransportSMTPUsername == "" {
					c, err := smtp.Dial(mailConfig.TransportSMTPServer)
					if err != nil {
						log.Fatalf("Error on connect to SMTP Server (%s): %s", mailConfig.TransportSMTPServer, err.Error())
					}

					defer c.Close()

					c.Mail(mailConfig.FromMailAddress)
					c.Rcpt(u.mail)

					wc, err := c.Data()
					if err != nil {
						log.Fatalf("Error open Stream to SMTP Server (%s): %s", mailConfig.TransportSMTPServer, err.Error())
					}
					defer wc.Close()
					buf := bytes.NewBuffer(msg)
					if _, err = buf.WriteTo(wc); err != nil {
						log.Fatalf("Error write Stream to SMTP Server (%s): %s", mailConfig.TransportSMTPServer, err.Error())
					}

				} else {
					err := smtp.SendMail(mailConfig.TransportSMTPServer, auth, host, to, msg)
					if err != nil {
						log.Fatalf("Error on Sending mail: %s", err.Error())
					}
				}

			}

		}

	}

}
