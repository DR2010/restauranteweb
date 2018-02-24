// main web application program for restauranteweb
// -----------------------------------------------
// .../src/restauranteweb/restauranteweb.go
// -----------------------------------------------
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	dishes "restauranteapi/models"
	"restauranteweb/areas/belnorthhandler"
	"restauranteweb/areas/btcmarketshandler"
	cachehandler "restauranteweb/areas/cachehandler"
	disheshandler "restauranteweb/areas/disheshandler"
	helper "restauranteweb/areas/helper"
	"restauranteweb/areas/ordershandler"
	"restauranteweb/areas/security"

	"github.com/go-redis/redis"
	// _ "github.com/go-sql-driver/mysql"
)

var mongodbvar helper.DatabaseX

// var credentials helper.Credentials

var db *sql.DB
var err error
var redisclient *redis.Client

// ListOfPayments is really
var ListOfPayments []belnorthhandler.Payment

// Looks after the main routing
//
func main() {

	// db, err = sql.Open("mysql", "root:oculos18@/gufcdraws")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	panic(err.Error())
	// }

	redisclient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	loadreferencedatainredis()

	// Read variables from server
	//
	envirvar := new(helper.RestEnvVariables)
	envirvar.APIAPIServerIPAddress, _ = redisclient.Get("Web.APIServer.IPAddress").Result()
	envirvar.APIAPIServerPort, _ = redisclient.Get("Web.APIServer.Port").Result()
	envirvar.WEBServerPort, _ = redisclient.Get("WEBServerPort").Result()
	envirvar.WEBDebug, _ = redisclient.Get("Web.Debug").Result()
	envirvar.RecordCurrencyTick, _ = redisclient.Get("RecordCurrencyTick").Result()
	envirvar.RunningFromServer, _ = redisclient.Get("RunningFromServer").Result()
	envirvar.AppBelnorthEnabled, _ = redisclient.Get("AppBelnorthEnabled").Result()
	envirvar.AppBitcoinEnabled, _ = redisclient.Get("AppBitcoinEnabled").Result()
	envirvar.AppFestaJuninaEnabled, _ = redisclient.Get("AppFestaJuninaEnabled").Result()

	// btcmarketshandler.SendEmail(redisclient, "StartingSystemNow"+envirvar.RunningFromServer)

	fmt.Println(">>> Web Server: restauranteweb.exe running.")
	fmt.Println("Loading reference data in cache - Redis")

	mongodbvar.Location = "localhost"
	mongodbvar.Database = "restaurante"

	// mongodbvar.APIServer = "http://192.168.2.180:1520/"
	// mongodbvar.APIServer = "http://localhost:1520/"

	fmt.Println("Running... Web Server Listening to :" + envirvar.WEBServerPort)
	fmt.Println("API Server: " + envirvar.APIAPIServerIPAddress + " Port: " + envirvar.APIAPIServerPort)

	router := XNewRouter()

	// handle using the router mux
	//
	http.Handle("/", router) // setting router rule

	http.Handle("/html/", http.StripPrefix("/html", http.FileServer(http.Dir("./"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./js"))))
	http.Handle("/ts/", http.StripPrefix("/ts", http.FileServer(http.Dir("./ts"))))
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts", http.FileServer(http.Dir("./fonts"))))
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("./images"))))

	// test
	//  belnorthhandler.Capitalfootball(redisclient)
	// PaymentStoreMemory(redisclient)

	err := http.ListenAndServe(":1510", nil) // setting listening port
	// err := http.ListenAndServe(envirvar.WEBServerPort, nil) // setting listening port
	if err != nil {
		//using the mux router
		log.Fatal("ListenAndServe: ", err)
	}

}

func loadreferencedatainredis() {
	variable := helper.Readfileintostruct()
	err = redisclient.Set("Web.MongoDB.Database", variable.APIMongoDBDatabase, 0).Err()
	err = redisclient.Set("Web.APIServer.Port", variable.APIAPIServerPort, 0).Err()
	err = redisclient.Set("WEBServerPort", variable.WEBServerPort, 0).Err()
	err = redisclient.Set("Web.MongoDB.Location", variable.APIMongoDBLocation, 0).Err()
	err = redisclient.Set("Web.APIServer.IPAddress", variable.APIAPIServerIPAddress, 0).Err()
	err = redisclient.Set("Web.Debug", variable.WEBDebug, 0).Err()
	err = redisclient.Set("RecordCurrencyTick", variable.RecordCurrencyTick, 0).Err()
	err = redisclient.Set("RunningFromServer", variable.RunningFromServer, 0).Err()
	err = redisclient.Set("AppFestaJuninaEnabled", variable.AppFestaJuninaEnabled, 0).Err()
	err = redisclient.Set("AppBitcoinEnabled", variable.AppBitcoinEnabled, 0).Err()
	err = redisclient.Set("AppBelnorthEnabled", variable.AppBelnorthEnabled, 0).Err()
}

func root(httpwriter http.ResponseWriter, req *http.Request) {

	// http.ServeFile(httpwriter, r, "index.html")
	// return

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	helper.HomePage(httpwriter, redisclient, credentials)

	// If credentials role is RESTAURANTEUSER
	// Show Order List By User
	//
}

// ----------------------------------------------------------
// Security section
// ----------------------------------------------------------

func signupPage(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	security.SignupPage(httpresponsewriter, httprequest, redisclient)
}

func logoutPage(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	security.LogoutPage(httpresponsewriter, httprequest)
}

func loginPageV4(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	security.LoginPage(httpresponsewriter, httprequest, redisclient)

}

// ----------------------------------------------------------
// Orders section
// ----------------------------------------------------------

func orderlist(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	if security.ValidateToken(redisclient, req) == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	// If user is not ADMIN, show only users order

	ordershandler.ListV2(httpwriter, redisclient, credentials)
}

func orderadddisplay(httpwriter http.ResponseWriter, req *http.Request) {
	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}
	ordershandler.LoadDisplayForAdd(httpwriter, redisclient, credentials)
}

