package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/gookit/ini/v2"
	"github.com/jordan-wright/email"
)

const (
	appVersion = "sZam5 sender 1.0"
)

var (
	addressTo  string
	subject    string
	pathAttach string
)

func init() {
	flag.StringVar(&addressTo, "to", "", "Адрес пользователя")
	flag.StringVar(&subject, "s", "sZam5", "Тема сообщения")
	flag.StringVar(&pathAttach, "f", "", "Файл архива со скриптом")
	flag.Parse()

	if addressTo == "" {
		log.Fatal("Не указан адрес пользователя")
	}
	if pathAttach == "" {
		log.Fatal("Не указан путь к файлу")
	}
}

func main() {
	fmt.Println(appVersion)

	err := ini.LoadFiles("smtpsend.ini")
	if err != nil {
		log.Fatal(err)
	}

	server := ini.String("server", "smtp.yandex.com")
	serverPort := ini.String("serverPort", "465")
	userMail := ini.String("user", "support@szam5.com")
	userName := ini.String("userName", "sZam5 Support")
	mailbody := ini.String("mailbody", "utf.txt")
	pass := ini.String("password", "")
	if pass == "" {
		log.Fatal("Пароль не может быть пустым!")
	}

	if _, err := os.Stat(pathAttach); err == os.ErrNotExist {
		log.Fatal(err)
	}

	bodyText, err := os.ReadFile(mailbody)
	if err != nil {
		log.Fatal(err)
	}

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", userName, userMail)
	e.To = []string{addressTo}
	e.Subject = subject
	e.Text = bodyText
	e.AttachFile(pathAttach)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         server,
	}

	if err := e.SendWithTLS(server+":"+serverPort, smtp.PlainAuth("", userMail, pass, server), tlsconfig); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(fmt.Sprintf(`Письмо "%s" для %s отправлено`, subject, addressTo))
	}
}
