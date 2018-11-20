package model

type Donor struct {
	Id         uint64 `bson:"_id" json:"_id"`
	Name       string `bson:"name" json:"name"`
	DOB        string   `bson:"dob" json:"dob"`
	Phone      string `bson:"phone" json:"phone"`
	Email      string `bson:"email" json:"email"`
	BloodGroup string `bson:"blood_group" json:"blood_group"`
	Address    string `bson:"address" json:"address"`
	District   string `bson:"district" json:"district"`
	Pincode    int    `bson:"pincode" json:"pincode"`
	Verified   bool   `bson:"verified" json:"verified"`
}