func orderadd(httpwriter http.ResponseWriter, req *http.Request) {
	error, _ := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}
	ordershandler.Add(httpwriter, req, redisclient)
}

func ordersettoserving(httpwriter http.ResponseWriter, req *http.Request) {
	error, _ := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}
	ordershandler.StartServing(httpwriter, req, redisclient)

	orderid := req.URL.Query().Get("orderid")
	backto := "/orderviewdisplay?orderid=" + orderid
	http.Redirect(httpwriter, req, backto, 303)
	// http.Redirect(httpwriter, req, "/orderlist", 303)
}

func ordersettoready(httpwriter http.ResponseWriter, req *http.Request) {
	error, _ := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}
	ordershandler.OrderisReady(httpwriter, req, redisclient)

	orderid := req.URL.Query().Get("orderid")
	backto := "/orderviewdisplay?orderid=" + orderid
	http.Redirect(httpwriter, req, backto, 303)
	// http.Redirect(httpwriter, req, "/orderlist", 303)
}

func ordercancel(httpwriter http.ResponseWriter, req *http.Request) {
	error, _ := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	error2 := ordershandler.OrderisCancelled(httpwriter, req, redisclient)

	if error2 == "401 Order being served" {
		http.Redirect(httpwriter, req, "/errorpage", 301)
		return

	}

	orderid := req.URL.Query().Get("orderid")
	backto := "/orderviewdisplay?orderid=" + orderid
	http.Redirect(httpwriter, req, backto, 303)
	// http.Redirect(httpwriter, req, "/orderlist", 303)
}

func ordersettocompleted(httpwriter http.ResponseWriter, req *http.Request) {
	error, _ := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	ordershandler.OrderisCompleted(httpwriter, req, redisclient)

	orderid := req.URL.Query().Get("orderid")
	backto := "/orderviewdisplay?orderid=" + orderid
	http.Redirect(httpwriter, req, backto, 303)
	// http.Redirect(httpwriter, req, "/orderlist", 303)
}

func orderviewdisplay(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}
	ordershandler.LoadDisplayForView(httpwriter, req, redisclient, credentials)
}

func orderStartServing(httpwriter http.ResponseWriter, req *http.Request) {
	error, _ := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}
	ordershandler.StartServing(httpwriter, req, redisclient)
}

// ----------------------------------------------------------
// BTC Markets section
// ----------------------------------------------------------

func btcpreorderadddisplay(httpwriter http.ResponseWriter, req *http.Request) {
	btcmarketshandler.LoadDisplayForAdd(httpwriter, redisclient)
}

func btcpreorderadd(httpwriter http.ResponseWriter, req *http.Request) {
	btcmarketshandler.BTCPreOrderAdd(httpwriter, req, redisclient)
}

func btcpreorderlist(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	btcmarketshandler.PreOrderList(httpwriter, redisclient, credentials)
}

