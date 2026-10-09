package email

import (
	"context"
	"fmt"
	"time"

	"github.com/resend/resend-go/v4"
	"github.com/spf13/viper"
)

type ResendGateway struct {
	client    *resend.Client
	fromEmail string
	appURL    string
}

func NewResendGateway(viper *viper.Viper) *ResendGateway {
	apiKey := viper.GetString("RESEND_API_KEY")
	client := resend.NewClient(apiKey)

	appPort := viper.GetString("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	return &ResendGateway{
		client:    client,
		fromEmail: viper.GetString("RESEND_SENDER_EMAIL"),
		appURL:    fmt.Sprintf("http://localhost:%s", appPort),
	}
}

func (g *ResendGateway) sendEmail(ctx context.Context, toEmail, subject, htmlBody string) error {
	if g.client == nil {
		return nil
	}

	params := &resend.SendEmailRequest{
		From:    g.fromEmail,
		To:      []string{toEmail},
		Subject: subject,
		Html:    htmlBody,
	}

	_, err := g.client.Emails.SendWithContext(ctx, params)
	return err
}

func (g *ResendGateway) SendVerificationEmail(ctx context.Context, toEmail, token string) error {
	verifyURL := fmt.Sprintf("%s/api/v1/auth/verify?token=%s", g.appURL, token)
	body := fmt.Sprintf(`
		<h2>Welcome to Nexus Video Games!</h2>
		<p>Please confirm your registration by clicking the link below:</p>
		<p>
			<a href="%s" style="
				padding: 10px 15px;
				background: #5865F2;
				color: #fff;
				text-decoration: none;
				border-radius: 4px;
			">
				Verify Account
			</a>
		</p>
	`, verifyURL)

	return g.sendEmail(ctx, toEmail, "Nexus Video Games - Verify Your Account", body)
}

func (g *ResendGateway) SendLoginNotificationEmail(ctx context.Context, toEmail, resetToken string) error {
	resetURL := fmt.Sprintf("%s/api/v1/auth/reset-password?token=%s", g.appURL, resetToken)
	body := fmt.Sprintf(`
		<h2>Security Notice: Recent Sign-In</h2>
		<p>A new login session was initiated for your Nexus Video Games account at %s.</p>
		<p>If this was not you, protect your account immediately by clicking the button below:</p>
		<p>
			<a href="%s" style="
				padding: 10px 15px;
				background: #ED4245;
				color: #fff;
				text-decoration: none;
				border-radius: 4px;
			">
				Reset Password
			</a>
		</p>
	`, time.Now().UTC().Format(time.RFC1123), resetURL)

	return g.sendEmail(ctx, toEmail, "Security Alert: New Sign-in to Nexus Video Games", body)
}

func (g *ResendGateway) SendTopupReceiptEmail(ctx context.Context, toEmail string, amount int64, invoiceID string) error {
	body := fmt.Sprintf(`
		<h2>Top-Up Confirmed</h2>
		<p>Your payment has been successfully processed.</p>
		<p><strong>Invoice ID:</strong> %s</p>
		<p><strong>Credit Added:</strong> IDR %d</p>
		<p>Your wallet balance is updated and ready for catalog game purchases.</p>
	`, invoiceID, amount)

	return g.sendEmail(ctx, toEmail, "Receipt: Nexus Video Games Wallet Top-Up", body)
}
