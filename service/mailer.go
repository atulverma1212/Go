package service

import (
	"GoProject/model"
	"strconv"
	"log"
	"google.golang.org/api/gmail/v1"
	"encoding/base64"
	"fmt"
)

type Request struct {
	from string
	to   string
	body string
}

const (
	consMIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	regSUBJ  = "DonorSpace Registration Successful"
)

func newRequest(to string) *Request {
	return &Request{
		to: to,
	}
}

func (r *Request) sendMail(message string) bool{
	body := "To: " + r.to + "\r\nSubject: " + regSUBJ + "\r\n" + consMIME + "\r\n" + message
	msg := gmail.Message{
		Raw: base64.StdEncoding.EncodeToString([]byte(body)),
	}
	if _, err := Mailer.Users.Messages.Send("verma.av1997@gmail.com", &msg).Do(); err!= nil {
		fmt.Println("Unable to send msg. " )
		fmt.Println(err)
	}
	return true
}

func SuccessMail(donor *model.Donor) {
	r := newRequest(donor.Email)
	msg := "Dear " + donor.Name + ", You have been successfully registered as a Donor at DonorSpace." +
		"Your registration id is DSpace_" + strconv.Itoa(int(donor.Id))
	if r.sendMail(msg )!=true {
		log.Println("Failed to send email to " + donor.Name)
	}
}
