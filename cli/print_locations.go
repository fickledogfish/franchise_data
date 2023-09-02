package cli

import (
	"example.com/franchises/cmd"
	"example.com/franchises/db"
)

type printLocations struct {
	Origin string `arg:"origin" short:"o" default:"" help:"Name to search for."`
}

func (self printLocations) Run(context *context) error {
	database, err := db.NewSqliteDb()
	if err != nil {
		return err
	}

	cmd := cmd.NewPrintLocationsCmd(database, self.Origin)
	return cmd.Run()
}
