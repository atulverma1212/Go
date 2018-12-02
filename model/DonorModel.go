package model

type Donor struct {
	Id         uint64 `bson:"_id" json:"_id" xlsx:"-"`
	Name       string `bson:"name" json:"name" xlsx:"0"`
	DOB        string `bson:"dob" json:"dob" xlsx:"1"`
	Phone      string `bson:"phone" json:"phone" xlsx:"5"`
	Email      string `bson:"email" json:"email" xlsx:"6"`
	BloodGroup string `bson:"blood_group" json:"blood_group" xlsx:"2"`
	Address    string `bson:"address" json:"address" xlsx:"4"`
	District   string `bson:"district" json:"district" xlsx:"3"`
	Pincode    int    `bson:"pincode" json:"pincode" xlsx:"7"`
	Verified   string   `bson:"verified" json:"verified" xlsx:"-"`
}
