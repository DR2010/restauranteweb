// Main web application program for restauranteweb
// -----------------------------------------------
// .../src/restauranteweb/restauranteweb.go
// -----------------------------------------------
package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"restauranteweb/areas/btcmarketshandler"
	cachehandler "restauranteweb/areas/cachehandler"
	disheshandler "restauranteweb/areas/disheshandler"
	helper "restauranteweb/areas/helper"
	"restauranteweb/areas/ordershandler"
	security "restauranteweb/areas/security"

	"github.com/go-redis/redis"
	// _ "github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"
)

var mongodbvar helper.DatabaseX

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

func root(httpwriter http.ResponseWriter, r *http.Request) {

	// create new template
	var listtemplate = `
		{{define "listtemplate"}}
	    This is our restaurant!
		{{end}}
		`

	t, _ := template.ParseFiles("templates/indextemplate.html")
	t, _ = t.Parse(listtemplate)

	t.Execute(httpwriter, listtemplate)
	return
}

func root2(httpwriter http.ResponseWriter, r *http.Request) {
	http.ServeFile(httpwriter, r, "index.html")

	return
}

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "templates/signup.html")

		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	passwordvalidate := req.FormValue("passwordvalidate")

	if password == passwordvalidate {

		// Call API to check if user exists and create
		var resultado = security.SignUp(redisclient, username, password, passwordvalidate)
		if resultado.ErrorCode == "200 OK" {

		} else {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		http.Redirect(res, req, "/", 301)
	} else {
		http.Error(res, "Passwords do not match.", 500)
		return
	}

}

func signupPageOLD(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "templates/signup.html")

		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var user string

	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
	}
}

func loginPage(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	if httprequest.Method != "POST" {
		http.ServeFile(httpresponsewriter, httprequest, "templates/login.html")
		return
	}

	// res httpresponsewriter
	// req httprequest

	username := httprequest.FormValue("username")
	password := httprequest.FormValue("password")

	// Check if the user is valid and issue reference token
	//
	var resultado = security.LoginUser(redisclient, username, password)

	if resultado.ErrorCode == "404 Error" {
		http.Redirect(httpresponsewriter, httprequest, "/loginPage", 301)
		return
	}

	// Store token somewhere in desktop
	// Thinking of only using API Key
	// The logon has to exist, but should not call an API instead access the database
	// However we should also allow the API key to access

	// At this point store the token somewhere, cookie or browser storage

	c := http.Cookie{
		Name:  "DanBTCjwt",
		Value: resultado.ReturnedValue,
	}
	http.SetCookie(httpresponsewriter, &c)

	http.Redirect(httpresponsewriter, httprequest, "/", 301)
	return

}

func loginPageOLD(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "templates/login.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(res, req, "/loginPage", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/loginPage", 301)
		return
	}

	res.Write([]byte("Hello" + databaseUsername))

}

// ----------------------------------------------------------
// Orders section
// ----------------------------------------------------------

func orderlist(httpwriter http.ResponseWriter, req *http.Request) {

	cookie, _ := req.Cookie("DanBTCjwt")

	if security.ValidateToken(redisclient, "Daniel", cookie.Value) == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 301)
	} else {
		ordershandler.List(httpwriter, redisclient)
	}
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

func btcmarketslistV1(httpwriter http.ResponseWriter, req *http.Request) {
	// var listofbit = btcmarketshandler.List(httpwriter, redisclient)
	btcmarketshandler.List(httpwriter, redisclient)

	// btcmarketshandler.Add(listofbit, redisclient)
}

// This is V2 - activated on 11:54AM 01-Jan-2018
// This version logs the rates every minute
func btcmarketslistV2(httpwriter http.ResponseWriter, req *http.Request) {
	var listofbit = btcmarketshandler.ListV2(httpwriter, redisclient)

	btcmarketshandler.Add(listofbit, redisclient)
}

