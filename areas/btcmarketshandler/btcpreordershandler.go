// Package btcmarketshandler Handler for dishes web
// -----------------------------------------------------------
// .../src/restauranteweb/areas/disherhandler/ordershandler.go
// -----------------------------------------------------------
package btcmarketshandler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	helper "restauranteweb/areas/helper"

	"github.com/go-redis/redis"
)

// PreOrder is to be exported
type PreOrder struct {
	Currency string // Currency
	Max      string // balance
	Min      string // Cotacao
	Email    string // date time
	Buy      string // date time
	Sell     string // date time
	Date     string
	DateTime string
}

// PreOrderList = assemble results of API call to dish list
//
func PreOrderList(httpwriter http.ResponseWriter, redisclient *redis.Client, credentials helper.Credentials) {

	// create new template
	// t, _ := template.ParseFiles("templates/btcmarkets/btcindextemplate.html", "templates/btcmarkets/btcmarketslisttemplate.html")
	t, _ := template.ParseFiles("html/homepage.html", "templates/btcmarkets/btcmarketslisttemplate.html")

	// Get list of orders (api call)
	//
	var list = PreOrderAPICallList(redisclient)

	// Assemble the display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Coins"
	items.Info.Currency = "SUMMARY"
	items.Info.UserID = credentials.UserID
	items.Info.Application = credentials.ApplicationID

	var numberoffields = 5

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Order ID"
	items.FieldNames[1] = "Name"
	items.FieldNames[2] = "Date"
	items.FieldNames[3] = "Status"
	items.FieldNames[4] = "Mode"

	// Set rows to be displayed
	items.PreOrders = make([]PreOrder, len(list))

	for i := 0; i < len(list); i++ {
		items.PreOrders[i] = PreOrder{}
		items.PreOrders[i] = list[i]
	}

	t.Execute(httpwriter, items)
}

// LoadDisplayForAdd is X
func LoadDisplayForAdd(httpwriter http.ResponseWriter, redisclient *redis.Client) {

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/btcmarkets/preorderadd.html")

	items := DisplayTemplate{}
	items.Info.Name = "Pre Order Add"

	t.Execute(httpwriter, items)

}

// Add is
func BTCPreOrderAdd(httpwriter http.ResponseWriter, req *http.Request, redisclient *redis.Client) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)

	ret := PreOrderAPICallAdd(redisclient, bodybyte)

	if ret.ID != "" {

		obj := &PreOrder{Currency: "ERROR"}
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
