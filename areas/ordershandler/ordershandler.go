// Package ordershandler Handler for dishes web
// -----------------------------------------------------------
// .../src/restauranteweb/areas/disherhandler/ordershandler.go
// -----------------------------------------------------------
package ordershandler

import (
	"encoding/json"
	helper "festajuninaweb/areas/helper"
	"festajuninaweb/areas/security"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	order "restauranteapi/models"

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
	Orders     []order.Order
	OrderItem  order.Order
	Pratos     []Dish
}

var mongodbvar helper.DatabaseX

// List = assemble results of API call to dish list
//
func List(httpwriter http.ResponseWriter, redisclient *redis.Client) {

	// create new template
	t, _ := template.ParseFiles("html/index.html", "templates/order/orderlisttemplate.html")

	// Get list of orders (api call)
	//
	var list = APICallList(redisclient)

	// Assemble the display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Order List"
	items.Info.UserID = "User"

	var numberoffields = 5

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Order ID"
	items.FieldNames[1] = "Name"
	items.FieldNames[2] = "Date"
	items.FieldNames[3] = "Status"
	items.FieldNames[4] = "Mode"

	// Set rows to be displayed
	items.Rows = make([]Row, len(list))
	items.Orders = make([]order.Order, len(list))
	// items.RowID = make([]int, len(dishlist))

	for i := 0; i < len(list); i++ {
		items.Rows[i] = Row{}
		items.Rows[i].Description = make([]string, numberoffields)
		items.Rows[i].Description[0] = list[i].ID
		items.Rows[i].Description[1] = list[i].ClientName
		items.Rows[i].Description[2] = list[i].Date
		items.Rows[i].Description[3] = list[i].Status
		items.Rows[i].Description[4] = list[i].EatMode

		items.Orders[i] = list[i]
	}

	t.Execute(httpwriter, items)
}

// ListV2 = assemble results of API call to dish list
func ListV2(httpwriter http.ResponseWriter, redisclient *redis.Client, credentials helper.Credentials) {

	// create new template
	t, _ := template.ParseFiles("templates/order/indexlistrefresh.html", "templates/order/orderlisttemplate.html")

	// Get list of orders (api call)
	//
	var list = APICallListV2(redisclient, credentials)

	// Assemble the display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Order List"
	items.Info.UserID = credentials.UserID
	items.Info.UserName = credentials.UserName
	items.Info.ApplicationID = credentials.ApplicationID
	items.Info.IsAdmin = credentials.IsAdmin

	var numberoffields = 5

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Order ID"
	items.FieldNames[1] = "Name"
	items.FieldNames[2] = "Date"
	items.FieldNames[3] = "Status"
	items.FieldNames[4] = "Mode"

	// Set rows to be displayed
	items.Rows = make([]Row, len(list))
	items.Orders = make([]order.Order, len(list))
	// items.RowID = make([]int, len(dishlist))

	for i := 0; i < len(list); i++ {
		items.Rows[i] = Row{}
		items.Rows[i].Description = make([]string, numberoffields)
		items.Rows[i].Description[0] = list[i].ID
		items.Rows[i].Description[1] = list[i].ClientName
		items.Rows[i].Description[2] = list[i].Date
		items.Rows[i].Description[3] = list[i].Status
		items.Rows[i].Description[4] = list[i].EatMode

		items.Orders[i] = list[i]
	}

	t.Execute(httpwriter, items)
}

