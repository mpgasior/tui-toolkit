package main

import (
	"log"

	"github.com/mpgasior/tui-toolkit/_example/cpu/ui"
	"github.com/mpgasior/tui-toolkit/mvu"
)

func main() {
	app := ui.New()

	if err := mvu.Run(app); err != nil {
		log.Fatal(err)
	}
}
