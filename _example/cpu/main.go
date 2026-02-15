package main

import (
	"log"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/app"
	"github.com/mpgasior/tui-toolkit/mvu"
)

func main() {
	app := app.New()

	if err := mvu.Run(app); err != nil {
		log.Fatal(err)
	}
}
