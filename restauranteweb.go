// Main web application program for restauranteweb
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
	"restauranteweb/areas/belnorthhandler"
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

// var credentials helper.Credentials

var db *sql.DB
var err error
var redisclient *redis.Client
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

	btcmarketshandler.SendEmail(redisclient, "StartingSystemNow"+envirvar.RunningFromServer)

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
	PaymentStoreMemory(redisclient)

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

func signupPage(httpresponsewriter http.ResponseWriter, req *http.Request) {

	type ControllerInfo struct {
		Name    string
		Message string
	}
	type DisplayTemplate struct {
		Info ControllerInfo
	}

	items := DisplayTemplate{}
	items.Info.Name = "Login Page"

	if req.Method != "POST" {

		t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = ""
		t.Execute(httpresponsewriter, items)

		// http.ServeFile(res, req, "templates/security/signup.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	passwordvalidate := req.FormValue("passwordvalidate")
	applicationid := req.FormValue("applicationid")

	if username == "" {
		t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "Please enter details."
		t.Execute(httpresponsewriter, items)
		return
	}

	if password == "" {
		t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "Please enter details."
		t.Execute(httpresponsewriter, items)
		return
	}

	if password == passwordvalidate {

		// Call API to check if user exists and create
		var resultado = security.SignUp(redisclient, username, password, passwordvalidate, applicationid)
		if resultado.ErrorCode == "200 OK" {

		} else {
			t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
			items.Info.Message = "Passwords mismatch."
			t.Execute(httpresponsewriter, items)
			return
		}

		http.Redirect(httpresponsewriter, req, "/", 303)
	} else {
		t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "Passwords do not match."
		t.Execute(httpresponsewriter, items)
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

func loginPageV4(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	type ControllerInfo struct {
		Name    string
		Message string
	}
	type DisplayTemplate struct {
		Info ControllerInfo
	}

	items := DisplayTemplate{}
	items.Info.Name = "Login Page"

	if httprequest.Method != "POST" {

		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = ""
		t.Execute(httpresponsewriter, items)

		// http.ServeFile(httpresponsewriter, httprequest, "templates/security/login.html")
		return
	}

	username := httprequest.FormValue("username")
	password := httprequest.FormValue("password")

	if username == "" {
		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "Please enter details."
		t.Execute(httpresponsewriter, items)
		return
	}

	if password == "" {
		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "Please enter details."
		t.Execute(httpresponsewriter, items)
		return
	}

	cookiekeyJWT := "DanBTCjwt"
	cookiekeyUSERID := "DanBTCuserid"

	cookieJWT, _ := httprequest.Cookie(cookiekeyJWT)
	cookieUSERID, _ := httprequest.Cookie(cookiekeyUSERID)

	if cookieJWT != nil {
		cokJWT := &http.Cookie{
			Name:     cookiekeyJWT,
			Value:    "X",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(httpresponsewriter, cokJWT)
	}

	if cookieUSERID != nil {
		cokUSERID := &http.Cookie{
			Name:     cookiekeyUSERID,
			Value:    "X",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(httpresponsewriter, cokUSERID)
	}

	// Check if the user is valid and issue reference token
	//
	var resultado = security.LoginUserV2(redisclient, username, password)

	if resultado.JWT == "Error" {
		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "Login error. Try again."
		t.Execute(httpresponsewriter, items)
		return
	}

	// Store Token in Cache
	var jwttoken = resultado.JWT
	year, month, day := time.Now().Date()
	var expiry = strconv.Itoa(int(year)) + strconv.Itoa(int(month)) + strconv.Itoa(int(day))

	rediskey := "DanBTCjwt" + username

	var credentials helper.Credentials
	credentials.UserID = username
	credentials.KeyJWT = rediskey
	credentials.JWT = jwttoken
	credentials.Expiry = expiry
	credentials.ClaimSet = resultado.ClaimSet
	credentials.ApplicationID = resultado.ApplicationID

	jsonval, _ := json.Marshal(credentials)
	jsonstring := string(jsonval)

	_ = redisclient.Set(rediskey, jsonstring, 0).Err()

	// store in cookie
	expiration := time.Now().Add(1 * 2 * time.Hour)

	cokJWT := &http.Cookie{
		Name:     cookiekeyJWT,
		Value:    jwttoken,
		Path:     "/",
		Expires:  expiration,
		MaxAge:   0,
		HttpOnly: true,
	}

	http.SetCookie(httpresponsewriter, cokJWT)

	cokUSERID := &http.Cookie{
		Name:     cookiekeyUSERID,
		Value:    username,
		Path:     "/",
		Expires:  expiration,
		MaxAge:   0,
		HttpOnly: true,
	}

	http.SetCookie(httpresponsewriter, cokUSERID)

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

	http.Redirect(httpwriter, req, "/orderlist", 303)
}

func ordersettoready(httpwriter http.ResponseWriter, req *http.Request) {
	error, _ := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}
	ordershandler.OrderisReady(httpwriter, req, redisclient)

	http.Redirect(httpwriter, req, "/orderlist", 303)
}

func orderviewdisplay(httpwriter http.ResponseWriter, req *http.Request) {
	error, _ := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}
	ordershandler.LoadDisplayForView(httpwriter, req, redisclient)
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

	error, credentials := security.ValidateTokenV2(redisclient, httprequest)

	if error == "NotOkToLogin" {
		http.Redirect(httpresponsewriter, httprequest, "/login", 303)
		return
	}

	disheshandler.LoadDisplayForUpdate(httpresponsewriter, httprequest, redisclient, credentials)
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

// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------
// -----------------  DELETE ----------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------

func loginPageV3(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	if httprequest.Method != "POST" {
		http.ServeFile(httpresponsewriter, httprequest, "templates/security/login.html")
		return
	}

	username := httprequest.FormValue("username")
	password := httprequest.FormValue("password")

	cookiekeyJWT := "DanBTCjwt"
	cookiekeyUSERID := "DanBTCuserid"

	cookieJWT, _ := httprequest.Cookie(cookiekeyJWT)
	cookieUSERID, _ := httprequest.Cookie(cookiekeyUSERID)

	if cookieJWT != nil {
		cokJWT := &http.Cookie{
			Name:     cookiekeyJWT,
			Value:    "X",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(httpresponsewriter, cokJWT)
	}

	if cookieUSERID != nil {
		cokUSERID := &http.Cookie{
			Name:     cookiekeyUSERID,
			Value:    "X",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(httpresponsewriter, cokUSERID)
	}

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

	rediskey := "DanBTCjwt" + username

	var credentials helper.Credentials
	credentials.UserID = username
	credentials.KeyJWT = rediskey
	credentials.JWT = jwttoken
	credentials.Expiry = expiry
	credentials.ClaimSet = resultado.ClaimSet
	credentials.ApplicationID = resultado.ApplicationID

	jsonval, _ := json.Marshal(credentials)
	jsonstring := string(jsonval)

	_ = redisclient.Set(rediskey, jsonstring, 0).Err()

	// store in cookie
	expiration := time.Now().Add(1 * 2 * time.Hour)

	cokJWT := &http.Cookie{
		Name:     cookiekeyJWT,
		Value:    jwttoken,
		Path:     "/",
		Expires:  expiration,
		MaxAge:   0,
		HttpOnly: true,
	}

	http.SetCookie(httpresponsewriter, cokJWT)

	cokUSERID := &http.Cookie{
		Name:     cookiekeyUSERID,
		Value:    username,
		Path:     "/",
		Expires:  expiration,
		MaxAge:   0,
		HttpOnly: true,
	}

	http.SetCookie(httpresponsewriter, cokUSERID)

	http.Redirect(httpresponsewriter, httprequest, "/", 303)
	return
}
