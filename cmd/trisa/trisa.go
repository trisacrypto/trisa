package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/trisacrypto/trisa/cmd/trisa/app"
)

func main() {
	if err := app.Run(os.Stdout, os.Stderr); err != nil {
		log.Fatal(err)
	}
}
