package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"goProject/dao"
)

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/donor/create", dao.Create).Methods("POST")
	router.HandleFunc("/donor/getAll", dao.GetAll).Methods("GET")
	router.HandleFunc("/donor/get/{id}", dao.GetDonor).Methods("GET")
	router.HandleFunc("/donor/remove/{id}", dao.RemoveDonor).Methods("GET")
	if err:= http.ListenAndServe(":3000", router); err!=nil{
		fmt.Println("Error while listening on port", err)
	}
}