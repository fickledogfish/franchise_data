package cmd

import (
	"errors"
	"io/fs"
	"os"

	"example.com/franchises/db"
	"example.com/franchises/log"
)

type cleanCmd struct {
}

func NewCleanCmd() cleanCmd {
	return cleanCmd{}
}

func (self cleanCmd) Run() error {
	generatedFiles := []string{
		db.SqliteDatabaseFileName,
	}

	for _, file := range generatedFiles {
		log.Debug("Cleaning %s", file)
		err := os.Remove(file)
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	return nil
}
