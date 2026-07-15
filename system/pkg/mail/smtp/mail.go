package smtp

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/mail"
	"net/smtp"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"system/pkg/config/env"
	"time"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/utils"
	"gopkg.in/gomail.v2"
)

var (
	auth smtp.Auth
	// config env.ConfigEnv
)

type Info struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Request struct {
	from     string
	to       []string
	subject  string
	body     string
	info     *Info
	template string
}

func (i *Info) NewRequest(from string, to []string, subject string, template string) *Request {
	return &Request{
		from:     from,
		to:       to,
		subject:  subject,
		info:     i,
		template: template,
	}
}

func (r *Request) Push(data interface{}) error {
	err := r.ParseTemplate(r.template, data)
	if err != nil {
		return err
	}

	ok, err := r.SendMail()
	if ok == true {
		return nil
	}

	return err
}

func (r *Request) SendMail() (bool, error) {
	smtpServer := r.info.Host + ":" + strconv.Itoa(r.info.Port)
	auth = smtp.PlainAuth("", r.info.Username, r.info.Password, r.info.Host)

	messageId := utils.UUID() + `@paxel.com>`

	header := make(map[string]string)
	currentTime := time.Now().Local()
	header["Message-Id"] = messageId
	header["Date"] = currentTime.Format(time.RFC1123Z)
	header["From"] = r.from
	header["To"] = r.to[0]
	header["Subject"] = r.subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + r.body

	if err := smtp.SendMail(
		smtpServer,
		auth,
		r.from,
		r.to,
		[]byte(message),
	); err != nil {
		return false, err
	}

	return true, nil
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()

	return nil
}

func PushByGoMail(from, fromAlias, to, subject, emailHost, emailUsername, emailPassword, fileName, pathFileName string, emailPort int, data interface{}) error {
	var err error
	message := gomail.NewMessage()

	message.SetAddressHeader("From", from, fromAlias)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	t := template.Must(template.New(fileName).ParseFiles(pathFileName))
	message.AddAlternativeWriter("text/html", func(w io.Writer) error {
		return t.Execute(w, data)
	})

	dialer := gomail.NewDialer(emailHost, emailPort, emailUsername, emailPassword)

	if err = dialer.DialAndSend(message); err != nil {
		log.Error("[Error] SendEmailByGoMail : ", err)
		return err
	}

	return err
}

func PushMultipleRecipientByGoMail(from string, fromAlias string, to []string, subject string, emailHost string, emailUsername string, emailPassword string, fileName string, pathFileName string, emailPort int, data interface{}) error {
	var err error
	message := gomail.NewMessage()

	message.SetAddressHeader("From", from, fromAlias)
	message.SetHeader("To", to...)
	message.SetHeader("Subject", subject)

	t := template.Must(template.New(fileName).ParseFiles(pathFileName))
	message.AddAlternativeWriter("text/html", func(w io.Writer) error {
		return t.Execute(w, data)
	})

	dialer := gomail.NewDialer(emailHost, emailPort, emailUsername, emailPassword)

	if err = dialer.DialAndSend(message); err != nil {
		fmt.Println("[Error] SendEmailByGoMail : ", err)
		return err
	}

	return err
}

func PushMultipleRecipientByGoMailWithAttachment(from string, fromAlias string, to []string, subject string, emailHost string, emailUsername string, emailPassword string, fileName string, pathFileName string, attachmentFile string, emailPort int, data interface{}) error {
	var err error
	message := gomail.NewMessage()

	message.SetAddressHeader("From", from, fromAlias)
	message.SetHeader("To", to...)
	message.SetHeader("Subject", subject)
	message.Attach(attachmentFile)

	t := template.Must(template.New(fileName).ParseFiles(pathFileName))
	message.AddAlternativeWriter("text/html", func(w io.Writer) error {
		return t.Execute(w, data)
	})

	dialer := gomail.NewDialer(emailHost, emailPort, emailUsername, emailPassword)

	if err = dialer.DialAndSend(message); err != nil {
		fmt.Println("[Error] SendEmailByGoMail : ", err)
		return err
	}

	return err
}

func BaseSendEmailWithTemplate(config env.ConfigEnv, emailSendTo, subject, fileName string, data interface{}) error {

	// log.Println("config.Get(MAIL_FROM_ADDRESS):", config.Get("MAIL_FROM_ADDRESS"))
	// log.Println("config.Get(MAIL_HOST):", config.Get("MAIL_HOST"))
	// log.Println("config.Get(MAIL_FROM_NAME):", config.Get("MAIL_FROM_NAME"))
	// log.Println("config.Get(MAIL_USERNAME):", config.Get("MAIL_USERNAME"))
	// log.Println("config.Get(MAIL_PASSWORD):", config.Get("MAIL_PASSWORD"))
	// log.Println("config.Get(MAIL_PORT):", config.Get("MAIL_PORT"))
	// log.Println("BaseSendEmailWithTemplate")
	var err error
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	mailPort, err := strconv.Atoi(config.Get("MAIL_PORT"))
	if err != nil {
		log.Error("mailPort, err:", err.Error())
		return err
	}

	err = PushByGoMail(
		config.Get("MAIL_FROM_ADDRESS"),
		config.Get("MAIL_FROM_NAME"),
		emailSendTo,
		subject,
		config.Get("MAIL_HOST"),
		config.Get("MAIL_USERNAME"),
		config.Get("MAIL_PASSWORD"),
		fileName,
		basepath+"/template/"+fileName,
		mailPort,
		data)

	if err != nil {
		log.Error("err:", err.Error())
	}

	return err

}

/*
func BaseSendEmailMultipleRecipientWithTemplate(config env.ConfigEnv, emailSendTo []string, subject string, fileName string, data interface{}) error {

	var err error
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	mailPort, err := strconv.Atoi(config.Get("MAIL_PORT"))
	if err != nil {
		return err
	}

	err = PushMultipleRecipientByGoMail(
		config.Get("MAIL_FROM_ADDRESS"),
		config.Get("MAIL_FROM_NAME"),
		emailSendTo,
		subject,
		config.Get("MAIL_HOST"),
		config.Get("MAIL_USERNAME"),
		config.Get("MAIL_PASSWORD"),
		fileName,
		basepath+"/template/"+fileName,
		mailPort,
		data)

	return err

}

func BaseSendEmailMultipleRecipientWithTemplateAndAttachment(config env.ConfigEnv, emailSendTo []string, subject string, fileName string, attachmentFile string, data interface{}) error {

	var err error
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	mailPort, err := strconv.Atoi(config.Get("MAIL_PORT"))
	if err != nil {
		return err
	}

	err = PushMultipleRecipientByGoMailWithAttachment(
		config.Get("MAIL_FROM_ADDRESS"),
		config.Get("MAIL_FROM_NAME"),
		emailSendTo,
		subject,
		config.Get("MAIL_HOST"),
		config.Get("MAIL_USERNAME"),
		config.Get("MAIL_PASSWORD"),
		fileName,
		basepath+"/template/"+fileName,
		attachmentFile,
		mailPort,
		data)

	return err

}
*/
