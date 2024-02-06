package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"main.go/entity"
	"main.go/repository"
)

var repo repository.MasterRepository = repository.NewMasterRepository()

func getMasterss(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	masters, err := repo.FindAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Error getting the masters"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(masters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Error encoding the masters"}`))
		return
	}
}

func postMaster(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var master entity.Master
	err := json.NewDecoder(r.Body).Decode(&master)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Error unmarshalling the request"}`))
		return
	}

	input, err := json.Marshal(master)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Error Marshal the master"}`))
		return
	}
	hash := sha256.Sum256(input)

	master.Id = hex.EncodeToString(hash[:])
	_, err = repo.Save(&master)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Error saving the master"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(master)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Error encoding the master"}`))
		return
	}
}
