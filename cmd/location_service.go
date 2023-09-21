package cmd

import "example.com/franchises/domain"

type LocationService interface {
	SleepBetweenRequests()
	SearchLocation(query string) (domain.Paginated[domain.Location], error)
}
