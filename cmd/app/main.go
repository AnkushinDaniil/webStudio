package main

import (
	"log"

	"main.go/internal/server"
)

func main() {
	srv := server.NewServer()
	if err := srv.Run("8000"); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
