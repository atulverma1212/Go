package dao

import (
	"testing"
	"GoProject/model"
	"fmt"
)

func TestInsert(t *testing.T) {
	str := "Atul verma"
	_ = len(str)
	fmt.Print(str[len(str)-1:])
}

func getDonor() model.Donor{
	return model.Donor{
		Id: next(Collection),
		Name: "Atul",
		DOB: "12-12-1997",
		Phone: "7508086142",
		BloodGroup: "B+",
		Address: "Jatinder Chownk",
		District: "Faridkot",
		Pincode: 151203,
	}
}