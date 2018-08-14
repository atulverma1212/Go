package dao

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"goProject/model"
)

type DonorDAO struct {
	Server string
	Database string
}

var db *mgo.Database

const(
	Collection = "Donor"
)

func (m *DonorDAO) Connect(){
	session,err := mgo.Dial(m.Server)
	if err!=nil{
		fmt.Println("Error while connecting DB")
		return
	}
	db = session.DB(m.Database)
}

func (m *DonorDAO) Create(student model.Donor) error {
	err := db.C(Collection).Insert(&student)
	return err
}

func (m *DonorDAO) GetAll() ([]model.Donor, error){
	var students[] model.Donor
	err := db.C(Collection).Find(bson.M{}).All(&students)
	return students,err
}

func (m *DonorDAO) GetDonor(id uint64) (model.Donor, error) {
	var student model.Donor
	err := db.C(Collection).Find(bson.M{"_id":id}).One(&student)
	return student, err
}

func (m *DonorDAO) Delete(id uint64) error {
	err := db.C(Collection).Remove(bson.M{"_id":id})
	return err
}


