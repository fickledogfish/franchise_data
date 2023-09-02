package main

import (
	"log"

	"example.com/franchises/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}
}
