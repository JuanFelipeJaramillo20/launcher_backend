package email

import (
	"fmt"
	"github.com/resend/resend-go/v2"
	"os"
	"sync"
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
