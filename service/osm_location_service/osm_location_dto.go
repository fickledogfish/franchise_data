package locationservice

import "example.com/franchises/domain"

type osmLocationDto struct {
	PlaceId     uint64        `json:"place_id"`
	Name        string        `json:"name"`
	DisplayName string        `json:"display_name"`
	Address     osmAddressDto `json:"address"`
}

func (self osmLocationDto) asLocation() domain.Location {
	return domain.NewLocation(
		self.PlaceId,
		OsmDataOrigin,
		self.Name,
		self.Address.asAddress(),
	)
}
