// Routes are defined here
// -----------------------------------------------
// .../src/restauranteweb/routes.go
// -----------------------------------------------
package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route is
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes are
type Routes []Route

// VNewRouter is
func VNewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

// XNewRouter is
func XNewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{"Index", "GET", "/", root},
	Route{"login", "GET", "/login", loginPage},
	Route{"signup", "GET", "/signup", signupPage},
	Route{"dishlist", "GET", "/dishlist", dishlist},
	Route{"dishadddisplay", "POST", "/dishadddisplay", dishadddisplay},
	Route{"dishupdatedisplay", "POST", "/dishupdatedisplay", dishupdatedisplay},
	Route{"dishdeletedisplay", "POST", "/dishdeletedisplay", dishdeletedisplay},
	Route{"dishdeletemultiple", "POST", "/dishdeletemultiple", dishdeletemultiple},
	Route{"dishadd", "POST", "/dishadd", dishadd},
	Route{"dishupdate", "POST", "/dishupdate", dishupdate},
	Route{"dishdelete", "POST", "/dishdelete", dishdelete},
	Route{"showcache", "GET", "/showcache", showcache},
	Route{"errorpage", "POST", "/errorpage", errorpage},
	Route{"orderlist", "GET", "/orderlist", orderlist},
	Route{"orderadddisplay", "POST", "/orderadddisplay", orderadddisplay},
	Route{"orderadd", "POST", "/orderadd", orderadd},
	Route{"orderviewdisplay", "GET", "/orderviewdisplay", orderviewdisplay},
	// Route{"orderviewdisplay", "POST", "/orderviewdisplay", orderviewdisplay},
}
