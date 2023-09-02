package cmd

import (
	"errors"
	"io/fs"
	"os"

	"example.com/franchises/db"
)

type cleanCmd struct {
}

func NewCleanCmd() Command {
	return cleanCmd{}
}

func (self cleanCmd) Run() error {
	generatedFiles := []string{
		db.SqliteDatabaseFileName,
	}

	for _, file := range generatedFiles {
		err := os.Remove(file)
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	return nil
}
