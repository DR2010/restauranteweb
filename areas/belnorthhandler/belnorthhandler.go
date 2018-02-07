// Package belnorthhandler Handler for dishes web
// -----------------------------------------------------------
// .../src/restauranteweb/areas/belnorthhandler/belnorthhandler.go
// -----------------------------------------------------------
package belnorthhandler

import (
	"html/template"
	"net/http"
	helper "restauranteweb/areas/helper"

	"github.com/go-redis/redis"
)

// This is the template to display as part of the html template
//

// ControllerInfo is
type ControllerInfo struct {
	UserID      string
	Name        string
	Message     string
	Currency    string
	FromDate    string
	ToDate      string
	Application string
}

// Row is
type Row struct {
	Description []string
}

// Coin is
type Coin struct {
	Short string
	Name  string
}

// SecurityClaim is
type SecurityClaim struct {
	Type  string
	Value string
}

// DisplayTemplate is
type DisplayTemplate struct {
	Info       ControllerInfo
	FieldNames []string
	Rows       []Row
	Players    []Player
	Payments   []Payment
}

var mongodbvar helper.DatabaseX

// HListGradingPlayers get list of grading players
//
func HListGradingPlayers(httpwriter http.ResponseWriter, redisclient *redis.Client, credentials helper.Credentials, ListOfPayments []Payment) {

	// create new template
	// t, _ := template.ParseFiles("templates/btcmarkets/btcbasic.html", "templates/btcmarkets/btcmarketstemplate.html")
	t, _ := template.ParseFiles("html/homepage.html", "templates/belnorth/pagebodytemplateBelnorth.html")

	// Get list of orders (api call)
	//
	var list = ListGradingPlayers(redisclient)

	// Assembly display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Coins"
	items.Info.Currency = "SUMMARY"
	items.Info.UserID = credentials.UserID
	items.Info.Application = credentials.ApplicationID

	items.Players = make([]Player, len(list))

	for i := 0; i < len(list); i++ {

		items.Players[i] = Player{}
		items.Players[i].Ffanumber = list[i].Ffanumber
		items.Players[i].Firstname = list[i].Firstname
		items.Players[i].Lastname = list[i].Lastname
		items.Players[i].Agegroupdob = list[i].Agegroupdob
		items.Players[i].Dateofbirth = list[i].Dateofbirth
		items.Players[i].Gender = list[i].Gender
		items.Players[i].Trialagegroup = list[i].Trialagegroup
		items.Players[i].Trialopengirls = list[i].Trialopengirls
		items.Players[i].Haveyouregistered = list[i].Haveyouregistered
		items.Players[i].Igivepermission = list[i].Igivepermission
		items.Players[i].BIBnumber = list[i].BIBnumber
		items.Players[i].Emailaddress = list[i].Emailaddress
		items.Players[i].PosBack = list[i].PosBack
		items.Players[i].PosMidfield = list[i].PosMidfield
		items.Players[i].PosAttack = list[i].PosAttack
		items.Players[i].PosKeeper = list[i].PosKeeper
		items.Players[i].PosMidfield = list[i].PosMidfield
		items.Players[i].Mobile = list[i].Mobile
		items.Players[i].Olderinterested = list[i].Olderinterested
		items.Players[i].Olderagetrial = list[i].Olderagetrial

		// playernamecf, _ := redisclient.Get(list[i].Ffanumber).Result()

		// if playernamecf != "" {
		// 	items.Players[i].RegisteredWithCF = "Y"
		// }

		for y := range ListOfPayments {
			if ListOfPayments[y].FFA == items.Players[i].Ffanumber {
				// Found!
				items.Players[i].RegisteredWithCF = "Y"
				items.Players[i].TransAmount = ListOfPayments[y].Amount
				items.Players[i].TransReference = ListOfPayments[y].Reference
				break
			}
		}
		// fmt.Println(items.Players[i].Ffanumber + " " + items.Players[i].RegisteredWithCF + " " + playernamecf)

	}

	t.Execute(httpwriter, items)

	return
}

