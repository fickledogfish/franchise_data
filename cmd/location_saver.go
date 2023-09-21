package cmd

import "example.com/franchises/domain"

type LocationSaver interface {
	SaveLocation(domain.Location) error
}
