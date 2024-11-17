package email

import (
	"fmt"
	"github.com/resend/resend-go/v2"
	"os"
	"sync"
	"time"
)

type EmailClient struct {
	client *resend.Client
}

var instance *EmailClient
var once sync.Once

func GetEmailClient() *EmailClient {
	once.Do(func() {
		apiKey := os.Getenv("RESEND_API_KEY")
		if apiKey == "" {
			fmt.Println("Warning: RESEND_API_KEY environment variable is not set.")
		}
		client := resend.NewClient(apiKey)
		instance = &EmailClient{client: client}
	})
	return instance
}

func (e *EmailClient) SendEmail(params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
	return e.client.Emails.Send(params)
}

func (e *EmailClient) SendPasswordResetEmail(nickName string, email string, resetLink string) error {
	requestDate := time.Now().Format("January 2, 2006 at 3:04 PM")
	fmt.Println("rendering in email: ", resetLink, nickName, email)
	body, err := RenderTemplate("password_reset/password_reset.html", map[string]string{
		"Username":    nickName,
		"ResetLink":   resetLink,
		"RequestDate": requestDate,
	})
	if err != nil {
		return fmt.Errorf("failed to render reset email template: %v", err)
	}

	params := &resend.SendEmailRequest{
		From:    "Support <support@jjar.lat>",
		To:      []string{email},
		Html:    body,
		Subject: fmt.Sprintf("Password Reset Request - %s", time.Now().Format("Jan 2, 2006 3:04 PM")),
	}

	_, err = e.client.Emails.Send(params)
	return err
}
