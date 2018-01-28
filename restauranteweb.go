// Main web application program for restauranteweb
// -----------------------------------------------
// .../src/restauranteweb/restauranteweb.go
// -----------------------------------------------
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"restauranteweb/areas/btcmarketshandler"
	cachehandler "restauranteweb/areas/cachehandler"
	disheshandler "restauranteweb/areas/disheshandler"
	helper "restauranteweb/areas/helper"
	"restauranteweb/areas/ordershandler"
	"restauranteweb/areas/security"

	"strconv"
	"time"

	"github.com/go-redis/redis"
	// _ "github.com/go-sql-driver/mysql"
)

var mongodbvar helper.DatabaseX
var credentials helper.Credentials

var db *sql.DB
var err error
var redisclient *redis.Client

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

	btcmarketshandler.SendEmail(redisclient, "StartingSystemNow"+envirvar.RunningFromServer)

	fmt.Println(">>> Web Server: restauranteweb.exe running.")
	fmt.Println("Loading reference data in cache - Redis")

	mongodbvar.Location = "localhost"
	mongodbvar.Database = "restaurante"

	credentials.UserID = "No User"
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

	err := http.ListenAndServe(":1510", nil) // setting listening port
	// err := http.ListenAndServe(envirvar.WEBServerPort, nil) // setting listening port
	if err != nil {
		//using the mux router
		log.Fatal("ListenAndServe: ", err)
	}

}

func loadreferencedatainredis() {
	// err = client.Set("MongoDB.Location", "{\"MongoDB.Location\":\"192.168.2.180\"}", 0).Err()

	// err = redisclient.Set("Web.MongoDB.Database", "restaurante", 0).Err()
	// err = redisclient.Set("Web.APIServer.Port", ":1520", 0).Err()

	// rodando from raspberry
	// err = redisclient.Set("Web.MongoDB.Location", "localhost", 0).Err()
	//err = redisclient.Set("Web.MongoDB.Location", "192.168.2.180", 0).Err()

	// err = redisclient.Set("Web.APIServer.IPAddress", "http://localhost:1520/", 0).Err()
	// err = redisclient.Set("Web.APIServer.IPAddress", "http://192.168.2.180:1520/", 0).Err()

	// err = redisclient.Set("Web.Debug", "Y", 0).Err()

	variable := helper.Readfileintostruct()
	err = redisclient.Set("Web.MongoDB.Database", variable.APIMongoDBDatabase, 0).Err()
	err = redisclient.Set("Web.APIServer.Port", variable.APIAPIServerPort, 0).Err()
	err = redisclient.Set("WEBServerPort", variable.WEBServerPort, 0).Err()
	err = redisclient.Set("Web.MongoDB.Location", variable.APIMongoDBLocation, 0).Err()
	err = redisclient.Set("Web.APIServer.IPAddress", variable.APIAPIServerIPAddress, 0).Err()
	err = redisclient.Set("Web.Debug", variable.WEBDebug, 0).Err()
	err = redisclient.Set("RecordCurrencyTick", variable.RecordCurrencyTick, 0).Err()
	err = redisclient.Set("RunningFromServer", variable.RunningFromServer, 0).Err()

}

// func root(httpwriter http.ResponseWriter, r *http.Request) {
func root(httpwriter http.ResponseWriter, req *http.Request) {

	if security.ValidateToken(redisclient, req) == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	// create new template
	var listtemplate = `
		{{define "listtemplate"}}
		This is my web site, Daniel - aka D#.
		<p/>
		<p/>
		<picture>
			<img src="images/avatar.png" alt="Avatar" width="400" height="400">
		</picture>
		{{end}}
		`

	t, _ := template.ParseFiles("html/index.html")
	t, _ = t.Parse(listtemplate)

	t.Execute(httpwriter, listtemplate)
	return
}

func root2(httpwriter http.ResponseWriter, r *http.Request) {
	http.ServeFile(httpwriter, r, "index.html")

	return
}

func root3(httpwriter http.ResponseWriter, req *http.Request) {

	if security.ValidateToken(redisclient, req) == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	helper.HomePage(httpwriter, redisclient, credentials)
}

