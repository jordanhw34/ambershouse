package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jordanhw34/ambershouse/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	// Anomymous Function
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}() // parenthesis are used to pass in run-time parameters, if it had params, would declare above and pass in here
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 12 * time.Second
	server.SendTimeout = 12 * time.Second

	client, err := server.Connect()
	if err != nil {
		log.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Body)
	} else {
		data, err := os.ReadFile(fmt.Sprintf("./templates-email/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}
		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Body, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(" > Email Sent! :-)")
	}
}
