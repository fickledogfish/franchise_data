package cmd

import (
	"fmt"

	"example.com/franchises/db"
	"example.com/franchises/domain"
)

type printLocationsCmd struct {
	loader db.LocationLoader

	selectOrigin string
}

func NewPrintLocationsCmd(
	loader db.LocationLoader,
	origin string,
) Command {
	return printLocationsCmd{
		loader:       loader,
		selectOrigin: origin,
	}
}

func (self printLocationsCmd) Run() error {
	var locations []domain.Location
	var err error

	if self.selectOrigin != "" {
		locations, err = self.loader.GetSavedLocationsFrom(self.selectOrigin)
	} else {
		locations, err = self.loader.GetSavedLocations()
	}

	if err != nil {
		return err
	}

	return dump(locations)
}

func dump(locations []domain.Location) error {
	fmt.Println(
		"\"Index\"," +
			"\"Data Origin\"," +
			"\"Id\"," +
			"\"Name\"," +
			"\"Street\"," +
			"\"City\"," +
			"\"State\"," +
			"\"Country\"," +
			"\"Postal Code\"",
	)

	for index, location := range locations {
		fmt.Printf(
			"\"%d\",\"%s\",\"%d\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
			index,
			location.Origin,
			location.Id,
			location.Name,
			location.Address.Street,
			location.Address.City,
			location.Address.State,
			location.Address.Country,
			location.Address.PostalCode,
		)
	}

	return nil
}
