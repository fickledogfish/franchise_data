package cmd

import "example.com/franchises/domain"

type LocationLoader interface {
	GetSavedLocations() ([]domain.Location, error)
	GetSavedLocationsFrom(origin string) ([]domain.Location, error)
}
