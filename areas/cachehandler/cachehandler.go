// Package cachehandler to handle call to cache
// --------------------------------------------------------------
// .../src/restauranteweb/areas/cachehandler/cachehandler.go
// --------------------------------------------------------------
package cachehandler

import (
	"html/template"
	"net/http"

	"github.com/go-redis/redis"
)

// This is the template to display as part of the html template
//

// ControllerInfo is
type ControllerInfo struct {
	Name string
}

// Row is
type Row struct {
	Description []string
}

// DisplayTemplate is
type DisplayTemplate struct {
	Info       ControllerInfo
	FieldNames []string
	Rows       []Row
}

// List = assemble results of API call to dish list
//
func List(httpwriter http.ResponseWriter, redisclient *redis.Client) {

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/listtemplate.html")

	// Get list of dishes (api call)
	//
	var cachelist = ListEntries(redisclient)

	// Assemble the display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Cache List"

	var numberoffields = 2

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Key"
	items.FieldNames[1] = "Value"

	// Set rows to be displayed
	items.Rows = make([]Row, len(cachelist))

	for i := 0; i < len(cachelist); i++ {
		items.Rows[i] = Row{}
		items.Rows[i].Description = make([]string, numberoffields)
		items.Rows[i].Description[0] = cachelist[i].Key
		items.Rows[i].Description[1] = cachelist[i].Value
	}

	t.Execute(httpwriter, items)
}
