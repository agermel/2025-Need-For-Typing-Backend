package utils

import (
	"log"

	"type/config"

	"gopkg.in/gomail.v2"
)

func SendMail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.AllConfig.Email.Address)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject) // 邮件标题
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.qq.com", 465, "agermel@foxmail.com", "enmarcqbdpigcjed")

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}
