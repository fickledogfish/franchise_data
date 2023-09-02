package main

import (
	"example.com/franchises/cli"
	"github.com/alecthomas/kong"
)

func main() {
	var c cli.Cli
	context := kong.Parse(&c)

	err := context.Run(cli.NewContext())
	context.FatalIfErrorf(err)
}
