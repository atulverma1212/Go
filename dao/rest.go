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
	donor.Verified = getToken(&donor)[0:10]
	if err := dao.create(&donor); err != nil {
		log.Println("Error while saving data")
		http.Error(w, "Error while saving data", http.StatusNotAcceptable)
		return
	}
	response := utils.MakeResponse(&donor)
	response["Response"] = "You are on the way. Confirm your email by clicking the link sent in the Verification email sent to you."
	res, err := json.Marshal(response)
	if err!= nil {
		log.Printf("Error while marshaling json: " + err.Error())
	}
	w.Write(res)
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
		id, err:= dao.verify(keys[1],keys[0]);
		if err!=nil {
			errResp := utils.MakeError(err.Error())
			http.Error(w,errResp,http.StatusNotAcceptable)
		} else {
			resp := "DonorSpace Registration Successfull. You DonorSpace Id is " + id
			w.Write([]byte(resp))
		}
	} else {
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
	}
}

func RequestDonor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	request := make(map[string] string)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error while decoding request body", http.StatusBadRequest)
		return
	}
	log.Printf("[RequestDonor]: Request Form Received: %v", request)
	requestId := request["id"]
	doctor := request["Doctor"]
	dcontact := request["Contact"]
	if len(requestId)<7 || requestId[0:7]!="DSpace_" || len(doctor)==0 || len(dcontact)<10{
		err := utils.MakeError("Invalid Request. Please check if all the fields are valid")
		http.Error(w, string(err), http.StatusNotAcceptable)
		return
	}
	requestee,err := findDonor(requestId[7:])
	if err!=nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	go requestDonation(requestee)
	response := utils.MakeDonationResponse(doctor, requestee.Name)
	w.Write(response)
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
			service.RequestMail(&donor, requestee)
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
			go dao.create(&donor)
		}
	}
	w.Write([]byte("Your file has been processed successfully"))
}

func getToken(donor *model.Donor) string {
	payload := donor.Name + donor.District + donor.Phone + time.Now().String()
	return base64.StdEncoding.EncodeToString([]byte(payload))
}

