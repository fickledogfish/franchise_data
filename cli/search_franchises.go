package cli

import (
	"errors"

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
	if self.Name == "" {
		return errors.New("Cannot search for an empty string.")
	}

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
	var services []cmd.LocationService

	if err := self.makeOsmService(
		&services,
		database,
		self.Language,
	); err != nil {
		return []cmd.LocationService{}, err
	}

	return services, nil
}

func (self searchFranchises) makeOsmService(
	list *[]cmd.LocationService,
	database cmd.LocationLoader,
	preferredLanguage string,
) error {
	osmLocations, err := database.GetSavedLocationsFrom("osm")
	if err != nil {
		return err
	}

	idsToExclude := make([]domain.LocationId, len(osmLocations))
	for index, location := range osmLocations {
		idsToExclude[index] = location.Id
	}

	*list = append(*list, locationservice.NewOsmLocationService(
		osmUserAgent,
		preferredLanguage,
		[]string{"br"},
		idsToExclude,
	))

	return nil
}
