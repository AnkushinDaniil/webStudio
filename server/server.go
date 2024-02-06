package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()
	const PORT string = ":8000"
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Up and running...")
	})
	router.HandleFunc("/articles", getMasterss).Methods("GET")
	router.HandleFunc("/articles", postMaster).Methods("POST")
	log.Println("Server listening on port ", PORT)
	log.Fatalln(http.ListenAndServe(PORT, router))
}
