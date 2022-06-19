package main

import (
	"log"
	"net/http"

	"github.com/bradishungry/golang-basic-api/internal/db"
	"github.com/bradishungry/golang-basic-api/pkg/api"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

// func corsMiddle() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
// 		if request.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 		}
// 	})
// }

func main() {
	err := db.InitDB("postgres://postgres@localhost")
	if err != nil {
		log.Print(err)
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/getCode", api.GetCodeEP).Methods("GET", "OPTIONS")
	router.HandleFunc("/getPeople", api.GetPeopleEP).Methods("GET", "OPTIONS")
	router.HandleFunc("/addPeople", api.AddPeopleEP).Methods("POST", "OPTIONS")
	router.HandleFunc("/updatePerson", api.UpdatePersonEP).Methods("PUT", "OPTIONS")
	router.HandleFunc("/deletePerson", api.DeletePersonEP).Methods("DELETE", "OPTIONS")

	//Cors took me a bit here, this is taken from a couple different stack overflow answers.
	//Originally I was just adding a response header, but preflight options needed something more
	//in depth. I imagine theres a much better way to deal with it.
	c := cors.New(cors.Options{
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedOrigins:     []string{"*"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"},
		OptionsPassthrough: true,
	})
	//Planned to use this as middleware
	//c.Handler(corsMiddle())
	handler := c.Handler(router)

	log.Println("Server started.")

	http.ListenAndServe(":8090", handler)
}
