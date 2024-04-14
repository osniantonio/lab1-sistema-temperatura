package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	cepPkg "github.com/osniantonio/lab1-sistema-temperatura/internal/pkg/cep"
	climatePkg "github.com/osniantonio/lab1-sistema-temperatura/internal/pkg/climate"
)

func main() {
	fmt.Println("Start running app...")
	router := mux.NewRouter()
	router.HandleFunc("/temperatures/{cep}", handleClimaRequest)
	fmt.Println("endpoint /temperatures/{cep} was created.")
	http.ListenAndServe(":8080", router)
}

func handleClimaRequest(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	cep := mux.Vars(r)["cep"]

	isValid, StatusCode, InvalidMessage, Address := cepPkg.SearchAddress(ctx, cep)
	if !isValid {
		http.Error(w, InvalidMessage, StatusCode)
		return
	}

	isValid, StatusCode, InvalidMessage, Temperature := climatePkg.SearchTemperature(ctx, Address.City)
	if !isValid {
		http.Error(w, InvalidMessage, StatusCode)
		return
	}

	jsonBytes, err := json.Marshal(Temperature)
	if err != nil {
		http.Error(w, "Fail to generate the JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
