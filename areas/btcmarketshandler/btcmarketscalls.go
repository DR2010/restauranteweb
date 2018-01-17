// Package btcmarketshandler API calls for dishes web
// --------------------------------------------------------------
// .../src/restauranteweb/areas/btcmarkets/btcmarketscalls.go
// --------------------------------------------------------------
package btcmarketshandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"restauranteapi/helper"
	"strings"

	"github.com/go-redis/redis"
)

// BTCCoin is to be
type BTCCoin struct {
	Balance      int    // balance
	PendingFunds int    // pend
	Currency     string // curren
}

// BTCCoinR is to be
type BTCCoinR struct {
	Balance      string // balance
	PendingFunds string // pend
	Currency     string // curren
}

// CurrencyTick is to be
type CurrencyTick struct {
	Balance      int    // balance
	PendingFunds int    // pend
	Currency     string // curren
	LastPrice    string // decimal 1.45
	AUDCashNow   string // decimal 1.45
}

// BalanceCrypto e
type BalanceCrypto struct {
	Balance        string //
	Currency       string //
	CotacaoAtual   string //
	ValueInCashAUD string //
	BestBid        string
	BestAsk        string
	LastPrice      string
	Instrument     string
	Volume24       string
	DateTime       string
	Rotina         string
}

// ListCoinsIhave works
func ListCoinsIhave(redisclient *redis.Client) []BalanceCrypto {

	var apiserver string
	var emptydisplay []BalanceCrypto

	apiserver, _ = redisclient.Get("Web.APIServer.IPAddress").Result()
	urlrequest := apiserver + "/orderlist"

	// OLD urlrequest = "http://pontinhoapi.azurewebsites.net/api/btcmarkets/ALL"
	// NEW bypass portal urlrequest = "http://pontinhoapi.azurewebsites.net/api/btcmarkets/getbyid?id=ALL"

	// NEW via portal    urlrequest = https://danielapimanagement.azure-api.net/pontinhoapi.azurewebsites.net/api/btcmarkets/getbyid?id=ALL
	// Have to send the subscription key :-)
	// urlrequest = "https://danielapimanagement.azure-api.net/pontinhoapi.azurewebsites.net/api/btcmarkets/getbyid?id=ALL"
	urlrequest = "http://pontinhoapi.azurewebsites.net/api/btcmarkets/getbyid?id=ALL"

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
	var list []BalanceCrypto

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Println(err)
	}

	return list
}

// SendEmail works
func SendEmail(redisclient *redis.Client, emailstr string) {

	urlrequest := "http://pontinhoapi.azurewebsites.net/api/email/" + emailstr

	url := fmt.Sprintf(urlrequest)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	defer resp.Body.Close()

	return
}

// ListCoinsHistory works
func ListCoinsHistory(redisclient *redis.Client, currency string, rows string) []BalanceCrypto {

	var apiserver string
	apiserver, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	urlrequest := apiserver + "/btccotacaolist?currency=" + currency + "&rows=" + rows

	// urlrequest = "http://localhost:1520/btccotacaolist?currency=ALL&rows=50"

	url := fmt.Sprintf(urlrequest)

	var emptydisplay []BalanceCrypto

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

	var dishlist []BalanceCrypto

	if err := json.NewDecoder(resp.Body).Decode(&dishlist); err != nil {
		log.Println(err)
	}

	return dishlist
}

// ListCoinsHistoryDate works
func ListCoinsHistoryDate(redisclient *redis.Client, currency string, yeardaymonth string, yeardaymonthend string) []BalanceCrypto {

	var apiserver string
	apiserver, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	urlrequest := apiserver + "/btccotacaolistdate?currency=" + currency + "&yeardaymonth=" + yeardaymonth + "&yeardaymonthend=" + yeardaymonthend

	// urlrequest = "http://localhost:1520/btccotacaolist?currency=ALL&rows=50"

	url := fmt.Sprintf(urlrequest)

	var emptydisplay []BalanceCrypto

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

	var dishlist []BalanceCrypto

	if err := json.NewDecoder(resp.Body).Decode(&dishlist); err != nil {
		log.Println(err)
	}

	return dishlist
}

// APIcallAdd is based on Dishes Add - different from Order Add
// They are different because the order one handles the entire form - is manually handled
func APIcallAdd(redisclient *redis.Client, cryptoInsert BalanceCrypto, rotina string) helper.Resultado {

	envirvar := new(helper.RestEnvVariables)

	envirvar.APIAPIServerIPAddress, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	apiURL := envirvar.APIAPIServerIPAddress
	resource := "/btccotacaoadd"

	data := url.Values{}
	data.Add("cryptoBalance", cryptoInsert.Balance)
	data.Add("cryptoCotacaoAtual", cryptoInsert.CotacaoAtual)
	data.Add("cryptoCurrency", cryptoInsert.Currency)
	data.Add("cryptoValueInCashAUD", cryptoInsert.ValueInCashAUD)
	data.Add("cryptoBestAsk", cryptoInsert.BestAsk)
	data.Add("cryptoBestBid", cryptoInsert.BestBid)
	data.Add("cryptoVolume24", cryptoInsert.Volume24)
	data.Add("rotina", cryptoInsert.Rotina)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(data.Encode())
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	var emptydisplay helper.Resultado
	emptydisplay.ErrorCode = resp2.Status

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
	} else {
		emptydisplay.IsSuccessful = "N"
	}

	return emptydisplay

}

// PreOrderAPICallList works
// Order List
func PreOrderAPICallList(redisclient *redis.Client) []PreOrder {

	var apiserver string
	var emptydisplay []PreOrder

	apiserver, _ = redisclient.Get("Web.APIServer.IPAddress").Result()
	urlrequest := apiserver + "/btcpreorderlist"

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
	var list []PreOrder

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Println(err)
	}

	return list
}

// RespAddOrder is
type RespAddOrder struct {
	ID string
}

// PreOrderAPIcallAdd is based on Dishes Add - different from Order Add
// They are different because the order one handles the entire form - is manually handled
func PreOrderAPICallAdd(redisclient *redis.Client, bodybyte []byte) RespAddOrder {

	envirvar := new(helper.RestEnvVariables)
	bodystr := string(bodybyte[:])

	envirvar.APIAPIServerIPAddress, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	// mongodbvar.APIServer = "http://localhost:1520/"

	apiURL := envirvar.APIAPIServerIPAddress
	resource := "/btcpreorderadd"

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
