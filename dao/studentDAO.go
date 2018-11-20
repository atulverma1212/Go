package dao

import (
	"fmt"
	"GoProject/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"GoProject/service"
	"log"
)

type DonorDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	Collection = "donor"
)

func (m *DonorDAO) Connect() {
	return
	session, err := mgo.Dial(m.Server)
	if err != nil {
		fmt.Println("Error while connecting DB")
		return
	}
	db = session.DB(m.Database)
}

func (m *DonorDAO) Create(donor model.Donor) error {
	err := db.C(Collection).Insert(&donor)
	go service.SuccessMail(&donor)
	return err
}

func (m *DonorDAO) GetAll() ([]model.Donor, error) {
	var students []model.Donor
	err := db.C(Collection).Find(bson.M{}).All(&students)
	return students, err
}

func (m *DonorDAO) GetDonor(id uint64) (model.Donor, error) {
	var donor model.Donor
	err := db.C(Collection).Find(bson.M{"_id": id}).One(&donor)
	return donor, err
}

func (m *DonorDAO) Delete(id uint64) error {
	err := db.C(Collection).Remove(bson.M{"_id": id})
	return err
}

func (m *DonorDAO) Verify(id uint64) error {
	donor,err := dao.GetDonor(id)
	if err!= nil {
		log.Printf("Failed to fetch donor with id %v", id)
		return err
	}
	donor.Verified = true
	if err = db.C(Collection).UpdateId(id, donor); err!= nil {
		log.Printf("Failed to Verify donor with id %v", id)
		return err
	}
	return nil
}
