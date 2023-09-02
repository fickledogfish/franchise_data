package domain

type Address struct {
	Street     string
	City       string
	State      string
	Country    string
	PostalCode string
}

func NewAddress(
	street string,
	city string,
	state string,
	country string,
	postalCode string,
) Address {
	return Address{
		Street:     street,
		City:       city,
		State:      state,
		Country:    country,
		PostalCode: postalCode,
	}
}