func btcmarketslistV3(httpwriter http.ResponseWriter, req *http.Request) {

	envirvar := new(helper.RestEnvVariables)
	envirvar.RecordCurrencyTick, _ = redisclient.Get("RecordCurrencyTick").Result()

	var listofbit = btcmarketshandler.ListV2(httpwriter, redisclient)

	if envirvar.RecordCurrencyTick == "Y" {
		btcmarketshandler.Add(listofbit, redisclient)
	}
}

func btclistcoinshistory(httpwriter http.ResponseWriter, req *http.Request) {

	params := req.URL.Query()
	var currency = params.Get("currency")
	if currency == "" {
		currency = "ALL"
	}

	var rows = params.Get("rows")
	if rows == "" {
		rows = "100"
	}

	btcmarketshandler.HListHistory(httpwriter, redisclient, currency, rows)

}

// ----------------------------------------------------------
// Dishes section
// ----------------------------------------------------------

func dishlist(httpwriter http.ResponseWriter, req *http.Request) {
	disheshandler.List(httpwriter, redisclient)
}

func dishadddisplay(httpwriter http.ResponseWriter, req *http.Request) {
	disheshandler.LoadDisplayForAdd(httpwriter)
}

func dishadd(httpwriter http.ResponseWriter, req *http.Request) {
	disheshandler.Add(httpwriter, req, redisclient)
}

func dishupdatedisplay(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {
	disheshandler.LoadDisplayForUpdate(httpresponsewriter, httprequest, redisclient)
}

func dishupdate(httpwriter http.ResponseWriter, req *http.Request) {

	dishtoadd := disheshandler.Dish{}

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")

	ret := disheshandler.DishupdateAPI(redisclient, dishtoadd)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/dishlist", 301)
	} else {
		http.Redirect(httpwriter, req, "/errorpage", 301)
	}
	return
}

func dishdeletedisplay(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	disheshandler.LoadDisplayForDelete(httpresponsewriter, httprequest, redisclient)

}
func dishdeletedisplayTBD(httpwriter http.ResponseWriter, req *http.Request) {

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
		DishItem   disheshandler.Dish
	}

	// create new template
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/dishdelete.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Delete"

	items.DishItem = disheshandler.Dish{}
	items.DishItem.Name = dishselected[0]

	var dishfind = disheshandler.Dish{}
	var dishname = items.DishItem.Name

	dishfind = disheshandler.FindAPI(redisclient, dishname)
	items.DishItem = dishfind

	t.Execute(httpwriter, items)

	return

}

func dishdelete(httpwriter http.ResponseWriter, req *http.Request) {

	dishtoadd := disheshandler.Dish{}

	dishtoadd.Name = req.FormValue("dishname") // This is the key, must be unique
	dishtoadd.Type = req.FormValue("dishtype")
	dishtoadd.Price = req.FormValue("dishprice")
	dishtoadd.GlutenFree = req.FormValue("dishglutenfree")
	dishtoadd.DairyFree = req.FormValue("dishdairyfree")
	dishtoadd.Vegetarian = req.FormValue("dishvegetarian")

	ret := disheshandler.Dishdelete(mongodbvar, dishtoadd)

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

	ret := helper.Resultado{}

	ret = disheshandler.DishDeleteMultipleAPI(redisclient, dishselected)

	if ret.IsSuccessful == "Y" {
		// http.ServeFile(httpwriter, req, "success.html")
		http.Redirect(httpwriter, req, "/dishlist", 301)
		return
	}

	http.Redirect(httpwriter, req, "/dishlist", 301)
	return

}

func showcache(httpwriter http.ResponseWriter, req *http.Request) {
	cachehandler.List(httpwriter, redisclient)
}

func errorpage(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {
	// create new template
	var listtemplate = `
	{{define "listtemplate"}}

	{{end}}
	`
	t, _ := template.ParseFiles("templates/error.html")
	t, _ = t.Parse(listtemplate)

	t.Execute(httpresponsewriter, listtemplate)
	return
}
