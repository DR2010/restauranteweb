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
	// ----------------------------------------------------------- Root
	Route{"Index", "GET", "/", root},
	// ----------------------------------------------------------- Error
	// Route{"errorpage", "GET", "/login", errorpage},
	// ----------------------------------------------------------- Security
	Route{"login", "GET", "/login", loginPageV4},
	Route{"login", "POST", "/login", loginPageV4},
	Route{"logout", "GET", "/logout", logoutPage},
	Route{"signup", "GET", "/signup", signupPage},
	Route{"signup", "POST", "/signup", signupPage},
	// ----------------------------------------------------------- Dishes
	Route{"dishlist", "GET", "/dishlist", dishlist},
	Route{"dishlistpictures", "GET", "/dishlistpictures", dishlistpictures},
	Route{"dishadddisplay", "POST", "/dishadddisplay", dishadddisplay},
	Route{"dishupdatedisplay", "POST", "/dishupdatedisplay", dishupdatedisplay},
	Route{"dishdeletedisplay", "POST", "/dishdeletedisplay", dishdeletedisplay},
	Route{"dishdeletemultiple", "POST", "/dishdeletemultiple", dishdeletemultiple},
	Route{"dishadd", "POST", "/dishadd", dishadd},
	Route{"dishupdate", "POST", "/dishupdate", dishupdate},
	Route{"dishdelete", "POST", "/dishdelete", dishdelete},
	Route{"showcache", "GET", "/showcache", showcache},
	Route{"errorpage", "POST", "/errorpage", errorpage},
	// ----------------------------------------------------------- Order
	Route{"orderlist", "GET", "/orderlist", orderlist},
	Route{"orderadddisplay", "POST", "/orderadddisplay", orderadddisplay},
	Route{"orderadddisplay", "GET", "/orderadddisplay", orderadddisplay},
	Route{"orderadd", "POST", "/orderadd", orderadd},
	Route{"ordersettoserving", "GET", "/ordersettoserving", ordersettoserving},
	Route{"ordersettoready", "GET", "/ordersettoready", ordersettoready},
	Route{"ordercancel", "GET", "/ordercancel", ordercancel},
	Route{"orderviewdisplay", "GET", "/orderviewdisplay", orderviewdisplay},
	// ----------------------------------------------------------- Order
	Route{"btcmarketslist", "GET", "/btcmarketslist", btcmarketslistV3},
	Route{"btcmarketshistorylist", "GET", "/btcmarketshistorylist", btclistcoinshistory},
	Route{"btcmarketshistorylistdate", "GET", "/btcmarketshistorylistdate", btclistcoinshistorydate},
	Route{"btcrecordtick", "GET", "/btcrecordtick", btcrecordtick},
	// Route{"orderviewdisplay", "POST", "/orderviewdisplay", orderviewdisplay},
	Route{"btcpreorderadddisplay", "POST", "/btcpreorderadddisplay", btcpreorderadddisplay},
	Route{"btcpreorderlist", "GET", "/btcpreorderlist", btcpreorderlist},
	Route{"btcpreorderadd", "POST", "/btcpreorderadd", btcpreorderadd},
	// -------------------------------------------------------------
	Route{"belnorthgradinglist", "GET", "/belnorthgradinglist", gradingList},
	Route{"competitionllayerlist", "GET", "/competitionllayerlist", competitionPlayerList},
	Route{"payments", "GET", "/payments", paymentList},
}
