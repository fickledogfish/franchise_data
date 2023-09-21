package cli

import (
	"example.com/franchises/cmd"
	"example.com/franchises/db"
	"example.com/franchises/domain"
	locationservice "example.com/franchises/service/location_service"
)

const (
	osmUserAgent = "Learning go apis"
)

type searchFranchises struct {
	Name     string `arg:"" help:"Name to search for."`
	Language string `short:"l" help:"Preferred language for the results."`
}

func (self searchFranchises) Run() error {
	database, err := db.NewSqliteDb()
	if err != nil {
		return err
	}

	services, err := self.makeServices(database)
	if err != nil {
		return err
	}

	cmd := cmd.NewSearchLocationCmd(self.Name, database, services)
	return cmd.Run()
}

func (self searchFranchises) makeServices(
	database cmd.LocationLoader,
) ([]cmd.LocationService, error) {
	osmService, err := makeOsmService(database, self.Language)
	if err != nil {
		return []cmd.LocationService{}, err
	}

	return []cmd.LocationService{
		osmService,
	}, nil
}

func makeOsmService(
	database cmd.LocationLoader,
	preferredLanguage string,
) (cmd.LocationService, error) {
	osmLocations, err := database.GetSavedLocationsFrom("osm")
	if err != nil {
		return nil, err
	}

	idsToExclude := make([]domain.LocationId, len(osmLocations))
	for index, location := range osmLocations {
		idsToExclude[index] = location.Id
	}

	return locationservice.NewOsmLocationService(
		osmUserAgent,
		preferredLanguage,
		[]string{"br"},
		idsToExclude,
	), nil
}
