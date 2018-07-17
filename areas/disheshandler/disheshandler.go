// Package disheshandler Handler for dishes web
// -----------------------------------------------------------
// .../src/restauranteweb/areas/disherhandler/disheshandler.go
// -----------------------------------------------------------
package disheshandler

import (
	helper "festajuninaweb/areas/helper"
	"html/template"
	"net/http"
	dishes "restauranteapi/models"

	"github.com/go-redis/redis"
)

// This is the template to display as part of the html template
//

// ControllerInfo is
type ControllerInfo struct {
	Name          string
	Message       string
	UserID        string
	UserName      string
	ApplicationID string //
	IsAdmin       string //
	IsAnonymous   string //
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
	Pratos     []dishes.Dish
}

var mongodbvar helper.DatabaseX

// List = assemble results of API call to dish list
//
func List(httpwriter http.ResponseWriter, redisclient *redis.Client, credentials helper.Credentials) {

	// create new template
	t, _ := template.ParseFiles("html/index.html", "templates/listtemplate.html")

	// Get list of dishes (api call)
	//
	var dishlist = listdishes(redisclient)

	// Assemble the display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Dish List"
	items.Info.UserID = credentials.UserID
	items.Info.UserName = credentials.UserName
	items.Info.ApplicationID = credentials.ApplicationID
	items.Info.IsAdmin = credentials.IsAdmin

	var numberoffields = 8

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Name"
	items.FieldNames[1] = "Type"
	items.FieldNames[2] = "Price"
	items.FieldNames[3] = "GlutenFree"
	items.FieldNames[4] = "DairyFree"
	items.FieldNames[5] = "Vegetarian"
	items.FieldNames[6] = "Initial"
	items.FieldNames[7] = "Available"

	// Set rows to be displayed
	items.Rows = make([]Row, len(dishlist))
	// items.RowID = make([]int, len(dishlist))

	for i := 0; i < len(dishlist); i++ {
		items.Rows[i] = Row{}
		items.Rows[i].Description = make([]string, numberoffields)
		items.Rows[i].Description[0] = dishlist[i].Name
		items.Rows[i].Description[1] = dishlist[i].Type
		items.Rows[i].Description[2] = dishlist[i].Price
		items.Rows[i].Description[3] = dishlist[i].GlutenFree
		items.Rows[i].Description[4] = dishlist[i].DairyFree
		items.Rows[i].Description[5] = dishlist[i].Vegetarian
		items.Rows[i].Description[6] = dishlist[i].InitialAvailable
		items.Rows[i].Description[7] = dishlist[i].CurrentAvailable
	}

	t.Execute(httpwriter, items)
}

// ListPictures shows dishes
func ListPictures(httpwriter http.ResponseWriter, redisclient *redis.Client, credentials helper.Credentials) {

	// create new template
	t, _ := template.ParseFiles("templates/dish/dishindex.html", "templates/dish/dishavailablelist.html")

	// Get list of dishes (api call)
	//
	var dishlist = listdishes(redisclient)

	// Assemble the display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Dish List Pictures"
	items.Info.UserID = credentials.UserID
	items.Info.UserName = credentials.UserName
	items.Info.ApplicationID = credentials.ApplicationID
	items.Info.IsAdmin = credentials.IsAdmin
	items.Info.IsAnonymous = credentials.IsAnonymous

	var numberoffields = 4

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Name"
	items.FieldNames[1] = "Image"
	items.FieldNames[2] = "Description"
	items.FieldNames[3] = "Price"

	items.Pratos = make([]dishes.Dish, len(dishlist))

	for i := 0; i < len(dishlist); i++ {
		items.Pratos[i] = dishes.Dish{}
		items.Pratos[i].Name = dishlist[i].Name
		items.Pratos[i].ImageName = dishlist[i].ImageName
		items.Pratos[i].Description = dishlist[i].Description
		items.Pratos[i].Price = dishlist[i].Price
	}

	t.Execute(httpwriter, items)
}

