package mailer

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

type GoMailer struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

func NewGoMailer(from, host string, port int, username, password string) *GoMailer {
	return &GoMailer{
		From:     from,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (m *GoMailer) SendLoginCode(toEmail string, code string) error {
	subject := "🔐 Ваш код входа в YoBank"

	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="ru">
		<head>
			<meta charset="UTF-8">
			<title>Код входа</title>
		</head>
		<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f4f4f4; padding: 20px;">
			<div style="max-width: 480px; margin: auto; background-color: #ffffff; border-radius: 12px; padding: 24px; text-align: center; box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);">
				<h2 style="color: #0088cc; margin-bottom: 16px;">Вход в YoBank</h2>
				<p style="font-size: 16px; color: #333;">
					Ваш одноразовый код для входа:
				</p>

				<div style="margin: 24px 0;">
					<div style="display: inline-block; background-color: #f1f1f1; border-radius: 8px; padding: 14px 28px; font-size: 28px; font-weight: bold; letter-spacing: 6px; color: #000;">
						%s
					</div>
				</div>

				<p style="font-size: 14px; color: #666;">
					Код действителен в течение 10 минут.<br>
					Если вы не запрашивали код, просто проигнорируйте это письмо.
				</p>

				<hr style="margin: 32px 0; border: none; border-top: 1px solid #eee;">

				<p style="font-size: 13px; color: #999;">
					С уважением,<br>Команда YoBank
				</p>
			</div>
		</body>
		</html>
	`, code)

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.From)
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	dialer := gomail.NewDialer(m.Host, m.Port, m.Username, m.Password)

	return dialer.DialAndSend(msg)
}
