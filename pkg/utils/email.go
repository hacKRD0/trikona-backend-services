package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/mailjet/mailjet-apiv3-go"
)

// EmailConfig holds the Mailjet configuration
type EmailConfig struct {
	APIKey    string
	SecretKey string
	FromEmail string
	FromName  string
}

// NewEmailConfig creates a new email configuration from environment variables
func NewEmailConfig() *EmailConfig {
	return &EmailConfig{
		APIKey:    os.Getenv("MAILJET_API_KEY"),
		SecretKey: os.Getenv("MAILJET_SECRET_KEY"),
		FromEmail: os.Getenv("MAILJET_FROM_EMAIL"),
		FromName:  os.Getenv("MAILJET_FROM_NAME"),
	}
}

// SendEmailVerification sends an email verification link to the user
func SendEmailVerification(email, token string) error {
	config := NewEmailConfig()
	
	// Initialize Mailjet client
	mailjetClient := mailjet.NewMailjetClient(config.APIKey, config.SecretKey)
	if mailjetClient == nil {
		return errors.New("failed to create Mailjet client")
	}

	// Email content
	subject := "Verify your email address"
	textBody := fmt.Sprintf(`
		Hello,
		
		Please click the following link to verify your email address:
		%s/register?token=%s
		
		This link will expire in 24 hours.
		
		If you did not request this verification, please ignore this email.
	`, os.Getenv("FRONTEND_URL"), token)

	htmlBody := fmt.Sprintf(`
		<table width="100%%" cellpadding="0" cellspacing="0" border="0">
			<tr>
				<td style="padding: 20px; font-family: Arial, sans-serif; line-height: 1.6;">
					<h2 style="color: #333333; margin-bottom: 20px;">Email Verification</h2>
					<p style="margin-bottom: 20px;">Hello,</p>
					<p style="margin-bottom: 20px;">Thank you for registering with us. Please verify your email address by clicking the button below:</p>
					<table cellpadding="0" cellspacing="0" border="0" style="margin: 20px 0;">
						<tr>
							<td align="center" bgcolor="#4CAF50" style="border-radius: 5px;">
								<a href="%s/register?token=%s" target="_blank" style="padding: 10px 20px; font-size: 16px; color: #ffffff; text-decoration: none; display: inline-block;">Verify Email Address</a>
							</td>
						</tr>
					</table>
					<p style="margin-bottom: 20px;">Or copy and paste this link into your browser:</p>
					<p style="margin-bottom: 20px; word-break: break-all;">%s/register?token=%s</p>
					<p style="margin-bottom: 20px; color: #666666; font-size: 14px;">This link will expire in 24 hours.</p>
					<p style="margin-bottom: 20px; color: #666666; font-size: 14px;">If you did not request this verification, please ignore this email.</p>
				</td>
			</tr>
		</table>
	`, os.Getenv("FRONTEND_URL"), token, os.Getenv("FRONTEND_URL"), token)

	// Create email message
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: config.FromEmail,
				Name:  config.FromName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: email,
				},
			},
			Subject:  subject,
			TextPart: textBody,
			HTMLPart: htmlBody,
		},
	}

	// Send the email
	_, err := mailjetClient.SendMailV31(&mailjet.MessagesV31{Info: messagesInfo})
	return err
}

// SendPasswordResetEmail sends a password reset link to the user
func SendPasswordResetEmail(email, firstName, token string) error {
	config := NewEmailConfig()
	
	// Initialize Mailjet client
	mailjetClient := mailjet.NewMailjetClient(config.APIKey, config.SecretKey)

	// Email content
	subject := "Reset your password"
	textBody := fmt.Sprintf(`
		Hello %s,
		
		You have requested to reset your password. Please click the following link to reset it:
		%s/reset-password?token=%s
		
		This link will expire in 3 hours.
		
		If you did not request this password reset, please ignore this email.
	`, firstName, os.Getenv("FRONTEND_URL"), token)

	htmlBody := fmt.Sprintf(`
		<table width="100%%" cellpadding="0" cellspacing="0" border="0">
			<tr>
				<td style="padding: 20px; font-family: Arial, sans-serif; line-height: 1.6;">
					<h2 style="color: #333333; margin-bottom: 20px;">Password Reset Request</h2>
					<p style="margin-bottom: 20px;">Hello %s,</p>
					<p style="margin-bottom: 20px;">We received a request to reset your password. Click the button below to create a new password:</p>
					<table cellpadding="0" cellspacing="0" border="0" style="margin: 20px 0;">
						<tr>
							<td align="center" bgcolor="#2196F3" style="border-radius: 5px;">
								<a href="%s/reset-password?token=%s" target="_blank" style="padding: 10px 20px; font-size: 16px; color: #ffffff; text-decoration: none; display: inline-block;">Reset Password</a>
							</td>
						</tr>
					</table>
					<p style="margin-bottom: 20px;">Or copy and paste this link into your browser:</p>
					<p style="margin-bottom: 20px; word-break: break-all;">%s/reset-password?token=%s</p>
					<p style="margin-bottom: 20px; color: #666666; font-size: 14px;">This link will expire in 3 hours.</p>
					<p style="margin-bottom: 20px; color: #666666; font-size: 14px;">If you did not request this password reset, please ignore this email.</p>
				</td>
			</tr>
		</table>
	`, firstName, os.Getenv("FRONTEND_URL"), token, os.Getenv("FRONTEND_URL"), token)

	// Create email message
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: config.FromEmail,
				Name:  config.FromName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: email,
				},
			},
			Subject:  subject,
			TextPart: textBody,
			HTMLPart: htmlBody,
		},
	}

	// Send the email
	_, err := mailjetClient.SendMailV31(&mailjet.MessagesV31{Info: messagesInfo})
	return err
} 