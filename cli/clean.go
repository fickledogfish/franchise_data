package cli

import (
	"errors"
	"io/fs"
	"os"

	"example.com/franchises/db"
)

type clean struct {
}

func (self clean) Run(context *context) error {
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