// ListV3OnlyPlaced = assemble results of API call to dish list
func ListV3OnlyPlaced(httpwriter http.ResponseWriter, redisclient *redis.Client, credentials helper.Credentials) {

	// create new template
	t, _ := template.ParseFiles("templates/order/indexlistrefresh.html", "templates/order/orderlisttemplate.html")

	// Get list of orders (api call)
	//
	var list = APICallListV2(redisclient, credentials)

	// Assemble the display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Order List"
	items.Info.UserID = credentials.UserID
	items.Info.UserName = credentials.UserName
	items.Info.ApplicationID = credentials.ApplicationID
	items.Info.IsAdmin = credentials.IsAdmin

	var numberoffields = 5

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Order ID"
	items.FieldNames[1] = "Name"
	items.FieldNames[2] = "Date"
	items.FieldNames[3] = "Status"
	items.FieldNames[4] = "Mode"

	// Set rows to be displayed
	items.Rows = make([]Row, len(list))
	items.Orders = make([]order.Order, len(list))
	// items.RowID = make([]int, len(dishlist))

	r := 0
	for i := 0; i < len(list); i++ {

		if list[i].Status == "Placed" {
			items.Rows[r] = Row{}
			items.Rows[r].Description = make([]string, numberoffields)
			items.Rows[r].Description[0] = list[i].ID
			items.Rows[r].Description[1] = list[i].ClientName
			items.Rows[r].Description[2] = list[i].Date
			items.Rows[r].Description[3] = list[i].Status
			items.Rows[r].Description[4] = list[i].EatMode

			items.Orders[r] = list[i]
			r++
		}
	}

	t.Execute(httpwriter, items)
}

// ListCompleted = assemble results of API call to dish list
func ListCompleted(httpwriter http.ResponseWriter, redisclient *redis.Client, credentials helper.Credentials) {

	// if credentials.IsAdmin != "Yes" {
	// 	return
	// }

	// create new template
	t, _ := template.ParseFiles("templates/order/indexlistrefresh.html", "templates/order/orderlisttemplate.html")

	// Get list of orders (api call)
	//
	var list = APICallListCompleted(redisclient, credentials)

	// Assemble the display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Order List"
	items.Info.UserID = credentials.UserID
	items.Info.UserName = credentials.UserName
	items.Info.ApplicationID = credentials.ApplicationID
	items.Info.IsAdmin = credentials.IsAdmin

	var numberoffields = 5

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Order ID"
	items.FieldNames[1] = "Name"
	items.FieldNames[2] = "Date"
	items.FieldNames[3] = "Status"
	items.FieldNames[4] = "Mode"

	// Set rows to be displayed
	items.Rows = make([]Row, len(list))
	items.Orders = make([]order.Order, len(list))
	// items.RowID = make([]int, len(dishlist))

	for i := 0; i < len(list); i++ {
		items.Rows[i] = Row{}
		items.Rows[i].Description = make([]string, numberoffields)
		items.Rows[i].Description[0] = list[i].ID
		items.Rows[i].Description[1] = list[i].ClientName
		items.Rows[i].Description[2] = list[i].Date
		items.Rows[i].Description[3] = list[i].Status
		items.Rows[i].Description[4] = list[i].EatMode

		items.Orders[i] = list[i]
	}

	t.Execute(httpwriter, items)
}

// ListStatus = assemble results of API call to dish list
func ListStatus(httprequest *http.Request, httpwriter http.ResponseWriter, redisclient *redis.Client, credentials helper.Credentials) {

	status := httprequest.URL.Query().Get("status")

	// create new template
	t, _ := template.ParseFiles("templates/order/indexlistrefresh.html", "templates/order/orderlisttemplate.html")

	// Get list of orders (api call)
	//
	var list = APICallListStatus(redisclient, credentials, status)

	// Assemble the display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Order List"
	items.Info.UserID = credentials.UserID
	items.Info.UserName = credentials.UserName
	items.Info.ApplicationID = credentials.ApplicationID
	items.Info.IsAdmin = credentials.IsAdmin

	var numberoffields = 5

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Order ID"
	items.FieldNames[1] = "Name"
	items.FieldNames[2] = "Date"
	items.FieldNames[3] = "Status"
	items.FieldNames[4] = "Mode"

	// Set rows to be displayed
	items.Rows = make([]Row, len(list))
	items.Orders = make([]order.Order, len(list))
	// items.RowID = make([]int, len(dishlist))

	for i := 0; i < len(list); i++ {
		items.Rows[i] = Row{}
		items.Rows[i].Description = make([]string, numberoffields)
		items.Rows[i].Description[0] = list[i].ID
		items.Rows[i].Description[1] = list[i].ClientName
		items.Rows[i].Description[2] = list[i].Date
		items.Rows[i].Description[3] = list[i].Status
		items.Rows[i].Description[4] = list[i].EatMode

		items.Orders[i] = list[i]
	}

	t.Execute(httpwriter, items)
}

