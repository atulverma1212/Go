package service

import (
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sendgrid/sendgrid-go"
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