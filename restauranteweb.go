package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	cachehandler "restauranteweb/areas/cachehandler"
	disheshandler "restauranteweb/areas/disheshandler"
	helper "restauranteweb/areas/helper"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"

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

	fmt.Println(">>> Web Server: restauranteweb.exe running.")
	fmt.Println("Loading reference data in cache - Redis")

	mongodbvar.Location = "localhost"
	mongodbvar.Database = "restaurante"
	mongodbvar.APIServer = "http://localhost:1520/"

	fmt.Println("Running... Listening to :1515 - print")
	fmt.Println("MongoDB location: " + mongodbvar.Location)
	fmt.Println("MongoDB database: " + mongodbvar.Database)
	fmt.Println("API Server: " + mongodbvar.APIServer)

	router := XNewRouter()

	// handle using the router mux
	//
	http.Handle("/", router) // setting router rule

	http.Handle("/html/", http.StripPrefix("/html", http.FileServer(http.Dir("./"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./js"))))
	http.Handle("/ts/", http.StripPrefix("/ts", http.FileServer(http.Dir("./ts"))))
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts", http.FileServer(http.Dir("./fonts"))))

	err := http.ListenAndServe(":1515", nil) // setting listening port
	if err != nil {
		//using the mux router
		log.Fatal("ListenAndServe: ", err)
	}
}

func loadreferencedatainredis() {
	// err = client.Set("MongoDB.Location", "{\"MongoDB.Location\":\"192.168.2.180\"}", 0).Err()
	// err = redisclient.Set("Web.MongoDB.Location", "192.168.2.180", 0).Err()
	err = redisclient.Set("Web.MongoDB.Location", "localhost", 0).Err()
	err = redisclient.Set("Web.MongoDB.Database", "restaurante", 0).Err()
	err = redisclient.Set("Web.APIServer.IPAddress", "http://localhost:1520/", 0).Err()
	err = redisclient.Set("Web.APIServer.Port", ":1520", 0).Err()
}

func root(httpwriter http.ResponseWriter, r *http.Request) {

	// create new template
	var listtemplate = `
		{{define "listtemplate"}}
	
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

func loginPage(res http.ResponseWriter, req *http.Request) {
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

func dishdeletedisplay(httpwriter http.ResponseWriter, req *http.Request) {

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

	dishtodelete := disheshandler.Dish{}

	ret := helper.Resultado{}

	for x := 0; x < len(dishselected); x++ {

		dishtodelete.Name = dishselected[x]

		ret = disheshandler.Dishdelete(mongodbvar, dishtodelete)
	}

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

// ----------------------------------------------------------------------
// ----------------------------------------------------------------------
// ----------------------------------------------------------------------
// ----------------------------------------------------------------------
// This is the section of methods to be deleted when it is all working
// ----------------------------------------------------------------------
// ----------------------------------------------------------------------
// ----------------------------------------------------------------------
// ----------------------------------------------------------------------

func dishupdatedisplayTBD(httpwriter http.ResponseWriter, req *http.Request) {

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
	t, _ := template.ParseFiles("templates/indextemplate.html", "templates/dishupdate.html")

	items := DisplayTemplate{}
	items.Info.Name = "Dish Add"

	items.DishItem = disheshandler.Dish{}
	items.DishItem.Name = dishselected[0]

	var dishfind = disheshandler.Dish{}
	var dishname = items.DishItem.Name

	dishfind = disheshandler.FindAPI(redisclient, dishname)
	items.DishItem = dishfind

	t.Execute(httpwriter, items)

	return

}
