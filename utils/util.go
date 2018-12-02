package utils

import (
	"GoProject/model"
	"strconv"
)

func MakeResponse(donor *model.Donor) map[string] string{
	response :=  make(map[string] string)
	response["Name"] = donor.Name
	response["Address"] = donor.Address + ", " + donor.District
	response["Id"] = "DSpace_" + strconv.Itoa(int(donor.Id))
	response["Blood Group"] = donor.BloodGroup
	response["Phone"] = donor.Phone
	response["Email"] = donor.Email
	return response
}