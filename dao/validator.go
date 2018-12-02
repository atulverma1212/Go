package dao

import (
	"GoProject/model"
	"regexp"
	"strconv"
	"strings"
	"log"
)

const (
	REGEX_NAME = "^[a-zA-Z_ ]*$"
	REGEX_PHONE = "^[0-9]{10}$"
	REGEX_PINCODE = "^[0-9]{6}$"
    REGEX_BG = "^[ABO]{1,2}[+-]$"
    REGEX_DOB= "^(0[1-9]|1[0,1,2])-(0[1-9]|[1,2][0-9]|3[0,1])-(19|20)\\d\\d$"
	REGEX_EMAIL = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

func validateDonor(donor *model.Donor) bool {
	if len(donor.Name) == 0 || len(donor.Phone) == 0 || len(donor.BloodGroup) == 0 || len(donor.District) == 0 || len(donor.DOB) == 0 || donor.Pincode == 0 || len(donor.Address)==0{
		return false
	}
	group := donor.BloodGroup
	donor.BloodGroup = strings.ToUpper(group[:len(group)-1])+group[len(group)-1:]
	if !validateName(&donor.Name) || !validatePhone(&donor.Phone) || !validatePincode(donor.Pincode) || !validateBG(&donor.BloodGroup) || !validateDOB(&donor.DOB){
		return false
	}
	if len(donor.Email)==0 || !validateEmail(&donor.Email) {
		log.Println("[validateDonor]: Invalid email for " + donor.Name + ". Removing email from entry.")
		donor.Email = ""
	}
	return true
}

func validateName(name *string) bool{
	var rgx = regexp.MustCompile(REGEX_NAME)
	return rgx.MatchString(*name)
}

func validateEmail(email *string) bool{
	rgx := regexp.MustCompile(REGEX_EMAIL)
	return rgx.MatchString(*email)
}

func validatePhone(phone *string) bool{
	rgx := regexp.MustCompile(REGEX_PHONE)
	return rgx.MatchString(*phone)
}

func validatePincode(pin int) bool{
	rgx := regexp.MustCompile(REGEX_PINCODE)
	return rgx.MatchString(strconv.Itoa(pin))
}

func validateBG(bg *string) bool{
	rgx := regexp.MustCompile(REGEX_BG)
	return rgx.MatchString(*bg)
}

func validateDOB(dob *string) bool{
	rgx := regexp.MustCompile(REGEX_DOB)
	return rgx.MatchString(*dob)
}





