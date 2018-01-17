// Package cachehandler to handle call to cache
// --------------------------------------------------------------
// .../src/restauranteweb/areas/cachehandler/cacheapicalls.go
// --------------------------------------------------------------
package cachehandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"restauranteweb/areas/helper"

	"github.com/go-redis/redis"
)

// Cache represents the cache data
type Cache struct {
	Key   string // cache key
	Value string // value in cache
}

// ListEntries works
func ListEntries(redisclient *redis.Client) []Cache {

	debug, _ := redisclient.Get("Web.Debug").Result()

	var apiserver string
	apiserver, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	urlrequest := apiserver + "/getcachedvalues"
	fmt.Println("urlrequest: ", urlrequest)

	url := fmt.Sprintf(urlrequest)

	var emptydisplay []Cache

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		if debug == "Y" {
			fmt.Println("NewRequest: ", err)
		}
		log.Fatal("NewRequest: ", err)

		return emptydisplay
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		if debug == "Y" {
			fmt.Println("client.Do(req): ", err)
		}
		log.Fatal("Do: ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	var cachelist []Cache

	if err := json.NewDecoder(resp.Body).Decode(&cachelist); err != nil {
		fmt.Println("json decoder error: ", err)
		log.Println(err)

	}

	return cachelist
}

// ListEntriesWeb works
func ListEntriesWeb(redisclient *redis.Client) []Cache {

	envirvar := new(helper.RestEnvVariables)
	envirvar.APIAPIServerIPAddress, _ = redisclient.Get("Web.APIServer.IPAddress").Result()
	envirvar.APIAPIServerPort, _ = redisclient.Get("Web.APIServer.Port").Result()
	envirvar.WEBServerPort, _ = redisclient.Get("WEBServerPort").Result()
	envirvar.WEBDebug, _ = redisclient.Get("Web.Debug").Result()
	envirvar.RecordCurrencyTick, _ = redisclient.Get("RecordCurrencyTick").Result()
	envirvar.RunningFromServer, _ = redisclient.Get("RunningFromServer").Result()

	cachelist := make([]Cache, 5)
	cachelist[0].Key = "Web.APIServer.IPAddress"
	cachelist[0].Value = envirvar.APIAPIServerIPAddress
	cachelist[1].Key = "Web.APIServer.Port"
	cachelist[1].Value = envirvar.APIAPIServerPort
	cachelist[2].Key = "RunningFromServer"
	cachelist[2].Value = envirvar.RunningFromServer
	cachelist[3].Key = "RecordCurrencyTick"
	cachelist[3].Value = envirvar.RecordCurrencyTick
	cachelist[4].Key = "WEBServerPort"
	cachelist[4].Value = envirvar.WEBServerPort

	return cachelist
}
