package dao

import (
	"fmt"
	"GoProject/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"GoProject/service"
	"github.com/pkg/errors"
)

type DonorDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	CollectionDonor = "donor"
	CollectionMap = "Punjab"
)

func (m *DonorDAO) connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		fmt.Println("Error while connecting DB")
		return
	}
	db = session.DB(m.Database)
}

func (m *DonorDAO) create(donor *model.Donor) error {
	donor.Id = next(CollectionDonor)
	err := db.C(CollectionDonor).Insert(&donor)
	go service.SuccessMail(donor)
	return err
}

func (m *DonorDAO) getAll() ([]model.Donor, error) {
	var donors []model.Donor
	err := db.C(CollectionDonor).Find(bson.M{}).All(&donors)
	return donors, err
}

func (m *DonorDAO) find(id uint64) (*model.Donor, error) {
	var donor model.Donor
	err := db.C(CollectionDonor).Find(bson.M{"_id": id}).One(&donor)
	return &donor, err
}

func findNearest(distt string,bg string) ([]model.Donor,error){
	var donors []model.Donor
	err := db.C(CollectionDonor).Find(bson.M{"blood_group":bg,"district":distt}).All(&donors)
	return donors, err
}

func (m *DonorDAO) delete(id uint64) error {
	err := db.C(CollectionDonor).Remove(bson.M{"_id": id})
	return err
}

func (m *DonorDAO) verify(id uint64, token string) error {
	donor,err := dao.find(id)
	if err!= nil {
		log.Printf("Failed to fetch donor with id %v", id)
		return err
	}
	if token == donor.Verified {
		donor.Verified = "Yes"
	} else {
		log.Printf("Invalid verification token for donor with Id %v", donor.Id)
		return errors.New("Invalid verification token")
	}
	if err = db.C(CollectionDonor).UpdateId(id, donor); err!= nil {
		log.Printf("Failed to verify donor with id %v", id)
		return err
	}
	return nil
}

func next(id string) uint64 {
	result := bson.M{}
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"seq": 1}},
		ReturnNew: true,
	}
	_, err := db.C("Counters").Find(bson.M{"_id": id}).Apply(change, &result)
	if err != nil {
		log.Println("Error while incrementing the counter")
		return 0
	}
	return uint64(result["seq"].(float64))
}
