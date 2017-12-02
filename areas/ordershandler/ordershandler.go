// Package ordershandler Handler for dishes web
// -----------------------------------------------------------
// .../src/restauranteweb/areas/disherhandler/ordershandler.go
// -----------------------------------------------------------
package ordershandler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
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
	Orders     []Order
}

var mongodbvar helper.DatabaseX

// List = assemble results of API call to dish list
//
func List(httpwriter http.ResponseWriter, redisclient *redis.Client) {

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/orderlisttemplate.html")

	// Get list of orders (api call)
	//
	var list = APICallList(redisclient)

	// Assemble the display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Order List"

	var numberoffields = 4

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Order ID"
	items.FieldNames[1] = "Name"
	items.FieldNames[2] = "Date"
	items.FieldNames[3] = "Mode"

	// Set rows to be displayed
	items.Rows = make([]Row, len(list))
	items.Orders = make([]Order, len(list))
	// items.RowID = make([]int, len(dishlist))

	for i := 0; i < len(list); i++ {
		items.Rows[i] = Row{}
		items.Rows[i].Description = make([]string, numberoffields)
		items.Rows[i].Description[0] = list[i].ID
		items.Rows[i].Description[1] = list[i].ClientName
		items.Rows[i].Description[2] = list[i].Date
		items.Rows[i].Description[3] = list[i].Foodeatplace

		items.Orders[i] = list[i]
	}

	t.Execute(httpwriter, items)
}

// LoadDisplayForAdd is X
func LoadDisplayForAdd(httpwriter http.ResponseWriter) {

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/order/orderadd.html")

	items := DisplayTemplate{}
	items.Info.Name = "Order Add"

	t.Execute(httpwriter, items)

}

// LoadDisplayForView is
func LoadDisplayForView(httpwriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client) {

	httprequest.ParseForm()

	// Get all selected records
	orderselected := httprequest.Form["dishes"]

	orderid := httprequest.URL.Query().Get("orderid")

	if orderid == "" {
		var numrecsel = len(orderselected)

		if numrecsel <= 0 {
			http.Redirect(httpwriter, httprequest, "/orderlist", 301)
			return
		}

		orderid = orderselected[0]
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
		OrderItem  Order
	}

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/order/orderview.html")

	items := DisplayTemplate{}
	items.Info.Name = "Order View"

	items.OrderItem = Order{}
	items.OrderItem.ID = orderid
	// items.OrderItem.ID = orderselected[0]

	var orderfind = Order{}
	var ordername = items.OrderItem.ID

	orderfind = FindAPI(redisclient, ordername)
	items.OrderItem = orderfind

	t.Execute(httpwriter, items)

	return
}

// Add is
func Add(httpwriter http.ResponseWriter, req *http.Request, redisclient *redis.Client) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)

	ret := APICallAdd(redisclient, bodybyte)

	if ret.ID != "" {

		obj := &RespAddOrder{ID: ret.ID}
		bresp, _ := json.Marshal(obj)

		fmt.Fprintf(httpwriter, string(bresp)) // write data to response

	} else {

		// create new template
		t, _ := template.ParseFiles("templates/indextemplate.html", "templates/error.html")

		items := DisplayTemplate{}
		items.Info.Name = "Error"
		items.Info.Message = "Order already registered."

		t.Execute(httpwriter, items)

	}
	return
}

// LoadDisplayForUpdate is
func LoadDisplayForUpdate(httpwriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client) {

	httprequest.ParseForm()

	// Get all selected records
	orderselected := httprequest.Form["dishes"]

	var numrecsel = len(orderselected)

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
		OrderItem  Order
	}

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/dishupdate.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Add"

	items.OrderItem = Order{}
	items.OrderItem.ID = orderselected[0]

	var objectfind = Order{}
	var orderid = items.OrderItem.ID

	objectfind = APICallFind(redisclient, orderid)
	items.OrderItem = objectfind

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
		DishItem   Order
	}

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/dishdelete.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Delete"

	items.DishItem = Order{}
	items.DishItem.ClientID = dishselected[0]

	var dishfind = Order{}
	var dishname = items.DishItem.ClientID

	dishfind = APICallFind(redisclient, dishname)
	items.DishItem = dishfind

	t.Execute(httpwriter, items)

	return

}
