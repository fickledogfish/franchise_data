package cli

import (
	"io"
	"os"

	"example.com/franchises/cmd"
	"example.com/franchises/db"
)

type printLocations struct {
	OutFile string `short:"o" help:"File to print to, will overwrite if it exists. Defaults to stdout."`
	Origin  string `short:"f" help:"Only export data from the given origin."`
}

func (self printLocations) Run() error {
	writer, err := makeWriter(self.OutFile)
	if err != nil {
		return err
	}

	database, err := db.NewSqliteDb()
	if err != nil {
		return err
	}

	cmd := cmd.NewPrintLocationsCmd(
		writer,
		database,
		self.Origin,
	)
	return cmd.Run()
}

func makeWriter(fileName string) (writer io.WriteCloser, err error) {
	if fileName == "" {
		writer = os.Stdout
	} else {
		// os.Create truncates the file if it already exists.
		writer, err = os.Create(fileName)
	}

	return
}
