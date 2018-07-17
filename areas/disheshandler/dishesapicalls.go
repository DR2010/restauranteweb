// Package disheshandler API calls for dishes web
// --------------------------------------------------------------
// .../src/restauranteweb/areas/disherhandler/dishesapicalls.go
// --------------------------------------------------------------
package disheshandler

import (
	"encoding/json"
	helper "festajuninaweb/areas/helper"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-redis/redis"

	dishes "restauranteapi/models"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// // Dish is to be exported
// type Dish struct {
// 	SystemID         bson.ObjectId `json:"id"        bson:"_id,omitempty"`
// 	Name             string        // name of the dish - this is the KEY, must be unique
// 	Type             string        // type of dish, includes drinks and deserts
// 	Price            string        // preco do prato multiplicar por 100 e nao ter digits
// 	GlutenFree       string        // Gluten free dishes
// 	DairyFree        string        // Dairy Free dishes
// 	Vegetarian       string        // Vegeterian dishes
// 	InitialAvailable string        // Number of items initially available
// 	CurrentAvailable string        // Currently available
// 	ImageName        string        // Image Name
// }

// ListDishes works
func listdishes(redisclient *redis.Client) []dishes.Dish {

	var apiserver string
	apiserver, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	urlrequest := apiserver + "/dishlist"

	// urlrequest = "http://localhost:1520/dishlist"

	url := fmt.Sprintf(urlrequest)

	var emptydisplay []dishes.Dish

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return emptydisplay
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	var dishlist []dishes.Dish

	if err := json.NewDecoder(resp.Body).Decode(&dishlist); err != nil {
		log.Println(err)
	}

	return dishlist
}

// APIcallAdd is
func APIcallAdd(redisclient *redis.Client, dishInsert dishes.Dish) helper.Resultado {

	mongodbvar := new(helper.DatabaseX)

	mongodbvar.APIServer, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	// mongodbvar.APIServer = "http://localhost:1520/"

	apiURL := mongodbvar.APIServer
	resource := "/dishadd"

	data := url.Values{}
	data.Add("dishname", dishInsert.Name)
	data.Add("dishtype", dishInsert.Type)
	data.Add("dishprice", dishInsert.Price)
	data.Add("dishglutenfree", dishInsert.GlutenFree)
	data.Add("dishdairyfree", dishInsert.DairyFree)
	data.Add("dishvegetarian", dishInsert.Vegetarian)
	data.Add("dishinitialavailable", dishInsert.InitialAvailable)
	data.Add("dishcurrentavailable", dishInsert.CurrentAvailable)
	data.Add("dishimagename", dishInsert.ImageName)
	data.Add("dishdescription", dishInsert.Description)
	data.Add("dishdescricao", dishInsert.Descricao)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(data.Encode())
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	fmt.Println("resp2.Status:" + resp2.Status)
	fmt.Println("body:" + data.Encode())

	var emptydisplay helper.Resultado
	emptydisplay.ErrorCode = resp2.Status

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
	} else {
		emptydisplay.IsSuccessful = "N"
	}

	return emptydisplay

}

// FindAPI is to find stuff
func FindAPI(redisclient *redis.Client, dishFind string) dishes.Dish {

	var apiserver string
	apiserver, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	// This is essential! Because if the string has spaces it doesn't work without the escape
	// Bolo de Cenoura = Bolo+de+Cenoura   >>> Works as a dream!
	dishfindescaped := url.QueryEscape(dishFind)
	urlrequest := apiserver + "/dishfind?dishname=" + dishfindescaped

	urlrequestencoded, _ := url.ParseRequestURI(urlrequest)
	// url := fmt.Sprintf(urlrequest)
	url := urlrequestencoded.String()
	// tw.Text = strings.Replace(tw.Text, " ", "+", -1)
	// urlx := url.QueryEscape(urlrequest)

	var emptydisplay dishes.Dish

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return emptydisplay
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	var dishback dishes.Dish

	if err := json.NewDecoder(resp.Body).Decode(&dishback); err != nil {
		log.Println(err)
	}

	return dishback

}

// DishupdateAPI is
func DishupdateAPI(redisclient *redis.Client, dishUpdate dishes.Dish) helper.Resultado {

	mongodbvar := new(helper.DatabaseX)

	mongodbvar.APIServer, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	apiURL := mongodbvar.APIServer
	resource := "/dishupdate"

	data := url.Values{}
	data.Add("dishname", dishUpdate.Name)
	data.Add("dishtype", dishUpdate.Type)
	data.Add("dishprice", dishUpdate.Price)
	data.Add("dishglutenfree", dishUpdate.GlutenFree)
	data.Add("dishdairyfree", dishUpdate.DairyFree)
	data.Add("dishvegetarian", dishUpdate.Vegetarian)
	data.Add("dishinitialavailable", dishUpdate.InitialAvailable)
	data.Add("dishcurrentavailable", dishUpdate.CurrentAvailable)
	data.Add("dishimagename", dishUpdate.ImageName)
	data.Add("dishdescription", dishUpdate.Description)
	data.Add("dishdescricao", dishUpdate.Descricao)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(data.Encode())
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	fmt.Println("resp2.Status:" + resp2.Status)

	var emptydisplay helper.Resultado
	emptydisplay.ErrorCode = resp2.Status

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
	}

	return emptydisplay
}

// DishdeleteAPI is
func DishdeleteAPI(redisclient *redis.Client, dishUpdate dishes.Dish) helper.Resultado {

	mongodbvar := new(helper.DatabaseX)

	mongodbvar.APIServer, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	apiURL := mongodbvar.APIServer
	resource := "/dishdelete"

	data := url.Values{}
	data.Add("dishname", dishUpdate.Name)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(data.Encode())
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	fmt.Println("resp2.Status:" + resp2.Status)

	var emptydisplay helper.Resultado
	emptydisplay.ErrorCode = resp2.Status

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
	}

	return emptydisplay
}

// Dishdelete is
func Dishdelete(database helper.DatabaseX, objectDelete dishes.Dish) helper.Resultado {

	database.Collection = "dishes"

	session, err := mgo.Dial(database.Location)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(database.Database).C(database.Collection)

	err = collection.Remove(bson.M{"name": objectDelete.Name})

	if err != nil {
		log.Fatal(err)
	}

	var res helper.Resultado
	res.ErrorCode = "0001"
	res.ErrorDescription = "Something Happened"
	res.IsSuccessful = "Y"

	return res
}

// DishDeleteMultipleAPI is
func DishDeleteMultipleAPI(redisclient *redis.Client, dishestodelete []string) helper.Resultado {

	mongodbvar := new(helper.DatabaseX)

	mongodbvar.APIServer, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	apiURL := mongodbvar.APIServer
	resource := "/dishdelete"

	data := url.Values{}
	data.Add("dishname", dishestodelete[0])

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(data.Encode())
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	fmt.Println("resp2.Status:" + resp2.Status)

	var emptydisplay helper.Resultado
	emptydisplay.ErrorCode = resp2.Status

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
	}

	return emptydisplay
}
