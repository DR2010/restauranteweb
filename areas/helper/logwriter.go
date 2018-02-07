// Package helper API calls for dishes web
// --------------------------------------------------------------
// .../src/restauranteweb/areas/helper/logwriter.go
// --------------------------------------------------------------
package helper

import (
	"log"
)

// LogFile is to be exported
type LogFile struct {
	Message string // name of the dish - this is the KEY, must be unique
}

// Write works
func Write(message string) {

	log.Println(message)
}
