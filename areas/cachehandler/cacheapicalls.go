package cachehandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
