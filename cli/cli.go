package cli

import (
	"example.com/franchises/log"
	"github.com/alecthomas/kong"
)

type cli struct {
	Clean            clean            `cmd:"" aliases:"c" help:"Remove previous database file."`
	PrintFranchises  printLocations   `cmd:"" aliases:"p" help:"Print saved locations as CSV."`
	SearchFranchises searchFranchises `cmd:"" aliases:"f,franchises" help:"Search franchises."`

	Verbosity int `name:"" short:"v" type:"counter" help:"Verbosity"`
}

func Run() error {
	var c cli
	context := kong.Parse(
		&c,
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
	)

	c.setupLogger()

	return context.Run()
}

func (self cli) setupLogger() {
	switch self.Verbosity {
	case 0:
		log.SetMinLevel(log.LevelError)

	case 1:
		log.SetMinLevel(log.LevelWarning)

	case 2:
		log.SetMinLevel(log.LevelInfo)

	default:
		log.SetMinLevel(log.LevelDebug)
	}
}
