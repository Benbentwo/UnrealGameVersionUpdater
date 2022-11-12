package main

import (
	"github.com/Benbentwo/go-bin-generic/app"
	"os"
)

func main() {
	if err := app.Run(nil); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
