package cli

import "github.com/alecthomas/kong"

type cli struct {
	Clean            clean            `cmd:"clean" aliases:"c" help:"Remove previous database file."`
	PrintFranchises  printLocations   `cmd:"print" aliases:"p" help:"Print saved locations as CSV."`
	SearchFranchises searchFranchises `cmd:"franchises" aliases:"f,franchises" help:"Search franchises."`
}

func Run() error {
	var c cli
	context := kong.Parse(&c)

	return context.Run()
}
