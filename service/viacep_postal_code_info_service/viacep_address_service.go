package viaceppostalcodeinfoservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/franchises/domain"
	"example.com/franchises/log"
)

const (
	viaCepBaseUrl = "https://viacep.com.br/ws/%s/json"
)

type viaCepPostalCodeInfoService struct {
	client *http.Client
}

func NewViaCepAddressService() viaCepPostalCodeInfoService {
	return viaCepPostalCodeInfoService{
		client: &http.Client{},
	}
}

func (self viaCepPostalCodeInfoService) RecoverInfoFromPostalCode(
	postalCode string,
) (domain.Address, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf(viaCepBaseUrl, postalCode), nil)
	if err != nil {
		return domain.Address{}, err
	}

	log.Debug("Attempting to recover %s by reaching to %s", postalCode, request.URL.String())

	response, err := self.client.Do(request)
	if err != nil {
		return domain.Address{}, err
	}
	defer response.Body.Close()

	var addressDto viaCepAddressDto
	err = json.NewDecoder(response.Body).Decode(&addressDto)
	if err != nil {
		return domain.Address{}, err
	}

	log.Debug("ViaCEP returned %q", addressDto)

	return addressDto.asAddress(), nil
}
