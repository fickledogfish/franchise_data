package cli

import (
	"example.com/franchises/cmd"
	"example.com/franchises/db"
	"example.com/franchises/domain"
	"example.com/franchises/log"
	osmService "example.com/franchises/service/osm_location_service"
)

type searchFranchises struct {
	Name     string `arg:"" help:"Name to search for."`
	Language string `short:"l" help:"Preferred language for the results."`

	OsmOptions struct {
		UserAgent string `help:"User agent to send on requests to the Nominatim API." default:"Franchise store locator"`
	} `embed:"" prefix:"osm-"`
}

func (self searchFranchises) Run() error {
	if self.Name == "" {
		return SearchFranchisesEmptySearchError()
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
	osmLocations, err := database.GetSavedLocationsFrom(osmService.OsmDataOrigin)
	if err != nil {
		return err
	}

	log.Info(
		"Preloaded %d location ID(s) for the %q source",
		len(osmLocations),
		osmService.OsmDataOrigin,
	)

	idsToExclude := make([]domain.LocationId, len(osmLocations))
	for index, location := range osmLocations {
		idsToExclude[index] = location.Id
	}

	*list = append(*list, osmService.NewOsmLocationService(
		self.OsmOptions.UserAgent,
		preferredLanguage,
		[]string{"br"},
		idsToExclude,
	))

	return nil
}