// ----------------------------------------------------------
// Security section
// ----------------------------------------------------------

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "templates/security/signup.html")

		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	passwordvalidate := req.FormValue("passwordvalidate")
	applicationid := req.FormValue("applicationid")

	if password == passwordvalidate {

		// Call API to check if user exists and create
		var resultado = security.SignUp(redisclient, username, password, passwordvalidate, applicationid)
		if resultado.ErrorCode == "200 OK" {

		} else {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		http.Redirect(res, req, "/", 303)
	} else {
		http.Error(res, "Passwords do not match.", 500)
		return
	}

}

func logoutPage(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	cookie, _ := httprequest.Cookie("DanBTCjwt")
	if cookie != nil {
		c := &http.Cookie{
			Name:     "DanBTCjwt",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(httpresponsewriter, c)

	}

	http.Redirect(httpresponsewriter, httprequest, "/", 303)
}

func loginPage(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	if httprequest.Method != "POST" {
		http.ServeFile(httpresponsewriter, httprequest, "templates/security/login.html")
		return
	}

	cookie, _ := httprequest.Cookie("DanBTCjwt")
	if cookie != nil {
		c := &http.Cookie{
			Name:     "DanBTCjwt",
			Value:    "X",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(httpresponsewriter, c)

	}

	// res httpresponsewriter
	// req httprequest

	username := httprequest.FormValue("username")
	password := httprequest.FormValue("password")

	// Check if the user is valid and issue reference token
	//
	var resultado = security.LoginUser(redisclient, username, password)

	if resultado.ErrorCode == "404 Error" {
		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
		return
	}

	// ---------------------------------------------
	//  COOKIE with ERROR !!!!!!!!!!!!!!!!! 28/01/2018 8:45AM
	// -----------------------------------------------
	// -----------------------------------------------
	// -----------------------------------------------
	// -----------------------------------------------
	// -----------------------------------------------
	// -----------------------------------------------

	// Store Token in Cache
	var jwttoken = resultado.ReturnedValue
	year, month, day := time.Now().Date()
	var expiry = strconv.Itoa(int(year)) + strconv.Itoa(int(month)) + strconv.Itoa(int(day))

	// Store token somewhere in desktop
	// Thinking of only using API Key
	// The logon has to exist, but should not call an API instead access the database
	// However we should also allow the API key to access

	// At this point store the token somewhere, cookie or browser storage

	credentials.UserID = username
	credentials.KeyJWT = "DanBTCjwt"
	credentials.JWT = jwttoken
	credentials.Expiry = expiry
	credentials.Roles = []string{"BTC", "RestauranteOwner"}

	jsonval, _ := json.Marshal(credentials)
	jsonstring := string(jsonval)

	// The string had to be escaped because of double quotes.
	// It will have to be unescaped before using
	escapedjson := html.EscapeString(jsonstring)

	// store in redis - server
	_ = redisclient.Set(credentials.KeyJWT, jsonstring, 0).Err()

	// store in cookie
	expiration := time.Now().Add(1 * 2 * time.Hour)

	c := &http.Cookie{
		Name:     credentials.KeyJWT,
		Value:    escapedjson,
		Path:     "/",
		Expires:  expiration,
		MaxAge:   0,
		HttpOnly: true,
	}

	http.SetCookie(httpresponsewriter, c)

	http.Redirect(httpresponsewriter, httprequest, "/", 303)
	return
}

func loginPageV2(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	if httprequest.Method != "POST" {
		http.ServeFile(httpresponsewriter, httprequest, "templates/security/login.html")
		return
	}

	cookie, _ := httprequest.Cookie("DanBTCjwt")
	if cookie != nil {
		c := &http.Cookie{
			Name:     "DanBTCjwt",
			Value:    "X",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(httpresponsewriter, c)

	}

	// res httpresponsewriter
	// req httprequest

	username := httprequest.FormValue("username")
	password := httprequest.FormValue("password")

	// Check if the user is valid and issue reference token
	//
	var resultado = security.LoginUserV2(redisclient, username, password)

	if resultado.JWT == "Error" {

		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
		return
	}

	// Store Token in Cache
	var jwttoken = resultado.JWT
	year, month, day := time.Now().Date()
	var expiry = strconv.Itoa(int(year)) + strconv.Itoa(int(month)) + strconv.Itoa(int(day))

	// Store token somewhere in desktop
	// Thinking of only using API Key
	// The logon has to exist, but should not call an API instead access the database
	// However we should also allow the API key to access

	// At this point store the token somewhere, cookie or browser storage

	credentials.UserID = username
	credentials.KeyJWT = "DanBTCjwt"
	credentials.JWT = jwttoken
	credentials.Expiry = expiry
	credentials.Roles = []string{"BTC", "RestauranteOwner"}
	credentials.ClaimSet = resultado.ClaimSet
	credentials.ApplicationID = resultado.ApplicationID

	jsonval, _ := json.Marshal(credentials)
	jsonstring := string(jsonval)

	// The string had to be escaped because of double quotes.
	// It will have to be unescaped before using
	escapedjson := html.EscapeString(jsonstring)

	// store in redis - server
	_ = redisclient.Set(credentials.KeyJWT, jsonstring, 0).Err()
	_ = redisclient.Set("ApplicationID", credentials.ApplicationID, 0).Err()
	_ = redisclient.Set("UserID", username, 0).Err()

	// store in cookie
	expiration := time.Now().Add(1 * 2 * time.Hour)

	c := &http.Cookie{
		Name:     credentials.KeyJWT,
		Value:    escapedjson,
		Path:     "/",
		Expires:  expiration,
		MaxAge:   0,
		HttpOnly: true,
	}

	http.SetCookie(httpresponsewriter, c)

	http.Redirect(httpresponsewriter, httprequest, "/", 303)
	return
}

// ----------------------------------------------------------
// Orders section
// ----------------------------------------------------------

func orderlist(httpwriter http.ResponseWriter, req *http.Request) {

	if security.ValidateToken(redisclient, req) == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	ordershandler.List(httpwriter, redisclient)
}

func orderadddisplay(httpwriter http.ResponseWriter, req *http.Request) {
	ordershandler.LoadDisplayForAdd(httpwriter, redisclient)
}

func orderadd(httpwriter http.ResponseWriter, req *http.Request) {
	ordershandler.Add(httpwriter, req, redisclient)
}

func orderviewdisplay(httpwriter http.ResponseWriter, req *http.Request) {
	ordershandler.LoadDisplayForView(httpwriter, req, redisclient)
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

	if security.ValidateToken(redisclient, req) == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	btcmarketshandler.PreOrderList(httpwriter, redisclient)
}

func btcmarketslistV3(httpwriter http.ResponseWriter, req *http.Request) {

	if security.ValidateToken(redisclient, req) == "NotOkToLogin" {
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

	if security.ValidateToken(redisclient, req) == "NotOkToLogin" {
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

	if security.ValidateToken(redisclient, req) == "NotOkToLogin" {
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
// Dishes section
// ----------------------------------------------------------

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

func dishadd(httpwriter http.ResponseWriter, req *http.Request) {
	if security.ValidateToken(redisclient, req) == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	disheshandler.Add(httpwriter, req, redisclient)
}

func dishupdatedisplay(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	if security.ValidateToken(redisclient, httprequest) == "NotOkToLogin" {
		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
		return
	}

	disheshandler.LoadDisplayForUpdate(httpresponsewriter, httprequest, redisclient)
}

func dishupdate(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	if security.ValidateToken(redisclient, httprequest) == "NotOkToLogin" {
		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
		return
	}

	dishtoadd := disheshandler.Dish{}

	dishtoadd.Name = httprequest.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = httprequest.FormValue("dishtype")
	dishtoadd.Price = httprequest.FormValue("dishprice")
	dishtoadd.GlutenFree = httprequest.FormValue("dishglutenfree")
	dishtoadd.DairyFree = httprequest.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = httprequest.FormValue("dishvegetarian")

	ret := disheshandler.DishupdateAPI(redisclient, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpresponsewriter, httprequest, "/dishlist", 301)
	} else {
		http.Redirect(httpresponsewriter, httprequest, "/errorpage", 301)
	}
	return
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

	dishtoadd := disheshandler.Dish{}

	dishtoadd.Name = httprequest.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = httprequest.FormValue("dishtype")
	dishtoadd.Price = httprequest.FormValue("dishprice")
	dishtoadd.GlutenFree = httprequest.FormValue("dishglutenfree")
	dishtoadd.DairyFree = httprequest.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = httprequest.FormValue("dishvegetarian")

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
