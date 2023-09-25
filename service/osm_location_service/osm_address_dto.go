package osmlocationservice

import "example.com/franchises/domain"

type osmAddressDto struct {
	PostCode string `json:"postcode"`
	Road     string `json:"road"`
	Suburb   string `json:"suburb"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
}

func (self osmAddressDto) asAddress() domain.Address {
	return domain.NewAddress(
		self.Road,
		self.City,
		self.State,
		self.Country,
		self.PostCode,
	)
}
