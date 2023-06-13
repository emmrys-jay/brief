package utility

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
	Token     string
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func SendMail(receiverEmail, template string, data *EmailData) error {
	from := "emmrysjonathan@gmail.com"
	password := "emmrys@newRandomGmail1013"
	host, port := "smtp.gmail.com", "465"
	serverName := host + ":" + port

	body, err := parseTemplate(template, data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	auth := smtp.PlainAuth("", from, password, host)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}
	conn, err := tls.Dial("tcp", serverName, tlsconfig)
	if err != nil {
		fmt.Println(2, err)
		return err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		fmt.Println(3, err)
	}
	client.Auth(auth)

	client.Mail(from)
	client.Rcpt(receiverEmail)

	w, err := client.Data()
	if err != nil {
		return err
	}

	w.Write([]byte("subject:" + data.Subject + "\n"))
	w.Write([]byte("Content-type: text/html; charset=\"UTF-8\"\n"))
	w.Write([]byte(body))
	w.Close()

	client.Quit()

	return nil
}

func parseTemplate(template string, data *EmailData) (string, error) {
	body := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(body, template, data)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}
