package disheshandler

import (
	"html/template"
	"net/http"
	helper "restauranteweb/areas/helper"

	"github.com/go-redis/redis"
)

// This is the template to display as part of the html template
//

// ControllerInfo is
type ControllerInfo struct {
	Name    string
	Message string
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

var mongodbvar helper.DatabaseX

// List = assemble results of API call to dish list
//
func List(httpwriter http.ResponseWriter, redisclient *redis.Client) {

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/listtemplate.html")

	// Get list of dishes (api call)
	//
	var dishlist = listdishes(redisclient)

	// Assemble the display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Dish List"

	var numberoffields = 6

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Name"
	items.FieldNames[1] = "Type"
	items.FieldNames[2] = "Price"
	items.FieldNames[3] = "GlutenFree"
	items.FieldNames[4] = "DairyFree"
	items.FieldNames[5] = "Vegetarian"

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
	}

	t.Execute(httpwriter, items)
}

// LoadDisplayForAdd is X
func LoadDisplayForAdd(httpwriter http.ResponseWriter) {

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/dishadd.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Add"

	t.Execute(httpwriter, items)

}

// Add is
func Add(httpwriter http.ResponseWriter, req *http.Request, redisclient *redis.Client) {

	dishtoadd := Dish{}

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")

	ret := DishaddAPI(redisclient, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/dishlist", 301)
	} else {
		// http.ServeFile(httpwriter, req, "templates/error.html")
		// http.PostForm("templates/error.html", url.Values{"key": {"Value"}, "id": {"123"}})

		// create new template
		t, _ := template.ParseFiles("templates/indextemplate.html", "templates/error.html")

		items := DisplayTemplate{}
		items.Info.Name = "Error"
		items.Info.Message = "Dish already registered. Press back to make changes and resubmit."

		t.Execute(httpwriter, items)

	}
	return
}

// LoadDisplayForUpdate is
func LoadDisplayForUpdate(httpwriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client) {

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
		DishItem   Dish
	}

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/dishupdate.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Add"

	items.DishItem = Dish{}
	items.DishItem.Name = dishselected[0]

	var dishfind = Dish{}
	var dishname = items.DishItem.Name

	dishfind = FindAPI(redisclient, dishname)
	items.DishItem = dishfind

	t.Execute(httpwriter, items)

	return

}

// Update dish sent
func Update(redisclient *redis.Client, httpwriter http.ResponseWriter, req *http.Request) {

	dishtoadd := Dish{}

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")

	ret := DishupdateAPI(redisclient, dishtoadd)

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
		DishItem   Dish
	}

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/dishdelete.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Delete"

	items.DishItem = Dish{}
	items.DishItem.Name = dishselected[0]

	var dishfind = Dish{}
	var dishname = items.DishItem.Name

	dishfind = FindAPI(redisclient, dishname)
	items.DishItem = dishfind

	t.Execute(httpwriter, items)

	return

}

func dishdelete(httpwriter http.ResponseWriter, req *http.Request) {

	dishtoadd := Dish{}

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")

	ret := Dishdelete(mongodbvar, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}
}

func dishdeletemultiple(httpwriter http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	// Get all selected records
	dishselected := req.Form["dishes"]

	var numrecsel = len(dishselected)

	if numrecsel <= 0 {
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}

	dishtodelete := Dish{}

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

// ----------------------------------------------------------------------
// ----------------------------------------------------------------------
// ----------------------------------------------------------------------
// ----------------------------------------------------------------------
// This is the section of methods to be deleted when it is all working
// ----------------------------------------------------------------------
// ----------------------------------------------------------------------
// ----------------------------------------------------------------------
// ----------------------------------------------------------------------

func dishupdateTBD(redisclient *redis.Client, httpwriter http.ResponseWriter, req *http.Request) {

	dishtoadd := Dish{}

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")

	ret := DishupdateAPI(redisclient, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}
}