func btcmarketslistV3(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	envirvar := new(helper.RestEnvVariables)
	envirvar.RecordCurrencyTick, _ = redisclient.Get("RecordCurrencyTick").Result()

	var listofbit = btcmarketshandler.ListV2(httpwriter, redisclient, credentials)

	if envirvar.RecordCurrencyTick == "Y" {
		btcmarketshandler.RecordTick(redisclient, listofbit, "btcmarketslistV3")
	}
}

func btclistcoinshistory(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	params := req.URL.Query()
	var currency = params.Get("currency")
	if currency == "" {
		currency = "ALL"
	}

	var rows = params.Get("rows")
	if rows == "" {
		rows = "100"
	}

	btcmarketshandler.HListHistory(httpwriter, redisclient, credentials, currency, rows)

}

func btclistcoinshistorydate(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	params := req.URL.Query()
	var currency = params.Get("currency")
	var fromDate = params.Get("fromDate")
	var toDate = params.Get("toDate")

	if fromDate == "" {
		fromDate = "2018-01-06"
	}
	if toDate == "" {
		toDate = "2018-01-07"
	}

	if currency == "" {
		currency = "ALL"
	}

	btcmarketshandler.HListHistoryDate(httpwriter, redisclient, credentials, currency, fromDate, toDate)

}

func btcrecordtick(httpwriter http.ResponseWriter, req *http.Request) {

	// Fazer o record tick aceitar um parametro para gravar de onde a rotina foi chamada
	// .... btcrecordtick?rotina=CURLubuntuAUTO
	// .... btcrecordtick?rotina=WindowsPC
	// .... btcrecordtick?rotina=WindowsPCCURL

	params := req.URL.Query()
	var rotina = params.Get("rotina")
	if rotina == "" {
		rotina = "Not sure - web test most likely"
	}

	var listofbit = btcmarketshandler.GetBalance(redisclient)
	btcmarketshandler.RecordTick(redisclient, listofbit, rotina)

	jsonval, _ := json.Marshal(listofbit)
	jsonstring := string(jsonval)

	http.Error(httpwriter, jsonstring, 200)
}

// ----------------------------------------------------------
// Belnorth section
// ----------------------------------------------------------
func gradingList(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	envirvar := new(helper.RestEnvVariables)
	envirvar.RecordCurrencyTick, _ = redisclient.Get("RecordCurrencyTick").Result()

	belnorthhandler.HListGradingPlayers(httpwriter, redisclient, credentials, ListOfPayments)

}

func competitionPlayerList(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	envirvar := new(helper.RestEnvVariables)
	envirvar.RecordCurrencyTick, _ = redisclient.Get("RecordCurrencyTick").Result()

	belnorthhandler.HListCompetitionPlayers(httpwriter, redisclient, credentials, ListOfPayments)

}

func paymentList(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	envirvar := new(helper.RestEnvVariables)
	envirvar.RecordCurrencyTick, _ = redisclient.Get("RecordCurrencyTick").Result()

	belnorthhandler.HListPayments(httpwriter, redisclient, credentials)

}

// PaymentStoreMemory stores payments in memory
func PaymentStoreMemory(redisclient *redis.Client) []belnorthhandler.Payment {

	ListOfPayments = belnorthhandler.ListPayments(redisclient)
	return ListOfPayments
}

// ----------------------------------------------------------
// Dishes section
// ----------------------------------------------------------

func dishlistpictures(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}
	disheshandler.ListPictures(httpwriter, redisclient, credentials)
}

func dishlist(httpwriter http.ResponseWriter, req *http.Request) {

	if security.ValidateToken(redisclient, req) == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	disheshandler.List(httpwriter, redisclient)
}

func dishadddisplay(httpwriter http.ResponseWriter, req *http.Request) {

	if security.ValidateToken(redisclient, req) == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	disheshandler.LoadDisplayForAdd(httpwriter)
}

func dishadd(httpwriter http.ResponseWriter, httprequest *http.Request) {

	// Retornar credentials e passar para a rotina Add below
	//
	error, _ := security.ValidateTokenV2(redisclient, httprequest)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, httprequest, "/login", 303)
		return
	}

	disheshandler.Add(httpwriter, httprequest, redisclient)
}

func dishupdatedisplay(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, httprequest)

	if error == "NotOkToLogin" {
		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
		return
	}

	disheshandler.LoadDisplayForUpdate(httpresponsewriter, httprequest, redisclient, credentials)
}

func dishupdate(httpwriter http.ResponseWriter, httprequest *http.Request) {

	// Retornar credentials e passar para a rotina Add below
	//
	error, _ := security.ValidateTokenV2(redisclient, httprequest)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, httprequest, "/login", 303)
		return
	}

	disheshandler.Update(httpwriter, httprequest, redisclient)

}

