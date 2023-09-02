package db

import "example.com/franchises/domain"

type LocationSaver interface {
	SaveLocation(domain.Location) error
}

type LocationLoader interface {
	GetSavedLocations() ([]domain.Location, error)
	GetSavedLocationsFrom(origin string) ([]domain.Location, error)
}
