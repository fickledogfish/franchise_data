package service

import "example.com/franchises/domain"

type osmLocationDto struct {
	PlaceId     uint64        `json:"place_id"`
	Name        string        `json:"name"`
	DisplayName string        `json:"display_name"`
	Address     osmAddressDto `json:"address"`
}

func (dto osmLocationDto) asLocation() domain.Location {
	return domain.NewLocation(
		dto.PlaceId,
		"osm",
		dto.Name,
		dto.Address.asAddress(),
	)
}
