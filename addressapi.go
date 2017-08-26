package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"io/ioutil"
	"os"
	"fmt"
)

type Address struct {
	ID        string   `json:"id,omitempty`
	Firstname string   `json:"firstname,omitempty`
	Lastname  string   `json:"lastname,omitempty"`
	EmailAddress string  `json:"emailaddress,omitempty"`
	PhoneNumber string  `json:"phonenumber,omitempty"`
}


var addresses []Address

func GetAddressEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range addresses {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Address{})
}
func GetAddressesEndpoint(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(addresses)
}
func CreateAddressEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var address Address
	_ = json.NewDecoder(r.Body).Decode(&address)
	address.ID = params["id"]
	addresses = append(addresses, address)
	json.NewEncoder(w).Encode(addresses)
}
func DeleteAddressEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range addresses {
		if item.ID == params["id"] {
			addresses = append(addresses[:index], addresses[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(addresses)
	}
}

// 1MB
const MAX_MEMORY = 1 * 1024 * 1024

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("uploading file")
	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
	}

	for key, value := range r.MultipartForm.Value {
		fmt.Fprintf(w, "%s:%s ", key, value)
		log.Printf("%s:%s", key, value)
	}

	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()
			path := fmt.Sprintf("files/%s", fileHeader.Filename)
			buf, _ := ioutil.ReadAll(file)
			ioutil.WriteFile(path, buf, os.ModePerm)
			fmt.Println(path)
		}
	}
}

func main() {
	router := mux.NewRouter()
	addresses = append(addresses, Address{ID: "1", Firstname: "David", Lastname: "Davey", EmailAddress: "dave@test.com", PhoneNumber: "214-555-5551" })
	addresses = append(addresses, Address{ID: "2", Firstname: "Glen", Lastname: "Gleny", EmailAddress: "glen@test.com", PhoneNumber: "214-555-5552" })
	addresses = append(addresses, Address{ID: "3", Firstname: "Daniel", Lastname: "Danny", EmailAddress: "dan@test.com", PhoneNumber: "214-555-5553" })

	router.HandleFunc("/upload", upload)
	router.Handle("/", http.FileServer(http.Dir("static")))
	
	router.HandleFunc("/addresses", GetAddressesEndpoint).Methods("GET")
	router.HandleFunc("/addresses/{id}", GetAddressEndpoint).Methods("GET")
	router.HandleFunc("/addresses/{id}", CreateAddressEndpoint).Methods("POST")
	router.HandleFunc("/addresses/{id}", DeleteAddressEndpoint).Methods("DELETE")


	log.Fatal(http.ListenAndServe(":8001", router))
}