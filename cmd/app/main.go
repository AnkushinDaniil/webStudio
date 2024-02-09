package main

import (
	"log"

	"github.com/spf13/viper"
	"main.go/internal/handler"
	"main.go/internal/repository"
	"main.go/internal/server"
	"main.go/internal/service"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initialazing configs: %s", err.Error())
	}

	repo := repository.NewRepository()
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	srv := server.NewServer()
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
