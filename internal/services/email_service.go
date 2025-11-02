package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strconv"

	"tabimoney/internal/models"
)

type EmailService struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	fromEmail    string
	fromName     string
}

type EmailTemplate struct {
	Subject string
	Body    string
}

type EmailData struct {
	UserName     string
	Title        string
	Message      string
	ActionURL    string
	Priority     string
	NotificationType string
	Amount       float64
	CategoryName string
	BudgetName   string
	GoalName     string
	Progress     float64
	Date         string
	Time         string
}

func NewEmailService() *EmailService {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if port == 0 {
		port = 587
	}

	return &EmailService{
		smtpHost:     os.Getenv("SMTP_HOST"),
		smtpPort:     port,
		smtpUsername: os.Getenv("SMTP_USERNAME"),
		smtpPassword: os.Getenv("SMTP_PASSWORD"),
		fromEmail:    os.Getenv("SMTP_FROM_EMAIL"),
		fromName:     os.Getenv("SMTP_FROM_NAME"),
	}
}

// SendNotificationEmail sends a notification email to user
func (s *EmailService) SendNotificationEmail(user *models.User, notification *models.Notification, data EmailData) error {
	if s.smtpHost == "" || s.smtpUsername == "" {
		log.Printf("Email service not configured, skipping email for user %d", user.ID)
		return nil
	}

	// Get email template based on notification type
	template := s.getEmailTemplate(notification.NotificationType, notification.Priority)
	
	// Prepare email data
	emailData := data
	emailData.UserName = user.FirstName
	if emailData.UserName == "" {
		emailData.UserName = user.Username
	}
	emailData.Title = notification.Title
	emailData.Message = notification.Message
	emailData.ActionURL = notification.ActionURL
	emailData.Priority = notification.Priority
	emailData.NotificationType = notification.NotificationType
	emailData.Date = notification.CreatedAt.Format("02/01/2006")
	emailData.Time = notification.CreatedAt.Format("15:04")

	// Render email content
	subject, body, err := s.renderEmailTemplate(template, emailData)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// Send email
	return s.sendEmail(user.Email, subject, body)
}

// getEmailTemplate returns appropriate template based on notification type and priority
func (s *EmailService) getEmailTemplate(notificationType, priority string) EmailTemplate {
	switch notificationType {
	case "warning":
		if priority == "urgent" || priority == "high" {
			return EmailTemplate{
				Subject: "🚨 Cảnh báo khẩn cấp từ TabiMoney",
				Body: s.getUrgentWarningTemplate(),
			}
		}
		return EmailTemplate{
			Subject: "⚠️ Cảnh báo từ TabiMoney",
			Body: s.getWarningTemplate(),
		}
	case "error":
		return EmailTemplate{
			Subject: "❌ Thông báo lỗi từ TabiMoney",
			Body: s.getErrorTemplate(),
		}
	case "success":
		return EmailTemplate{
			Subject: "✅ Thành công từ TabiMoney",
			Body: s.getSuccessTemplate(),
		}
	case "reminder":
		return EmailTemplate{
			Subject: "🔔 Nhắc nhở từ TabiMoney",
			Body: s.getReminderTemplate(),
		}
	default:
		return EmailTemplate{
			Subject: "📊 Thông báo từ TabiMoney",
			Body: s.getInfoTemplate(),
		}
	}
}

// renderEmailTemplate renders email template with data
func (s *EmailService) renderEmailTemplate(emailTemplate EmailTemplate, data EmailData) (string, string, error) {
	// Render subject
	subjectTmpl, err := template.New("email_subject").Parse(emailTemplate.Subject)
	if err != nil {
		return "", "", err
	}
	
	var subjectBuf bytes.Buffer
	if err := subjectTmpl.Execute(&subjectBuf, data); err != nil {
		return "", "", err
	}
	subject := subjectBuf.String()

	// Render body
	bodyTmpl, err := template.New("email_body").Parse(emailTemplate.Body)
	if err != nil {
		return "", "", err
	}
	
	var bodyBuf bytes.Buffer
	if err := bodyTmpl.Execute(&bodyBuf, data); err != nil {
		return "", "", err
	}
	body := bodyBuf.String()

	return subject, body, nil
}