// LoadDisplayForAdd is X
func LoadDisplayForAdd(httpwriter http.ResponseWriter) {

	// create new template
	t, _ := template.ParseFiles("html/index.html", "templates/dish/dishadd.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Add"

	t.Execute(httpwriter, items)

}

// Add is
func Add(httpwriter http.ResponseWriter, req *http.Request, redisclient *redis.Client) {

	dishtoadd := dishes.Dish{}

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")
	dishtoadd.InitialAvailable = req.FormValue("initialavailable")
	// Set to the same value as initial available quantity
	dishtoadd.CurrentAvailable = req.FormValue("initialavailable")
	dishtoadd.ImageName = req.FormValue("imagename")
	dishtoadd.Description = req.FormValue("dishdescription")
	dishtoadd.Descricao = req.FormValue("dishdescricao")

	ret := APIcallAdd(redisclient, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/dishlist", 301)
	} else {
		// http.ServeFile(httpwriter, req, "templates/error.html")
		// http.PostForm("templates/error.html", url.Values{"key": {"Value"}, "id": {"123"}})

		// create new template
		t, _ := template.ParseFiles("html/index.html", "templates/error.html")

		items := DisplayTemplate{}
		items.Info.Name = "Error"
		items.Info.Message = "Dish already registered. Press back to make changes and resubmit."

		t.Execute(httpwriter, items)

	}
	return
}

// Update dish sent
func Update(httpwriter http.ResponseWriter, req *http.Request, redisclient *redis.Client) {

	dishtoadd := dishes.Dish{}

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")
	dishtoadd.InitialAvailable = req.FormValue("dishinitialavailable")
	dishtoadd.CurrentAvailable = req.FormValue("dishcurrentavailable")
	dishtoadd.ImageName = req.FormValue("imagename")
	dishtoadd.Description = req.FormValue("dishdescription")
	dishtoadd.Descricao = req.FormValue("dishdescricao")

	ret := DishupdateAPI(redisclient, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}
}

// LoadDisplayForUpdate is
func LoadDisplayForUpdate(httpwriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client, credentials helper.Credentials) {

	httprequest.ParseForm()

	// Get all selected records
	dishselected := httprequest.Form["dishes"]

	var numrecsel = len(dishselected)

	if numrecsel <= 0 {
		http.Redirect(httpwriter, httprequest, "/dishlist", 301)
		return
	}

	type ControllerInfo struct {
		Name        string
		Message     string
		UserID      string
		Currency    string
		Application string
	}
	type Row struct {
		Description []string
	}
	type DisplayTemplate struct {
		Info       ControllerInfo
		FieldNames []string
		Rows       []Row
		DishItem   dishes.Dish
	}

	// create new template
	t, _ := template.ParseFiles("html/index.html", "templates/dish/dishupdate.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Add"
	items.Info.Currency = "SUMMARY"
	items.Info.UserID = credentials.UserID
	items.Info.Application = credentials.ApplicationID

	items.DishItem = dishes.Dish{}
	items.DishItem.Name = dishselected[0]

	var dishfind = dishes.Dish{}
	var dishname = items.DishItem.Name

	dishfind = FindAPI(redisclient, dishname)
	items.DishItem = dishfind

	t.Execute(httpwriter, items)

	return

}

// LoadDisplayForDelete is
func LoadDisplayForDelete(httpwriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client) {

	httprequest.ParseForm()

	// Get all selected records
	dishselected := httprequest.Form["dishes"]

	var numrecsel = len(dishselected)

	if numrecsel <= 0 {
		http.Redirect(httpwriter, httprequest, "/dishlist", 301)
		return
	}

	type ControllerInfo struct {
		Name    string
		Message string
	}
	type Row struct {
		Description []string
	}
	type DisplayTemplate struct {
		Info       ControllerInfo
		FieldNames []string
		Rows       []Row
		DishItem   dishes.Dish
	}

	// create new template
	t, _ := template.ParseFiles("html/index.html", "templates/dishdelete.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Delete"

	items.DishItem = dishes.Dish{}
	items.DishItem.Name = dishselected[0]

	var dishfind = dishes.Dish{}
	var dishname = items.DishItem.Name

	dishfind = FindAPI(redisclient, dishname)
	items.DishItem = dishfind

	t.Execute(httpwriter, items)

	return

}

// Delete dish sent
func Delete(redisclient *redis.Client, httpwriter http.ResponseWriter, req *http.Request) {

	dishtoadd := dishes.Dish{}

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")
	dishtoadd.InitialAvailable = req.FormValue("dishinitialavailable")
	dishtoadd.CurrentAvailable = req.FormValue("dishcurrentavailable")
	dishtoadd.ImageName = req.FormValue("imagename")
	dishtoadd.Description = req.FormValue("dishdescription")
	dishtoadd.Descricao = req.FormValue("dishdescricao")

	ret := DishdeleteAPI(redisclient, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}
}

func dishdeletedisplay(httpwriter http.ResponseWriter, req *http.Request, redisclient *redis.Client) {

	req.ParseForm()

	// Get all selected records
	dishselected := req.Form["dishes"]

	var numrecsel = len(dishselected)

	if numrecsel <= 0 {
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}

	type ControllerInfo struct {
		Name string
	}
	type Row struct {
		Description []string
	}
	type DisplayTemplate struct {
		Info       ControllerInfo
		FieldNames []string
		Rows       []Row
		DishItem   dishes.Dish
	}

	// create new template
	t, _ := template.ParseFiles("html/index.html", "templates/dish/dishdelete.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Delete"

	items.DishItem = dishes.Dish{}
	items.DishItem.Name = dishselected[0]

	var dishfind = dishes.Dish{}
	var dishname = items.DishItem.Name

	dishfind = FindAPI(redisclient, dishname)
	items.DishItem = dishfind

	t.Execute(httpwriter, items)

	return

}

func dishdelete(httpwriter http.ResponseWriter, req *http.Request) {

	dishtoadd := dishes.Dish{}

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")
	dishtoadd.InitialAvailable = req.FormValue("dishinitialavailable")
	dishtoadd.CurrentAvailable = req.FormValue("dishcurrentavailable")
	dishtoadd.ImageName = req.FormValue("imagename")
	dishtoadd.Description = req.FormValue("dishdescription")
	dishtoadd.Descricao = req.FormValue("dishdescricao")

	ret := Dishdelete(mongodbvar, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}
}

// Dishdeletemultiple is to delete multiple dishes
func Dishdeletemultiple(httpwriter http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	// Get all selected records
	dishselected := req.Form["dishes"]

	var numrecsel = len(dishselected)

	if numrecsel <= 0 {
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}

	dishtodelete := dishes.Dish{}

	ret := helper.Resultado{}

	for x := 0; x < len(dishselected); x++ {

		dishtodelete.Name = dishselected[x]

		ret = Dishdelete(mongodbvar, dishtodelete)
	}

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}

	http.Redirect(httpwriter, req, "/dishlist", 301)
	return

}
