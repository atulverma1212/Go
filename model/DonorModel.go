package model

type Donor struct{
	Id 		uint64		`bson:"_id" json:"_id"`
	Name 	string	 	`bson:"name" json:"name"`
	Age 	int			`bson:"age" json:"age"`
	Phone	string		`bson:"phone" json:"phone"`
	Email	string		`bson:"email" json:"email"`
	BloodGroup	string	`bson:"blood_group" json:"blood_group"`
	Address	Address 	`bson:"address" json:"address"`
}

type Address struct {
	Street	string		`bson:"street" json:"street"`
	City	string		`bson:"city" json:"city"`
	State	string		`bson:"state" json:"state"`
	Pincode int			`bson:"pincode" json:"pincode"`
}
