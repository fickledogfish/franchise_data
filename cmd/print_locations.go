package cmd

import (
	"fmt"
	"io"

	"example.com/franchises/db"
	"example.com/franchises/domain"
)

type printLocationsCmd struct {
	writer io.Writer
	loader db.LocationLoader

	selectOrigin string
}

func NewPrintLocationsCmd(
	writer io.Writer,
	loader db.LocationLoader,
	origin string,
) Command {
	return printLocationsCmd{
		writer:       writer,
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

	return dump(self.writer, locations)
}

func dump(writer io.Writer, locations []domain.Location) error {
	if _, err := fmt.Fprintln(
		writer,
		"\"Index\","+
			"\"Data Origin\","+
			"\"Id\","+
			"\"Name\","+
			"\"Street\","+
			"\"City\","+
			"\"State\","+
			"\"Country\","+
			"\"Postal Code\"",
	); err != nil {
		return err
	}

	for index, location := range locations {
		if _, err := fmt.Fprintf(
			writer,
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
		); err != nil {
			return err
		}
	}

	return nil
}
