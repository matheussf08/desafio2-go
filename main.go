package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
}

type BrasilAPI struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"street"`
	Bairro     string `json:"neighborhood"`
	Localidade string `json:"city"`
	Uf         string `json:"state"`
}

func main() {
	channel1 := make(chan ViaCEP)
	channel2 := make(chan BrasilAPI)

	var cep string = "80250070"

	go BuscaViaCep(cep, channel1)
	go BuscaBrasilApi(cep, channel2)

	select {
		case viaCep := <-channel1:
			print("Resposta do Via CEP:\n\n")
			fmt.Printf("CEP: %s\nRua: %s\nBairro: %s\nCidade: %s-%s", viaCep.Cep, viaCep.Logradouro, viaCep.Bairro, viaCep.Localidade, viaCep.Uf)
		
		case brasilApi := <-channel2:
			print("Resposta do Brasil API:\n\n")
			fmt.Printf("CEP: %s\nRua: %s\nBairro: %s\nCidade: %s-%s", brasilApi.Cep, brasilApi.Logradouro, brasilApi.Bairro, brasilApi.Localidade, brasilApi.Uf)

		case <- time.After(time.Second * 1):
			println("Timeout!")
	}	
}

func BuscaViaCep(cep string, channel chan ViaCEP) {
	res, error := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if error != nil {
		panic(error)
	}
	defer res.Body.Close()

	body, error := io.ReadAll(res.Body)
	if error != nil {
		panic(error)
	}

	var viaCep ViaCEP
	error = json.Unmarshal(body, &viaCep)
	if error != nil {
		panic(error)
	}

	channel <- viaCep
}

func BuscaBrasilApi(cep string, channel chan BrasilAPI) {
	res, error := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if error != nil {
		panic(error)
	}
	defer res.Body.Close()

	body, error := io.ReadAll(res.Body)
	if error != nil {
		panic(error)
	}

	var cepBrasilAPI BrasilAPI
	error = json.Unmarshal(body, &cepBrasilAPI)
	if error != nil {
		panic(error)
	}

	channel <- cepBrasilAPI
}