// LoadDisplayForAdd is X
func LoadDisplayForAdd(httpwriter http.ResponseWriter, redisclient *redis.Client, credentials helper.Credentials) {

	// create new template
	t, _ := template.ParseFiles("templates/order/indexadd.html", "templates/order/orderadd.html")

	items := DisplayTemplate{}
	items.Info.Name = "Order Add"
	items.Info.UserID = credentials.UserID
	if credentials.UserName == "Anonymous" {
		items.Info.UserName = ""
	} else {
		items.Info.UserName = credentials.UserName
	}
	items.Info.ApplicationID = credentials.ApplicationID
	items.Info.IsAdmin = credentials.IsAdmin

	// Retrieve list of dishes by calling API to get dishes
	//
	var dishlist = Listdishes(redisclient)

	// Set rows to be displayed
	items.Pratos = make([]Dish, len(dishlist))

	for i := 0; i < len(dishlist); i++ {
		items.Pratos[i] = Dish{}
		items.Pratos[i].Name = dishlist[i].Name
		items.Pratos[i].Price = dishlist[i].Price
	}

	t.Execute(httpwriter, items)

}

// LoadDisplayForView is
func LoadDisplayForView(httpwriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client, credentials helper.Credentials) {

	httprequest.ParseForm()

	// Get all selected records
	orderselected := httprequest.Form["dishes"]

	// get the order id from the request
	orderid := httprequest.URL.Query().Get("orderid")

	if orderid == "" {
		var numrecsel = len(orderselected)

		if numrecsel <= 0 {
			http.Redirect(httpwriter, httprequest, "/orderlist", 301)
			return
		}

		orderid = orderselected[0]
	}

	// create new template
	// t, _ := template.ParseFiles("html/index.html", "templates/order/orderview.html")
	t, _ := template.ParseFiles("templates/order/indexview.html", "templates/order/orderview.html")

	items := DisplayTemplate{}
	items.Info.Name = "Order View"
	items.Info.UserID = credentials.UserID
	items.Info.UserName = credentials.UserName
	items.Info.ApplicationID = credentials.ApplicationID
	items.Info.IsAdmin = credentials.IsAdmin

	items.OrderItem = order.Order{}
	items.OrderItem.ID = orderid
	// items.OrderItem.ID = orderselected[0]

	var orderfind = order.Order{}
	var ordername = items.OrderItem.ID

	orderfind = FindAPI(redisclient, ordername)
	items.OrderItem = orderfind

	// f, err := strconv.ParseFloat("3.1415", 64)
	// sprintf
	// fmt.Printf("|%6.2f|%6.2f|\n", 1.2, 3.45)

	for x := 0; x < len(items.OrderItem.Items); x++ {

		vprice, _ := strconv.ParseFloat(items.OrderItem.Items[x].Price, 64)
		vsprice := fmt.Sprintf("%6.2f", vprice)
		items.OrderItem.Items[x].Price = vsprice

		vtotal, _ := strconv.ParseFloat(items.OrderItem.Items[x].Total, 64)
		vstotal := fmt.Sprintf("%6.2f", vtotal)
		items.OrderItem.Items[x].Total = vstotal
	}

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
		t, _ := template.ParseFiles("html/index.html", "templates/error.html")

		items := DisplayTemplate{}
		items.Info.Name = "Error"
		items.Info.Message = "Order already registered."

		t.Execute(httpwriter, items)

	}
	return
}

