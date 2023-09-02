package cli

import (
	"example.com/franchises/cmd"
)

type clean struct {
}

func (self clean) Run(context *context) error {
	cmd := cmd.NewCleanCmd()
	return cmd.Run()
}
