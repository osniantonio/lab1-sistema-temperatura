package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	cepPkg "github.com/osniantonio/lab1-sistema-temperatura/internal/pkg/cep"
	climaPkg "github.com/osniantonio/lab1-sistema-temperatura/internal/pkg/clima"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/temperatura/{cep}", handleClimaRequest)
	http.ListenAndServe(":8080", router)
}

func handleClimaRequest(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	cep := mux.Vars(r)["cep"]

	isValid, StatusCode, InvalidMessage, Endereco := cepPkg.BuscarEndereco(ctx, cep)
	if !isValid {
		http.Error(w, InvalidMessage, StatusCode)
		return
	}

	clima, err := climaPkg.BuscarTemperatura(Endereco.Localidade)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(clima)
	if err != nil {
		http.Error(w, "Erro ao gerar JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