// AddOrderClient is designed to add order and client for anonymous
func AddOrderClient(httpwriter http.ResponseWriter, req *http.Request, redisclient *redis.Client, credentials helper.Credentials) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)

	// Well, well
	// We need to check the cookie first before we call the Place Order
	// If the cookie is present the USER ID can and should be sent
	// The API can't create a new user
	// At this point the Cookie is create by the AnonymousLogin
	// We have to break the logic
	// 1) Check cookie, get USER ID
	// 2) Send to API Call
	// 2.... Inside the API call only create a new user if "Anonymous is sent"
	// ...........
	// ...........
	// ...........
	// ...........
	// ...........  25/03/2018 --- continuar.

	ret := APICallAddOrderClient(redisclient, bodybyte)

	if ret.ID != "" {

		obj := &RespAddOrder{ID: ret.ID, ClientID: ret.ClientID}
		bresp, _ := json.Marshal(obj)

		// initialthreechar := obj.ClientID[0:3]

		// Create cookie and prevent new clients from being created
		//
		// Ainda tenho que achar e manter o user name
		// esta tudo em bytes, nao tenho acesso, posso mandar de volta da API
		// so nao sei o que armazenar no redis cache

		// if initialthreechar == "USR" {
		// Nao funcionou.
		// Permite que o Admin faca pedidos mas quando ha' logoff a coisa complica.
		//
		username := "Anonymous"
		security.AnonymousLogin(httpwriter, req, redisclient, obj.ClientID, username)

		fmt.Fprintf(httpwriter, string(bresp)) // write data to response

	} else {

		// create new template
		t, _ := template.ParseFiles("html/index.html", "templates/error.html")

		items := DisplayTemplate{}
		items.Info.Name = "Error"
		items.Info.Message = "Order already registered."

		t.Execute(httpwriter, items)

	}
	return
}

// StartServing is test
func StartServing(httpwriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client) {

	orderid := httprequest.URL.Query().Get("orderid")

	orderfind := FindAPI(redisclient, orderid)
	orderfind.Status = "Serving"
	orderfind.TimeStartServing = time.Now().String()

	orderfindbyte, _ := json.Marshal(orderfind)

	APICallUpdate(redisclient, orderfindbyte)

	return
}

// OrderisReady is test
func OrderisReady(httpwriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client) {

	orderid := httprequest.URL.Query().Get("orderid")
	orderfind := FindAPI(redisclient, orderid)
	orderfind.Status = "Ready"
	orderfind.TimeCompleted = time.Now().String()

	orderfindbyte, _ := json.Marshal(orderfind)

	APICallUpdate(redisclient, orderfindbyte)

	return
}

// OrderisCompleted is test
func OrderisCompleted(httpwriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client) {

	orderid := httprequest.URL.Query().Get("orderid")
	orderfind := FindAPI(redisclient, orderid)
	orderfind.Status = "Completed"
	orderfind.TimeCompleted = time.Now().String()

	orderfindbyte, _ := json.Marshal(orderfind)

	APICallUpdate(redisclient, orderfindbyte)

	return
}

// OrderisCancelled is to cancel order
func OrderisCancelled(httpwriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client) string {

	orderid := httprequest.URL.Query().Get("orderid")
	orderfind := FindAPI(redisclient, orderid)

	if orderfind.Status == "Placed" {
		orderfind.Status = "Cancelled"
		orderfind.TimeCancelled = time.Now().String()

		orderfindbyte, _ := json.Marshal(orderfind)

		APICallUpdate(redisclient, orderfindbyte)
		return "200 OK"
	}

	return "401 Order being served"
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
		OrderItem  order.Order
	}

	// create new template
	t, _ := template.ParseFiles("html/index.html", "templates/dishupdate.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Add"

	items.OrderItem = order.Order{}
	items.OrderItem.ID = orderselected[0]

	var objectfind = order.Order{}
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
		DishItem   order.Order
	}

	// create new template
	t, _ := template.ParseFiles("html/index.html", "templates/dishdelete.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Delete"

	items.DishItem = order.Order{}
	items.DishItem.ClientID = dishselected[0]

	var dishfind = order.Order{}
	var dishname = items.DishItem.ClientID

	dishfind = APICallFind(redisclient, dishname)
	items.DishItem = dishfind

	t.Execute(httpwriter, items)

	return

}
