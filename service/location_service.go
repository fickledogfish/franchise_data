package service

import "example.com/franchises/domain"

type LocationService interface {
	SearchLocation(query string) (domain.Paginated[domain.Location], error)
}
