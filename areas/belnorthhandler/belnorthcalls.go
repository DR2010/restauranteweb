// Package belnorthhandler API calls for dishes web
// --------------------------------------------------------------
// .../src/areas/belnorthhandler/belnorthcalls.go
// --------------------------------------------------------------
package belnorthhandler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis"
)

// Player is player
type Player struct {
	Ffanumber         string
	Firstname         string
	Lastname          string
	Dateofbirth       string
	Gender            string
	Display           string
	Trialagegroup     string
	Trialopengirls    string
	Haveyouregistered string
	Igivepermission   string
	BIBnumber         string
	Emailaddress      string
	PosBack           string
	PosMidfield       string
	PosAttack         string
	PosKeeper         string
	Agegroupdob       string
	Mobile            string
	Olderinterested   string
	Olderagetrial     string
	ShirtSize         string
	RegisteredWithCF  string
	TransAmount       string
	TransReference    string
}

// Payment is player
type Payment struct {
	FKCompetition  string
	FFA            string
	Name           string
	Date           string
	CCFName        string
	TransactionRef string
	Reference      string
	DisbursementID string
	Amount         string
	CCNumber       string
	CardHolderName string
	Club           string
}

// ListGradingPlayers works
func ListGradingPlayers(redisclient *redis.Client) []Player {

	var apiserver string
	var emptydisplay []Player

	apiserver, _ = redisclient.Get("Web.APIServer.IPAddress").Result()
	urlrequest := apiserver + "/orderlist"

	urlrequest = "http://belnorth.com/api/bnapiplayerrepo.php?action=getPlayersGrading"

	url := fmt.Sprintf(urlrequest)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return emptydisplay
	}

	// via portal requires header
	// req.Header.Set("Ocp-Apim-Subscription-Key", "eb9f7b1620494fb2bdc7815705fd8c7e")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	// return list of crypto
	var list []Player

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Println(err)
	}

	return list
}

// ListCompetitionPlayers works
func ListCompetitionPlayers(redisclient *redis.Client, competition string) []Player {

	var apiserver string
	var emptydisplay []Player

	apiserver, _ = redisclient.Get("Web.APIServer.IPAddress").Result()
	urlrequest := apiserver + "/orderlist"

	urlrequest = "http://belnorth.com/api/bnapiplayerrepo.php?action=ListCompetitionPlayers&competition=" + competition

	url := fmt.Sprintf(urlrequest)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return emptydisplay
	}

	// via portal requires header
	// req.Header.Set("Ocp-Apim-Subscription-Key", "eb9f7b1620494fb2bdc7815705fd8c7e")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	// return list of crypto
	var list []Player

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Println(err)
	}

	return list
}

// ListPayments works
func ListPayments(redisclient *redis.Client) []Payment {

	var apiserver string
	var emptydisplay []Payment

	apiserver, _ = redisclient.Get("Web.APIServer.IPAddress").Result()
	urlrequest := apiserver + "/orderlist"

	urlrequest = "http://belnorth.com/api/bnapiplayerrepo.php?action=payments"

	url := fmt.Sprintf(urlrequest)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return emptydisplay
	}

	// via portal requires header
	// req.Header.Set("Ocp-Apim-Subscription-Key", "eb9f7b1620494fb2bdc7815705fd8c7e")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	// return list of crypto
	var list []Payment

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Println(err)
	}

	return list
}

// GetSinglePayment works
func GetSinglePayment(redisclient *redis.Client, ffanumber string) Payment {

	var apiserver string
	var emptydisplay Payment

	apiserver, _ = redisclient.Get("Web.APIServer.IPAddress").Result()
	urlrequest := apiserver + "/orderlist"

	urlrequest = "http://belnorth.com/api/bnapiplayerrepo.php?action=getpaymentdetails&ffanumber=" + ffanumber

	url := fmt.Sprintf(urlrequest)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return emptydisplay
	}

	// via portal requires header
	// req.Header.Set("Ocp-Apim-Subscription-Key", "eb9f7b1620494fb2bdc7815705fd8c7e")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return emptydisplay
	}

	defer resp.Body.Close()

	// return list of crypto
	var payment Payment

	if err := json.NewDecoder(resp.Body).Decode(&payment); err != nil {
		log.Println(err)
	}

	return payment
}

// PlayerRegistrationFile is
type PlayerRegistrationFile struct {
	FFA  string
	Name string
	DOB  string
}

// Capitalfootball is
func Capitalfootball(redisclient *redis.Client) []PlayerRegistrationFile {

	file, err := os.Open("capitalfootball.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var playerlist []PlayerRegistrationFile

	scanner := bufio.NewScanner(file)

	playerlist = make([]PlayerRegistrationFile, 52)

	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(scanner.Text())

		tmp := strings.Split(line, ",")

		i++
		playerlist[i] = PlayerRegistrationFile{}
		playerlist[i].FFA = strings.Trim(tmp[0], " ")
		playerlist[i].Name = strings.Trim(tmp[1], " ")
		playerlist[i].DOB = strings.Trim(tmp[2], " ")

		fmt.Println(playerlist[i].FFA)

		err = redisclient.Set(playerlist[i].FFA, playerlist[i].Name, 0).Err()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return playerlist
}