func dishdeletedisplay(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {
	if security.ValidateToken(redisclient, httprequest) == "NotOkToLogin" {
		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
		return
	}
	disheshandler.LoadDisplayForDelete(httpresponsewriter, httprequest, redisclient)

}

func dishdelete(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	// ---------------------------------------------------------------------
	//        Security - Authorisation Check
	// ---------------------------------------------------------------------
	if security.ValidateToken(redisclient, httprequest) == "NotOkToLogin" {
		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
		return
	}
	// ---------------------------------------------------------------------

	dishtoadd := dishes.Dish{}

	dishtoadd.Name = httprequest.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = httprequest.FormValue("dishtype")
	dishtoadd.Price = httprequest.FormValue("dishprice")
	dishtoadd.GlutenFree = httprequest.FormValue("dishglutenfree")
	dishtoadd.DairyFree = httprequest.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = httprequest.FormValue("dishvegetarian")
	dishtoadd.InitialAvailable = httprequest.FormValue("dishinitialavailable")
	dishtoadd.CurrentAvailable = httprequest.FormValue("dishcurrentavailable")

	ret := disheshandler.Dishdelete(mongodbvar, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpresponsewriter, httprequest, "/dishlist", 301)
		return
	}
}

func dishdeletemultiple(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {
	// ---------------------------------------------------------------------
	//        Security - Authorisation Check
	// ---------------------------------------------------------------------
	if security.ValidateToken(redisclient, httprequest) == "NotOkToLogin" {
		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
		return
	}
	// ---------------------------------------------------------------------

	httprequest.ParseForm()

	// Get all selected records
	dishselected := httprequest.Form["dishes"]

	var numrecsel = len(dishselected)

	if numrecsel <= 0 {
		http.Redirect(httpresponsewriter, httprequest, "/dishlist", 301)
		return
	}

	ret := helper.Resultado{}

	ret = disheshandler.DishDeleteMultipleAPI(redisclient, dishselected)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpresponsewriter, httprequest, "/dishlist", 301)
		return
	}

	http.Redirect(httpresponsewriter, httprequest, "/dishlist", 301)
	return

}

func showcache(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {
	// ---------------------------------------------------------------------
	//        Security - Authorisation Check
	// ---------------------------------------------------------------------
	if security.ValidateToken(redisclient, httprequest) == "NotOkToLogin" {
		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
		return
	}
	// ---------------------------------------------------------------------

	// Cache from API
	cachehandler.List(httpresponsewriter, redisclient)

}

func errorpage(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {
	// create new template
	var listtemplate = `
	{{define "listtemplate"}}
	{{ .Info.Name }}
	{{end}}
	`
	t, _ := template.ParseFiles("templates/error.html")
	t, _ = t.Parse(listtemplate)

	t.Execute(httpresponsewriter, listtemplate)
	return
}

// func dishupdateV1(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

// 	if security.ValidateToken(redisclient, httprequest) == "NotOkToLogin" {
// 		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
// 		return
// 	}

// 	dishtoadd := dishes.Dish{}

// 	dishtoadd.Name = httprequest.FormValue("dishname") // This is the key, must be unique
// 	dishtoadd.Type = httprequest.FormValue("dishtype")
// 	dishtoadd.Price = httprequest.FormValue("dishprice")
// 	dishtoadd.GlutenFree = httprequest.FormValue("dishglutenfree")
// 	dishtoadd.DairyFree = httprequest.FormValue("dishdairyfree")
// 	dishtoadd.Vegetarian = httprequest.FormValue("dishvegetarian")
// 	dishtoadd.InitialAvailable = httprequest.FormValue("dishinitialavailable")
// 	dishtoadd.CurrentAvailable = httprequest.FormValue("dishcurrentavailable")
// 	dishtoadd.ImageName = httprequest.FormValue("imagename")
// 	dishtoadd.Description = httprequest.FormValue("dishdescription")
// 	dishtoadd.Descricao = httprequest.FormValue("dishdescricao")

// 	ret := disheshandler.DishupdateAPI(redisclient, dishtoadd)

// 	if ret.IsSuccessful == "Y" {
// 		// http.ServeFile(httpwriter, req, "success.html")
// 		http.Redirect(httpresponsewriter, httprequest, "/dishlist", 301)
// 	} else {
// 		http.Redirect(httpresponsewriter, httprequest, "/errorpage", 301)
// 	}
// 	return
// }
