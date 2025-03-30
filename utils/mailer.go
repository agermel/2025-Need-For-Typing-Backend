package utils

import (
	"fmt"
	"log"
	"time"

	"type/config"

	"gopkg.in/gomail.v2"
)

func SendMail(to, subject, body, imageURL string) error {
	start := time.Now()
	m := gomail.NewMessage()
	m.SetHeader("From", config.AllConfig.Email.Address)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	// 规范整洁的 HTML 邮件内容
	htmlBody := fmt.Sprintf(`
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { padding: 20px; max-width: 600px; margin: auto; }
				.header { background: #f4f4f4; padding: 10px; text-align: center; font-size: 20px; }
				.content { padding: 20px; }
				.footer { font-size: 12px; text-align: center; margin-top: 20px; color: #888; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">%s</div>
				<div class="content">
					<p>%s</p>
					%s
				</div>
				<div class="footer">本邮件由系统自动发送，请勿回复。</div>
			</div>
		</body>
		</html>
	`, subject, body, func() string {
		if imageURL != "" {
			return fmt.Sprintf(`<img src="%s" alt="Image" style="max-width: 100%%; margin-top: 10px;"/>`, imageURL)
		}
		return ""
	}())

	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer("smtp.qq.com", 465, config.AllConfig.Email.Address, config.AllConfig.Email.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %s in %v", to, time.Since(start))
	return nil
}
