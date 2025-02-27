package main

import (
	"log"

	"github.com/joostvanmeeuwen/phpvm/internal/tui"
)

func main() {
	app := tui.NewApp()
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
