package dao

import (
	"GoProject/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"errors"
	"GoProject/service"
	"GoProject/utils"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"encoding/base64"
	"time"
	"strings"
)

var (
	dao    = DonorDAO{}
	config map[string]string
)

func init() {
	configData, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Println("Error while opening the config file: " + err.Error())
		return
	}
	if err := yaml.Unmarshal(configData, &config); err != nil {
		log.Println("Error while unmarshalling the config file")
		return
	}
	dao.Server = config["server"]
	dao.Database = config["database"]
	service.InitMailer(config)
	//dao.Server = "localhost"
	//dao.Database = "Lifecare"
	dao.connect()
}

func Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var donor model.Donor
	if err := json.NewDecoder(r.Body).Decode(&donor); err != nil {
		http.Error(w, "Error while decoding request body", http.StatusBadRequest)
		return
	}
	if validateDonor(&donor)!=true {
		log.Println("Invalid user details")
		http.Error(w, "Invalid user details", http.StatusNotAcceptable)
		return
	}
	donor.Verified = "Yes"
	if err := dao.create(donor); err != nil {
		log.Println("Error while saving data")
		http.Error(w, "Error while saving data", http.StatusNotAcceptable)
		return
	}
	response, _ := json.Marshal("Data saved successfully! ")
	w.Write(response)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	donors, err := dao.getAll()
	if err != nil {
		http.Error(w, "Error while fetching data", http.StatusBadRequest)
	}
	response, _ := json.Marshal(donors)
	w.Write(response)
}

func GetDonor(w http.ResponseWriter, r *http.Request) {
	if params := mux.Vars(r); params != nil {
		id, _ := strconv.Atoi(params["id"])
		donor, err := dao.find(uint64(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while fetching data with id: %d", params["id"]), http.StatusBadRequest)
		}
		response, _ := json.Marshal(donor)
		w.Write(response)
	} else {
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
	}
}

func RemoveDonor(w http.ResponseWriter, r *http.Request) {
	if params := mux.Vars(r); params != nil {
		id, _ := strconv.Atoi(params["id"])
		err := dao.delete(uint64(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while deleting data with id: %d", params["id"]), http.StatusBadRequest)
		}
		w.Write([]byte(fmt.Sprintf("Data with id: %d removed successfully", id)))
	} else {
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
	}
}

func VerifyDonor(w http.ResponseWriter, r *http.Request) {
	if params := mux.Vars(r); params != nil {
		token, _ := params["token"]
		keys := strings.Split(token, "|")
		donor, err := findDonor(keys[1])
		if err!=nil {
			w.Write([]byte("Donor not found."))
			return
		}
		if donor.Verified == keys[0] {
			donor.Verified = "Yes"
			go service.SuccessMail(donor)
			w.Write([]byte("Donor " + donor.Name + "verified successfully."))
		} else {
			w.Write([]byte("Invalid verify token"))
		}
	} else {
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
	}
}

func RequestDonor(w http.ResponseWriter, r *http.Request) {
	if params := mux.Vars(r); params != nil {
		requestId := params["id"]
		requestee,err := findDonor(requestId[7:])
		if err!=nil {
			http.Error(w, err.Error(), http.StatusNoContent)
		}
		requestDonation(requestee)
		w.Write([]byte(fmt.Sprintf("Donation request sent to nearest 5 donors")))
	} else {
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
	}
}



func findDonor(requestId string) (*model.Donor, error){
	id,err := strconv.Atoi(requestId)
	if err!=nil {
		return  &model.Donor{},errors.New("Invalid request ID")
	}
	return dao.find(uint64(id))
}

func requestDonation(requestee *model.Donor) error{
	data,_ := GetMap(requestee.District)
	sortedDistt := utils.SortMap(data)
	for _,pair:= range sortedDistt {
		donorsList,err := findNearest(pair.Key,requestee.BloodGroup)
		if err!=nil {
			log.Printf("Error while fetching donors of %v district and %v blood-group.",pair.Key,requestee.BloodGroup)
			continue
		}
		for _,donor := range donorsList {
			go service.RequestMail(&donor, requestee)
		}
	}
	return nil
}

func MedicalUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, fileHeader , err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	xlFile,_ := xlsx.OpenReaderAt(file,fileHeader.Size)
	for _, sheet := range xlFile.Sheets {
		for rowNum, row := range sheet.Rows {
			if rowNum==0 {
				continue
			}
			donor := model.Donor{}
			if err=row.ReadStruct(&donor); err!=nil {
				log.Printf("[MedicalUpload]: Error while converting row to struct. Error: %v", err)
				http.Error(w,"Invalid file format",http.StatusNotAcceptable)
				return
			}
			if validateDonor(&donor)!=true {
				log.Println("Invalid user details. Donor: ", donor)
				continue
			}
			donor.Verified="Yes"
			go dao.create(donor)
		}
	}
	w.Write([]byte("Your file has been processed successfully"))
}

func getToken(donor *model.Donor) string{
	payload := donor.Name + donor.District + donor.Phone + time.Now().String()
	return base64.StdEncoding.EncodeToString([]byte(payload))
}

