package main

import (
	"GoProject/dao"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/donor/create", dao.Create).Methods("POST")
	router.HandleFunc("/donor/getAll", dao.GetAll).Methods("GET")
	router.HandleFunc("/donor/get/{id}", dao.GetDonor).Methods("GET")
	router.HandleFunc("/donor/remove/{id}", dao.RemoveDonor).Methods("GET")
	router.HandleFunc("/donor/verify/{token}", dao.VerifyDonor).Methods("GET")
	router.HandleFunc("/donor/request",dao.RequestDonor).Methods("POST")
	router.HandleFunc("/donor/medical",dao.MedicalUpload).Methods("POST")
	log.Println("Starting DonorSpace server ...")
	if err := http.ListenAndServe(":3000", router); err != nil {
		fmt.Println("Error while listening on port", err)
	}
}
