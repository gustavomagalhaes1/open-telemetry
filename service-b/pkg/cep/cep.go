package cep

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type AddressResponse struct {
	CEP         string `json:"cep,omitempty"`
	Logradouro  string `json:"logradouro,omitempty"`
	Complemento string `json:"complemento,omitempty"`
	Bairro      string `json:"bairro,omitempty"`
	Localidade  string `json:"localidade,omitempty"`
	UF          string `json:"uf,omitempty"`
	Erro        bool   `json:"erro,omitempty"`
}

func GetAddressFromViaCEP(cep string) (*AddressResponse, error) {
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	resp, err := client.Get(url)
	
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var address AddressResponse
	err = json.NewDecoder(resp.Body).Decode(&address)
	if address.Erro {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &address, nil
}


