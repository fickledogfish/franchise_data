package viaceppostalcodeinfoservice

import "example.com/franchises/domain"

type viaCepAddressDto struct {
	PostalCode   string `json:"cep"`
	Street       string `json:"logradouro"`
	Neighborhood string `json:"bairro"`
	City         string `json:"localidade"`
	State        string `json:"uf"`
}

func (self viaCepAddressDto) asAddress() domain.Address {
	return domain.NewAddress(
		self.Street,
		self.City,
		self.State,
		"Brasil",
		self.PostalCode,
	)
}
