package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"main.go/internal/controller"
	router "main.go/internal/delivery/http"
	"main.go/internal/repository"
	"main.go/internal/service"
)

var (
	masterRepository repository.MasterRepository = repository.NewFirestoreRepository()
	masterService    service.MasterService       = service.NewMasterService(masterRepository)
	masterController controller.MasterController = controller.NewMasterController(masterService)
	httpRouter       router.Router               = router.NewChiRouter()
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	const PORT string = ":8000"
	httpRouter.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Up and running...")
	})

	httpRouter.GET("/masters", masterController.GetMasters)
	httpRouter.POST("/masters", masterController.PostMaster)

	httpRouter.SERVE(PORT)
}
