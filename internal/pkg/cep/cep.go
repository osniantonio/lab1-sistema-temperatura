package cep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

type Endereco struct {
	Bairro      string `json:"bairro"`
	Cep         string `json:"cep"`
	Complemento string `json:"complemento"`
	DDD         string `json:"ddd"`
	GIA         string `json:"gia"`
	IBGE        string `json:"ibge"`
	Localidade  string `json:"localidade"`
	Logradouro  string `json:"logradouro"`
	Siafi       string `json:"siafi"`
	UF          string `json:"uf"`
}

func cepValidoDigitos(cep string) bool {
	return regexp.MustCompile(`^\d{8}$`).MatchString(cep)
}

func cepFormatoValido(cep string) bool {
	return regexp.MustCompile(`^\d{5}-\d{3}$`).MatchString(cep)
}

func cepValido(cep string) bool {
	return cepValidoDigitos(cep) || cepFormatoValido(cep)
}

func BuscarEndereco(ctx context.Context, cep string) (bool, int, string, Endereco) {
	if !cepValido(cep) {
		return false, http.StatusUnprocessableEntity, "invalid zipcode", Endereco{}
	}

	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	client := http.Client{Timeout: time.Second}
	resp, err := client.Get(url)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return false, http.StatusRequestTimeout, "request timeout", Endereco{}
		}
		return false, http.StatusInternalServerError, "internal server error", Endereco{}
	}
	defer resp.Body.Close()

	var endereco Endereco
	if err := json.NewDecoder(resp.Body).Decode(&endereco); err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return false, http.StatusNotFound, "can not find zipcode", Endereco{}
		} else {
			return false, http.StatusInternalServerError, "invalid response format", Endereco{}
		}
	}

	if len(endereco.Cep) == 0 {
		return false, http.StatusNotFound, "can not find zipcode", Endereco{}
	}

	return true, http.StatusOK, "", endereco
}
