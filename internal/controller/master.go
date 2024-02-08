package controller

import (
	"encoding/json"
	"net/http"

	"main.go/internal/entity"
	"main.go/internal/errors"
	"main.go/internal/service"
)

var masterService service.MasterService

type MasterController interface {
	GetMasters(w http.ResponseWriter, r *http.Request)
	PostMaster(w http.ResponseWriter, r *http.Request)
}

type controller struct{}

func NewMasterController(service service.MasterService) *controller {
	masterService = service
	return &controller{}
}

func (*controller) GetMasters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	masters, err := masterService.FindAll()
	if err != nil {
		handleError(w, err, "Error getting the masters")
		return
	}
	writeResponse(w, masters)
}

func (*controller) PostMaster(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var master entity.Master
	err := json.NewDecoder(r.Body).Decode(&master)
	if err != nil {
		handleError(w, err, "Error unmarshalling the request")
		return
	}
	err = masterService.Validate(&master)
	if err != nil {
		handleError(w, err, err.Error())
		return
	}
	_, err = masterService.Create(&master)
	if err != nil {
		handleError(w, err, "Error saving the master")
		return
	}
	writeResponse(w, master)
}

func handleError(w http.ResponseWriter, err error, message string) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ServiceError{Message: message})
	}
}

func writeResponse(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}
