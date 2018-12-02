package service

import (
	"GoProject/model"
	"strconv"
	"log"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Request struct {
	from *mail.Email
	to   *mail.Email
	subject string
	body string
}

const (
	verSUBJ = "DonorSpace Verification Mail"
	regSUBJ  = "DonorSpace Registration Successful"
	reqSUBJ = "Blood Donation Request"
	HTMLfooter = "<br><br><br><p><hr>This is an auto-generated email. Please donot reply  to the sender. For any queries, contact team@DonorSpace</p>"
)

func newRequest(name,email,subject string) *Request {
	return &Request{
		to: mail.NewEmail(name, email),
		subject: subject,
		from: fromDonorSpace,
	}
}

func (r *Request) sendMail(msgString,msgHTML string) bool{
	message := mail.NewSingleEmail(r.from, r.subject, r.to, msgString, msgHTML)
	response, err := client.Send(message)
	if err != nil {
		log.Printf("[sendMail]: Error while sending mail. Request: %v. Error: %v. Response code: %v", r, err, response.StatusCode)
		return false
	}
	log.Println(response.Body)
	log.Printf("Mail sent successfully to %v. Response code: %v", r.to, response.StatusCode)
	return true
}

func VerificationMail(donor *model.Donor) {
	if len(donor.Email)==0 {
		log.Println("[VerificationMail]: Invalid email for " + donor.Name)
		return
	}
	r := newRequest(donor.Name, donor.Email, verSUBJ)
	link := "http://" + baseUrl + "/donor/verify/" + donor.Verified + "|" + strconv.Itoa(int(donor.Id))
	msg := "Dear " + donor.Name + ", Click the following link to verify your email at DonorSpace." + link
	htmlMsg := "<strong>" + msg + "</strong>" + HTMLfooter
	if r.sendMail(msg, htmlMsg)!=true {
		log.Println("[VerificationMail]: Failed to send verification email to " + donor.Name)
	}
}

func SuccessMail(donor *model.Donor) {
	if len(donor.Email)==0 {
		log.Println("[SuccessMail]: Invalid email for " + donor.Name)
		return
	}
	r := newRequest(donor.Name, donor.Email, regSUBJ)
	msg := "Dear " + donor.Name + ", You have been successfully registered as a Donor at DonorSpace." +
		"Your registration id is DSpace_" + strconv.Itoa(int(donor.Id))
	htmlMsg := "<strong>" + msg + "</strong>" + HTMLfooter
	if r.sendMail(msg, htmlMsg)!=true {
		log.Println("[SuccessMail]: Failed to send registration successfull email to " + donor.Name)
	}
}

func RequestMail(donor *model.Donor, req *model.Donor) {
	if len(donor.Email)==0 {
		log.Printf("Email not present for %v. Contact by phone: %v.",donor.Name,donor.Phone)
		return
	}
	r := newRequest(donor.Name, donor.Email, reqSUBJ)
	msg := "Dear " + donor.Name + ", You have been requested for blood donation by " + req.Name +
		"Kindly contact to " + req.Name + " by Phone: " + req.Phone + ". Address: " + req.Address +
		", " + req.District + ". Pincode: " + strconv.Itoa(req.Pincode)
	msgHtml := "<strong>" + msg + "<strong>" + HTMLfooter
	if r.sendMail(msg,msgHtml )!=true {
		log.Println("[RequestMail]: Failed to send request email to " + donor.Name)
	}
}


