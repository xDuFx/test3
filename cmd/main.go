package main

import (
	"encoding/json"
	"log"
	"os"
	"test3/package/api"
	"test3/package/models"
	"test3/package/repository"

	"github.com/gorilla/mux"
)
// "postgres://postgres:123@localhost:5432/postgres"
func main() {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := models.Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Println("error:", err)
	}		
	db, err := repository.New("postgres://"+ configuration.LoginBd +":"+ configuration.PassBd+"@localhost:"+ configuration.PortBd+"/"+ configuration.NameBd)
	if err != nil {
		log.Fatal(err.Error())
	}
	api := api.New(mux.NewRouter(), db)
	api.FillEndpoints()
	log.Fatal(api.ListenAndServe("localhost:" + configuration.ServerPort))
}