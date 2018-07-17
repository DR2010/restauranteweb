// main web application program for festajuninaweb
// -----------------------------------------------
// .../src/festajuninaweb/festajuninaweb.go
// -----------------------------------------------
package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	cachehandler "festajuninaweb/areas/cachehandler"
	"festajuninaweb/areas/ordershandler"
	"festajuninaweb/areas/security"
	"fmt"
	"html/template"
	"log"
	"net/http"
	// The Models are shared by WEB and API
	disheshandler "festajuninaweb/areas/disheshandler"
	helper "festajuninaweb/areas/helper"
	dishes "restauranteapi/models"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

// Message our message object
type Message struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

var mongodbvar helper.DatabaseX

var clients []Client

// var credentials helper.Credentials

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
	envirvar.AppBelnorthEnabled, _ = redisclient.Get("AppBelnorthEnabled").Result()
	envirvar.AppBitcoinEnabled, _ = redisclient.Get("AppBitcoinEnabled").Result()
	envirvar.AppFestaJuninaEnabled, _ = redisclient.Get("AppFestaJuninaEnabled").Result()

	fmt.Println(">>> Web Server: restauranteweb.exe running.")
	fmt.Println("Loading reference data in cache - Redis")

	mongodbvar.Location = "localhost"
	mongodbvar.Database = "restaurante"

	// mongodbvar.APIServer = "http://192.168.2.180:1520/"
	// mongodbvar.APIServer = "http://localhost:1520/"

	fmt.Println("FESTA JUNINA web server")
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

	err := http.ListenAndServe(":1710", nil) // setting listening port
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

	_, credentials := security.ValidateTokenV2(redisclient, req)

	helper.HomePage(httpwriter, redisclient, credentials)

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

func instructions(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	security.Instructions(httpresponsewriter, httprequest, redisclient)

}

// ----------------------------------------------------------
// Orders section
// ----------------------------------------------------------

// ----------------
// Anonymous
// ----------------
func orderlist(httpwriter http.ResponseWriter, req *http.Request) {
	_, credentials := security.ValidateTokenV2(redisclient, req)

	// If user is not ADMIN, show only users order

	ordershandler.ListV2(httpwriter, redisclient, credentials)
	// ordershandler.ListV3OnlyPlaced(httpwriter, redisclient, credentials)
}

func orderadddisplay(httpwriter http.ResponseWriter, req *http.Request) {
	_, credentials := security.ValidateTokenV2(redisclient, req)

	ordershandler.LoadDisplayForAdd(httpwriter, redisclient, credentials)
}

// orderadd
// is designed to place an order for a client that is logged on
//
func orderadd(httpwriter http.ResponseWriter, req *http.Request) {

	_, _ = security.ValidateTokenV2(redisclient, req)

	ordershandler.Add(httpwriter, req, redisclient)
}

// orderclientadd
// is designed to place an order for an anonymous client
// it creates a dummy client
func orderclientadd(httpwriter http.ResponseWriter, req *http.Request) {

	// Find token
	// Get user ID
	_, credentials := security.ValidateTokenV2(redisclient, req)

	ordershandler.AddOrderClient(httpwriter, req, redisclient, credentials)
}

func orderviewdisplay(httpwriter http.ResponseWriter, req *http.Request) {

	_, credentials := security.ValidateTokenV2(redisclient, req)

	ordershandler.LoadDisplayForView(httpwriter, req, redisclient, credentials)
}
func ordercancel(httpwriter http.ResponseWriter, req *http.Request) {
	_, _ = security.ValidateTokenV2(redisclient, req)

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

// ----------------

func orderlistcompleted(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	// Only Admin
	//
	if credentials.IsAdmin == "Yes" {
		ordershandler.ListCompleted(httpwriter, redisclient, credentials)
	}

}

func orderliststatus(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	// Only Admin
	//
	if credentials.IsAdmin == "Yes" {
		ordershandler.ListStatus(req, httpwriter, redisclient, credentials)
	}

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

func orderStartServing(httpwriter http.ResponseWriter, req *http.Request) {
	error, _ := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}
	ordershandler.StartServing(httpwriter, req, redisclient)
}

// ----------------------------------------------------------
// Dishes section
// ----------------------------------------------------------

func dishlistpictures(httpwriter http.ResponseWriter, req *http.Request) {

	_, credentials := security.ValidateTokenV2(redisclient, req)

	disheshandler.ListPictures(httpwriter, redisclient, credentials)
}

func dishlist(httpwriter http.ResponseWriter, req *http.Request) {

	error, credentials := security.ValidateTokenV2(redisclient, req)

	if error == "NotOkToLogin" {
		http.Redirect(httpwriter, req, "/login", 303)
		return
	}

	disheshandler.List(httpwriter, redisclient, credentials)
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

// ---------------------------------------------------
// Websockets
// ---------------------------------------------------
// ---------------------------------------------------
// func rootws(httpwriter http.ResponseWriter, httprequest *http.Request) {
// 	flag.Parse()
// 	hub := newHub()
// 	go hub.run()
// 	serveWs(hub, httpwriter, httprequest)

// 	// Upgrade initial GET request to a websocket

// 	wsconn, err := upgrader.Upgrade(httpwriter, httprequest, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Register our new client
// 	clients[wsconn] = true

// 	fmt.Println("Client subscribed")

// }

func rootws(httpwriter http.ResponseWriter, httprequest *http.Request) {
	conn, err := websocket.Upgrade(httpwriter, httprequest, nil, 1024, 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(msg))
	}

}

func broadcastmsg(httpwriter http.ResponseWriter, httprequest *http.Request) {
	fmt.Println("broadcast message")
}

func broadcastHandler(w http.ResponseWriter, r *http.Request) {
	msg := []byte("Message")

	for _, c := range clients {
		broadcastsingle(msg, c)
	}

	fmt.Fprintf(w, "Broadcasting %v", msg)
}

func addClientAndGreet(list []Client, client Client) []Client {
	clients = append(list, client)
	websocket.WriteJSON(client.conn, Message{"Server", "Welcome!"})
	return clients
}

func broadcastsingle(msg []byte, client Client) {
	fmt.Printf("Broadcasting %+v\n", msg)
	client.send <- msg
}
