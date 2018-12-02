package service

import (
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sendgrid/sendgrid-go"
	"log"
	"fmt"
)

var (
	client *sendgrid.Client
	fromDonorSpace *mail.Email
	baseUrl string
)

func InitMailer(config map[string]string) {
	client = sendgrid.NewSendClient(config["mailerKey"])
	baseUrl = config["url"]
	fromDonorSpace = mail.NewEmail("DonorSpace", "noReply@DonorSpace.com")
}


func Send() {
	from := mail.NewEmail("Example User", "test@example.com")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Example User", "verma.av1997@gmail.com")
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient("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")
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