// HListCompetitionPlayers get list of grading players
//
func HListCompetitionPlayers(httpwriter http.ResponseWriter, redisclient *redis.Client, credentials helper.Credentials, ListOfPayments []Payment) {

	// create new template
	// t, _ := template.ParseFiles("templates/btcmarkets/btcbasic.html", "templates/btcmarkets/btcmarketstemplate.html")
	t, _ := template.ParseFiles("html/homepage.html", "templates/belnorth/pagebodytbelnorthcompetition.html")

	var competition = "BELNORTHPRESEASON2018"

	// Get list of orders (api call)
	//
	var list = ListCompetitionPlayers(redisclient, competition)

	// Assembly display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Coins"
	items.Info.Currency = "SUMMARY"
	items.Info.UserID = credentials.UserID
	items.Info.Application = credentials.ApplicationID

	items.Players = make([]Player, len(list))

	for i := 0; i < len(list); i++ {

		items.Players[i] = Player{}
		items.Players[i].Ffanumber = list[i].Ffanumber
		items.Players[i].Firstname = list[i].Firstname
		items.Players[i].Lastname = list[i].Lastname
		items.Players[i].Agegroupdob = list[i].Agegroupdob
		items.Players[i].Dateofbirth = list[i].Dateofbirth
		items.Players[i].Gender = list[i].Gender
		items.Players[i].Haveyouregistered = list[i].Haveyouregistered
		items.Players[i].Igivepermission = list[i].Igivepermission
		items.Players[i].Emailaddress = list[i].Emailaddress
		items.Players[i].Mobile = list[i].Mobile
		items.Players[i].ShirtSize = list[i].ShirtSize

		// playernamecf, _ := redisclient.Get(list[i].Ffanumber).Result()

		// Call API is too slow
		// paymentinfo := GetSinglePayment(redisclient, items.Players[i].Ffanumber)

		// find in array
		//

		for y := range ListOfPayments {
			if ListOfPayments[y].FFA == items.Players[i].Ffanumber {
				// Found!
				items.Players[i].RegisteredWithCF = "Y"
				items.Players[i].TransAmount = ListOfPayments[y].Amount
				items.Players[i].TransReference = ListOfPayments[y].Reference
				break
			}
		}

		// if paymentinfo.Name != "" {
		// 	items.Players[i].RegisteredWithCF = paymentinfo.Amount
		// }

	}

	t.Execute(httpwriter, items)

	return
}

// HListPayments get list of grading players
//
func HListPayments(httpwriter http.ResponseWriter, redisclient *redis.Client, credentials helper.Credentials) {

	// create new template
	// t, _ := template.ParseFiles("templates/btcmarkets/btcbasic.html", "templates/btcmarkets/btcmarketstemplate.html")
	t, _ := template.ParseFiles("html/homepage.html", "templates/belnorth/payments.html")

	// Get list of orders (api call)
	//
	var list = ListPayments(redisclient)

	// Assembly display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Payments"
	items.Info.Currency = "SUMMARY"
	items.Info.UserID = credentials.UserID
	items.Info.Application = credentials.ApplicationID

	items.Payments = make([]Payment, len(list))

	for i := 0; i < len(list); i++ {

		items.Payments[i] = Payment{}
		items.Payments[i].Amount = list[i].Amount
		items.Payments[i].CardHolderName = list[i].CardHolderName
		items.Payments[i].CCFName = list[i].CCFName
		items.Payments[i].CCNumber = list[i].CCNumber
		items.Payments[i].Club = list[i].Club
		items.Payments[i].Date = list[i].Date
		items.Payments[i].DisbursementID = list[i].DisbursementID
		items.Payments[i].FFA = list[i].FFA
		items.Payments[i].FKCompetition = list[i].FKCompetition
		items.Payments[i].Name = list[i].Name
		items.Payments[i].Reference = list[i].Reference
		items.Payments[i].TransactionRef = list[i].TransactionRef

	}

	t.Execute(httpwriter, items)

	return
}

// Lpad is left pad
func Lpad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}
