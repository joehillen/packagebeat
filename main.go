package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/joehillen/packagebeat/beater"
)

var Name = "packagebeat"
var Version = "0.0.1"

func main() {
	if err := beat.Run(Name, Version, beater.New()); err != nil {
		os.Exit(1)
	}
}
