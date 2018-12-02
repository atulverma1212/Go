package utils

import (
	"GoProject/model"
	"encoding/json"
)

func MakeResponse(donor *model.Donor) map[string] string{
	response :=  make(map[string] string)
	response["Name"] = donor.Name
	response["Address"] = donor.Address + ", " + donor.District
	response["Blood Group"] = donor.BloodGroup
	response["Phone"] = donor.Phone
	response["Email"] = donor.Email
	return response
}

func MakeError(msg string) string {
	res := make(map[string] string)
	res["Error"] = msg
	response,_ := json.Marshal(res)
	return string(response)
}

func MakeDonationResponse(doctor, donor string) []byte {
	res := "Dear " + donor + ", Your Donation request has been received successfully. After approval from Dr. " + doctor
	res += ", DonorSpace team will contact 5 donors nearest to you."
	return []byte(res)
}