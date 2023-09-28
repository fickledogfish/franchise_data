package osmlocationservice

import (
	"strings"

	"example.com/franchises/domain"
)

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
		strings.TrimSpace(self.Road),
		strings.TrimSpace(self.City),
		strings.TrimSpace(self.State),
		strings.TrimSpace(self.Country),
		strings.TrimSpace(self.PostCode),
	)
}
