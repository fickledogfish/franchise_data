package domain

type LocationId = uint64

type Location struct {
	Id      LocationId
	Origin  string
	Name    string
	Address Address
}

func NewLocation(
	id LocationId,
	origin string,
	name string,
	address Address,
) Location {
	return Location{
		Id:      id,
		Origin:  origin,
		Name:    name,
		Address: address,
	}
}
