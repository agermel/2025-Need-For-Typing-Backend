package utils

import (
	"log"
	"time"

	"type/config"

	"gopkg.in/gomail.v2"
)

func SendMail(to, subject, body string) error {
	start := time.Now()
	m := gomail.NewMessage()
	m.SetHeader("From", config.AllConfig.Email.Address)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject) // 邮件标题
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.163.com", 465, config.AllConfig.Email.Address, config.AllConfig.Email.Password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %s in %v", to, time.Since(start))
	return nil
}