// sendEmail sends email using SMTP
func (s *EmailService) sendEmail(to, subject, body string) error {
	// Create message
	msg := fmt.Sprintf("From: %s <%s>\r\n", s.fromName, s.fromEmail)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=UTF-8\r\n"
	msg += "\r\n"
	msg += body

	// Setup authentication
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)

	// Connect to server
	addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)
	
	// Create TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.smtpHost,
	}

	// Connect with STARTTLS (for Gmail and most SMTP servers)
	conn, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Check if server supports STARTTLS
	if ok, _ := conn.Extension("STARTTLS"); ok {
		if err = conn.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	// Authenticate
	if err := conn.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set sender and recipient
	if err := conn.Mail(s.fromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err := conn.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send message
	writer, err := conn.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer writer.Close()

	if _, err := writer.Write([]byte(msg)); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

// Email Templates

func (s *EmailService) getUrgentWarningTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 600px; margin: 0 auto; background-color: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #ff4444, #cc0000); color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; }
        .alert-box { background-color: #fff3cd; border: 1px solid #ffeaa7; border-radius: 6px; padding: 15px; margin: 20px 0; }
        .button { display: inline-block; background-color: #dc3545; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 10px 0; }
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; color: #6c757d; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚨 Cảnh báo khẩn cấp</h1>
        </div>
        <div class="content">
            <h2>Xin chào {{.UserName}}!</h2>
            <div class="alert-box">
                <h3>{{.Title}}</h3>
                <p>{{.Message}}</p>
            </div>
            {{if .ActionURL}}
            <p><a href="{{.ActionURL}}" class="button">Xem chi tiết</a></p>
            {{end}}
            <p><strong>Thời gian:</strong> {{.Date}} lúc {{.Time}}</p>
        </div>
        <div class="footer">
            <p>Đây là email tự động từ TabiMoney. Vui lòng không trả lời email này.</p>
        </div>
    </div>
</body>
</html>`
}

func (s *EmailService) getWarningTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 600px; margin: 0 auto; background-color: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #ffc107, #ff8f00); color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; }
        .warning-box { background-color: #fff3cd; border: 1px solid #ffeaa7; border-radius: 6px; padding: 15px; margin: 20px 0; }
        .button { display: inline-block; background-color: #ffc107; color: #212529; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 10px 0; }
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; color: #6c757d; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>⚠️ Cảnh báo</h1>
        </div>
        <div class="content">
            <h2>Xin chào {{.UserName}}!</h2>
            <div class="warning-box">
                <h3>{{.Title}}</h3>
                <p>{{.Message}}</p>
            </div>
            {{if .ActionURL}}
            <p><a href="{{.ActionURL}}" class="button">Xem chi tiết</a></p>
            {{end}}
            <p><strong>Thời gian:</strong> {{.Date}} lúc {{.Time}}</p>
        </div>
        <div class="footer">
            <p>Đây là email tự động từ TabiMoney. Vui lòng không trả lời email này.</p>
        </div>
    </div>
</body>
</html>`
}

func (s *EmailService) getErrorTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 600px; margin: 0 auto; background-color: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #dc3545, #c82333); color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; }
        .error-box { background-color: #f8d7da; border: 1px solid #f5c6cb; border-radius: 6px; padding: 15px; margin: 20px 0; }
        .button { display: inline-block; background-color: #dc3545; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 10px 0; }
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; color: #6c757d; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>❌ Thông báo lỗi</h1>
        </div>
        <div class="content">
            <h2>Xin chào {{.UserName}}!</h2>
            <div class="error-box">
                <h3>{{.Title}}</h3>
                <p>{{.Message}}</p>
            </div>
            {{if .ActionURL}}
            <p><a href="{{.ActionURL}}" class="button">Xem chi tiết</a></p>
            {{end}}
            <p><strong>Thời gian:</strong> {{.Date}} lúc {{.Time}}</p>
        </div>
        <div class="footer">
            <p>Đây là email tự động từ TabiMoney. Vui lòng không trả lời email này.</p>
        </div>
    </div>
</body>
</html>`
}

func (s *EmailService) getSuccessTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 600px; margin: 0 auto; background-color: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #28a745, #20c997); color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; }
        .success-box { background-color: #d4edda; border: 1px solid #c3e6cb; border-radius: 6px; padding: 15px; margin: 20px 0; }
        .button { display: inline-block; background-color: #28a745; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 10px 0; }
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; color: #6c757d; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ Thành công</h1>
        </div>
        <div class="content">
            <h2>Xin chào {{.UserName}}!</h2>
            <div class="success-box">
                <h3>{{.Title}}</h3>
                <p>{{.Message}}</p>
            </div>
            {{if .ActionURL}}
            <p><a href="{{.ActionURL}}" class="button">Xem chi tiết</a></p>
            {{end}}
            <p><strong>Thời gian:</strong> {{.Date}} lúc {{.Time}}</p>
        </div>
        <div class="footer">
            <p>Đây là email tự động từ TabiMoney. Vui lòng không trả lời email này.</p>
        </div>
    </div>
</body>
</html>`
}

func (s *EmailService) getReminderTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 600px; margin: 0 auto; background-color: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #17a2b8, #138496); color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; }
        .reminder-box { background-color: #d1ecf1; border: 1px solid #bee5eb; border-radius: 6px; padding: 15px; margin: 20px 0; }
        .button { display: inline-block; background-color: #17a2b8; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 10px 0; }
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; color: #6c757d; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔔 Nhắc nhở</h1>
        </div>
        <div class="content">
            <h2>Xin chào {{.UserName}}!</h2>
            <div class="reminder-box">
                <h3>{{.Title}}</h3>
                <p>{{.Message}}</p>
            </div>
            {{if .ActionURL}}
            <p><a href="{{.ActionURL}}" class="button">Xem chi tiết</a></p>
            {{end}}
            <p><strong>Thời gian:</strong> {{.Date}} lúc {{.Time}}</p>
        </div>
        <div class="footer">
            <p>Đây là email tự động từ TabiMoney. Vui lòng không trả lời email này.</p>
        </div>
    </div>
</body>
</html>`
}

func (s *EmailService) getInfoTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 600px; margin: 0 auto; background-color: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #6f42c1, #5a32a3); color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; }
        .info-box { background-color: #e2e3e5; border: 1px solid #d6d8db; border-radius: 6px; padding: 15px; margin: 20px 0; }
        .button { display: inline-block; background-color: #6f42c1; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 10px 0; }
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; color: #6c757d; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📊 Thông báo</h1>
        </div>
        <div class="content">
            <h2>Xin chào {{.UserName}}!</h2>
            <div class="info-box">
                <h3>{{.Title}}</h3>
                <p>{{.Message}}</p>
            </div>
            {{if .ActionURL}}
            <p><a href="{{.ActionURL}}" class="button">Xem chi tiết</a></p>
            {{end}}
            <p><strong>Thời gian:</strong> {{.Date}} lúc {{.Time}}</p>
        </div>
        <div class="footer">
            <p>Đây là email tự động từ TabiMoney. Vui lòng không trả lời email này.</p>
        </div>
    </div>
</body>
</html>`
}

