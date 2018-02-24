// Package ordershandler API calls for dishes web
// --------------------------------------------------------------
// .../src/restauranteweb/areas/disherhandler/orderapicalls.go
// --------------------------------------------------------------
package ordershandler

import (
	"encoding/json"
	"fmt"
	"log"
	"mongodb/dishes"
	"net/http"
	"net/url"
	helper "restauranteweb/areas/helper"
	"strings"

	order "restauranteapi/models"

	"github.com/go-redis/redis"

	"gopkg.in/mgo.v2/bson"
)

// Dish is to be exported
type Dish struct {
	SystemID   bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	Name       string        // name of the dish - this is the KEY, must be unique
	Type       string        // type of dish, includes drinks and deserts
	Price      string        // preco do prato multiplicar por 100 e nao ter digits
	GlutenFree string        // Gluten free dishes
	DairyFree  string        // Dairy Free dishes
	Vegetarian string        // Vegeterian dishes
}

// SearchCriteria is what the client wants
type SearchCriteria struct {
	ID                   string // random ID for order, yet to define algorithm
	ClientName           string // Client Name
	ClientID             string // Client ID in case they logon
	Date                 string // Order Date
	Time                 string // Order Time
	Status               string // Open, Completed, Cancelled
	EatMode              string // EatIn, TakeAway, Delivery
	DeliveryMode         string // Internal, UberEats,
	DeliveryFee          string // Delivery Fee
	DeliveryLocation     string // Address
	DeliveryContactPhone string // Delivery phone number
}

// RespAddOrder is
type RespAddOrder struct {
	ID string
}

// FindAPI is to find stuff
func FindAPI(redisclient *redis.Client, orderFind string) order.Order {

	var apiserver string
	apiserver, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	urlrequest := apiserver + "/orderfind?orderid=" + orderFind

	url := fmt.Sprintf(urlrequest)

	var emptydisplay order.Order

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

	var orderback order.Order

	if err := json.NewDecoder(resp.Body).Decode(&orderback); err != nil {
		log.Println(err)
	}

	return orderback

}

// APICallList works
// Order List
func APICallList(redisclient *redis.Client) []order.Order {

	var apiserver string
	var emptydisplay []order.Order

	apiserver, _ = redisclient.Get("Web.APIServer.IPAddress").Result()
	urlrequest := apiserver + "/orderlist"

	url := fmt.Sprintf(urlrequest)

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

	// return list of orders
	var list []order.Order

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Println(err)
	}

	return list
}

// APICallListV2 works
// Order List
func APICallListV2(redisclient *redis.Client, credentials helper.Credentials) []order.Order {

	var apiserver string
	var emptydisplay []order.Order

	apiserver, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	urlrequest := apiserver + "/orderlist?clientid=" + credentials.UserID

	// Check if user is admin
	for x := 0; x < len(credentials.ClaimSet); x++ {
		if credentials.ClaimSet[x].Type == "USERTYPE" {
			if credentials.ClaimSet[x].Value == "ADMIN" {
				// list all if user is admin
				urlrequest = apiserver + "/orderlist"
				break
			}
		}
	}

	url := fmt.Sprintf(urlrequest)

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

	// return list of orders
	var list []order.Order

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Println(err)
	}

	return list
}

// APICallAdd is
func APICallAdd(redisclient *redis.Client, bodybyte []byte) RespAddOrder {

	envirvar := new(helper.RestEnvVariables)
	bodystr := string(bodybyte[:])

	envirvar.APIAPIServerIPAddress, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	// mongodbvar.APIServer = "http://localhost:1520/"

	apiURL := envirvar.APIAPIServerIPAddress
	resource := "/orderadd"

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(bodystr)
	resp2, err := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	var emptydisplay helper.Resultado
	emptydisplay.ErrorCode = resp2.Status

	defer resp2.Body.Close()
	var objectback RespAddOrder

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
		var resultado = resp2.Body
		log.Println(resultado)

		if err = json.NewDecoder(resp2.Body).Decode(&objectback); err != nil {
			log.Println(err)
		} else {

			var x = objectback.ID
			log.Println(x)
		}

	} else {
		emptydisplay.IsSuccessful = "N"

	}
	return objectback
}

// APICallFind is to find stuff
func APICallFind(redisclient *redis.Client, objectfind string) order.Order {

	var apiserver string
	apiserver, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	urlrequest := apiserver + "/orderfind?ID=" + objectfind

	url := fmt.Sprintf(urlrequest)

	var emptydisplay order.Order

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

	var objectback order.Order

	if err := json.NewDecoder(resp.Body).Decode(&objectback); err != nil {
		log.Println(err)
	}

	return objectback

}

// APICallUpdate is
func APICallUpdate(redisclient *redis.Client, bodybyte []byte) RespAddOrder {

	envirvar := new(helper.RestEnvVariables)
	bodystr := string(bodybyte[:])

	envirvar.APIAPIServerIPAddress, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	// mongodbvar.APIServer = "http://localhost:1520/"

	apiURL := envirvar.APIAPIServerIPAddress
	resource := "/orderupdate"

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(bodystr)
	resp2, err := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	var emptydisplay helper.Resultado
	emptydisplay.ErrorCode = resp2.Status

	defer resp2.Body.Close()
	var objectback RespAddOrder

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
		var resultado = resp2.Body
		log.Println(resultado)

		if err = json.NewDecoder(resp2.Body).Decode(&objectback); err != nil {
			log.Println(err)
		} else {

			var x = objectback.ID
			log.Println(x)
		}

	} else {
		emptydisplay.IsSuccessful = "N"

	}
	return objectback
}

// APICallDelete is
func APICallDelete(redisclient *redis.Client, dishUpdate Dish) helper.Resultado {

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

// APICallDeleteMany is
func APICallDeleteMany(redisclient *redis.Client, dishestodelete []string) helper.Resultado {

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

// Listdishes works
func Listdishes(redisclient *redis.Client) []dishes.Dish {

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
