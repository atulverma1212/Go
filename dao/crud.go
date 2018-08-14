package dao

import (
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"gopkg.in/mgo.v2"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
	"goProject/model"
)

var dao = DonorDAO{}

func init(){
	var config map[string] string
	configData, err := ioutil.ReadFile("config.yml")
	if err!=nil {
		log.Println("Error while opening the config file")
		return
	}
	if err:= yaml.Unmarshal(configData, &config); err!=nil {
		log.Println("Error while unmarshalling the config file")
		return
	}
	dao.Server = config["server"]
	dao.Database = config["database"]
	dao.Connect()
}

func Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var donor model.Donor
	if err := json.NewDecoder(r.Body).Decode(&donor); err!=nil {
		http.Error(w, "Error while decoding request body", http.StatusBadRequest)
		return
	}
	donor.Id = next(Collection)
	fmt.Printf("%T", donor.Id)
	fmt.Printf("%T", 8)
	if err:= dao.Create(donor); err!=nil {
		log.Println("Error while saving data")
		http.Error(w, "Error while saving data", http.StatusNotAcceptable)
		return
	}
	response,_ :=json.Marshal("Data saved successfully! ")
	w.Write(response)
}

func GetAll(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	donors,err := dao.GetAll()
	if(err!=nil){
		http.Error(w, "Error while fetching data", http.StatusBadRequest)
	}
	response,_ := json.Marshal(donors)
	w.Write(response)
}

func GetDonor(w http.ResponseWriter, r *http.Request){
	if params:= mux.Vars(r); params!=nil{
		id,_ := strconv.Atoi(params["id"])
		donor, err:=dao.GetDonor(uint64(id))
		if err!=nil{
			http.Error(w, fmt.Sprintf("Error while fetching data with id: %d", params["id"]), http.StatusBadRequest)
		}
		response,_ := json.Marshal(donor)
		w.Write(response)
	} else {
		http.Error(w,"Invalid request parameters", http.StatusBadRequest)
	}
}

func RemoveDonor(w http.ResponseWriter, r *http.Request){
	if params:= mux.Vars(r); params!=nil{
		id,_ := strconv.Atoi(params["id"])
		err:=dao.Delete(uint64(id))
		if err!=nil {
			http.Error(w, fmt.Sprintf("Error while deleting data with id: %d", params["id"]), http.StatusBadRequest)
		}
		w.Write([]byte(fmt.Sprintf("Data with id: %d removed successfully", id)))
	} else {
		http.Error(w,"Invalid request parameters", http.StatusBadRequest)
	}
}

func next(id string) uint64 {
	result := bson.M{}
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"seq": 1}},
		ReturnNew: true,
	}
	_, err:= db.C("Counters").Find(bson.M{"_id": id}).Apply(change, &result)
	if err!=nil {
		log.Println("Error while incrementing the counter")
		return 0
	}
	return uint64(result["seq"].(float64))
}