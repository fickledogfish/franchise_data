package service

import "example.com/franchises/domain"

type osmAddressDto struct {
	PostCode string `json:"postcode"`
	Road     string `json:"road"`
	Suburb   string `json:"suburb"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
}

func (dto osmAddressDto) asAddress() domain.Address {
	return domain.NewAddress(
		dto.Road,
		dto.City,
		dto.State,
		dto.Country,
		dto.PostCode,
	)
}
