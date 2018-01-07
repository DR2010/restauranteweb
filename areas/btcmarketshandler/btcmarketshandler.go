// Package btcmarketshandler Handler for dishes web
// -----------------------------------------------------------
// .../src/restauranteweb/areas/disherhandler/ordershandler.go
// -----------------------------------------------------------
package btcmarketshandler

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
	Name     string
	Message  string
	Currency string
}

// Row is
type Row struct {
	Description []string
}

// DisplayTemplate is
type DisplayTemplate struct {
	Info       ControllerInfo
	FieldNames []string
	Rows       []Row
	Btccoin    []BalanceCrypto
}

var mongodbvar helper.DatabaseX

// List = assemble results of API call to dish list
//
func List(httpwriter http.ResponseWriter, redisclient *redis.Client) {

	// create new template
	t, _ := template.ParseFiles("templates/basictemplate.html", "templates/btcmarkets/btcmarketstemplate.html")

	// Get list of orders (api call)
	//
	var list = ListCoinsIhave(redisclient)

	// Assembly display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Daniel Investment List"
	items.Info.Currency = "NA"

	var numberoffields = 4

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Currency"
	items.FieldNames[1] = "Balance"
	items.FieldNames[2] = "CotacaoAtual"
	items.FieldNames[3] = "ValueInCashAUD"

	// Set rows to be displayed
	items.Rows = make([]Row, len(list))
	items.Btccoin = make([]BalanceCrypto, len(list))

	for i := 0; i < len(list); i++ {

		items.Btccoin[i] = BalanceCrypto{}
		items.Btccoin[i].Balance = list[i].Balance
		items.Btccoin[i].Currency = list[i].Currency
		items.Btccoin[i].CotacaoAtual = list[i].CotacaoAtual
		items.Btccoin[i].ValueInCashAUD = list[i].ValueInCashAUD

	}

	t.Execute(httpwriter, items)

}

// ListV2 = assemble results of API call to dish list
//
func ListV2(httpwriter http.ResponseWriter, redisclient *redis.Client) []BalanceCrypto {

	// create new template
	t, _ := template.ParseFiles("templates/basictemplate.html", "templates/btcmarkets/btcmarketstemplate.html")

	// Get list of orders (api call)
	//
	var list = ListCoinsIhave(redisclient)

	// Assembly display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Coins"
	items.Info.Currency = "SUMMARY"

	var numberoffields = 4

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Currency"
	items.FieldNames[1] = "Balance"
	items.FieldNames[2] = "CotacaoAtual"
	items.FieldNames[3] = "ValueInCashAUD"

	// Set rows to be displayed
	items.Rows = make([]Row, len(list))
	items.Btccoin = make([]BalanceCrypto, len(list))

	var RetBtccoin []BalanceCrypto
	RetBtccoin = make([]BalanceCrypto, len(list))

	for i := 0; i < len(list); i++ {

		items.Btccoin[i] = BalanceCrypto{}
		items.Btccoin[i].Balance = list[i].Balance
		items.Btccoin[i].Currency = list[i].Currency
		items.Btccoin[i].CotacaoAtual = list[i].CotacaoAtual
		items.Btccoin[i].ValueInCashAUD = list[i].ValueInCashAUD

		// New code to return values to write to mongo every minute or every call
		// 31/12/2017
		//
		RetBtccoin[i] = BalanceCrypto{}
		RetBtccoin[i].Balance = list[i].Balance
		RetBtccoin[i].Currency = list[i].Currency
		RetBtccoin[i].CotacaoAtual = list[i].CotacaoAtual
		RetBtccoin[i].ValueInCashAUD = list[i].ValueInCashAUD

	}

	t.Execute(httpwriter, items)

	return RetBtccoin
}

// GetBalance = assemble results of API call to dish list
//
func GetBalance(redisclient *redis.Client) []BalanceCrypto {

	// Get list of orders (api call)
	//
	var list = ListCoinsIhave(redisclient)

	var RetBtccoin []BalanceCrypto
	RetBtccoin = make([]BalanceCrypto, len(list))

	for i := 0; i < len(list); i++ {

		// New code to return values to write to mongo every minute or every call
		// 31/12/2017
		//
		RetBtccoin[i] = BalanceCrypto{}
		RetBtccoin[i].Balance = list[i].Balance
		RetBtccoin[i].Currency = list[i].Currency
		RetBtccoin[i].CotacaoAtual = list[i].CotacaoAtual
		RetBtccoin[i].ValueInCashAUD = list[i].ValueInCashAUD

	}

	return RetBtccoin
}

// HListHistory = assemble results of API call to
//
func HListHistory(httpwriter http.ResponseWriter, redisclient *redis.Client, currency string, rows string) []BalanceCrypto {

	// create new template
	t, _ := template.ParseFiles("templates/basictemplate.html", "templates/btcmarkets/btcmarketstemplate.html")

	// Get list of orders (api call)
	//
	var list = ListCoinsHistory(redisclient, currency, rows)

	// Assembly display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Market Value - History"
	items.Info.Currency = currency

	var numberoffields = 5

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Currency"
	items.FieldNames[1] = "Balance"
	items.FieldNames[2] = "CotacaoAtual"
	items.FieldNames[3] = "ValueInCashAUD"
	items.FieldNames[4] = "DateTime"

	// Set rows to be displayed
	items.Rows = make([]Row, len(list))
	items.Btccoin = make([]BalanceCrypto, len(list))

	var RetBtccoin []BalanceCrypto
	RetBtccoin = make([]BalanceCrypto, len(list))

	for i := 0; i < len(list); i++ {

		items.Btccoin[i] = BalanceCrypto{}
		items.Btccoin[i].Balance = list[i].Balance
		items.Btccoin[i].Currency = list[i].Currency
		items.Btccoin[i].CotacaoAtual = list[i].CotacaoAtual
		items.Btccoin[i].ValueInCashAUD = list[i].ValueInCashAUD
		items.Btccoin[i].DateTime = list[i].DateTime

	}

	t.Execute(httpwriter, items)

	return RetBtccoin
}

// HListHistoryDate = assemble results of API call to
//
func HListHistoryDate(httpwriter http.ResponseWriter, redisclient *redis.Client, currency string, yeardaymonth string, yeardaymonthend string) []BalanceCrypto {

	// create new template
	t, _ := template.ParseFiles("templates/basictemplate.html", "templates/btcmarkets/btcmarketstemplate.html")

	// Get list of orders (api call)
	//
	var list = ListCoinsHistoryDate(redisclient, currency, yeardaymonth, yeardaymonthend)

	// Assembly display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Market Value - History - Date"
	items.Info.Currency = currency

	var numberoffields = 5

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Currency"
	items.FieldNames[1] = "Balance"
	items.FieldNames[2] = "CotacaoAtual"
	items.FieldNames[3] = "ValueInCashAUD"
	items.FieldNames[4] = "DateTime"

	// Set rows to be displayed
	items.Rows = make([]Row, len(list))
	items.Btccoin = make([]BalanceCrypto, len(list))

	var RetBtccoin []BalanceCrypto
	RetBtccoin = make([]BalanceCrypto, len(list))

	for i := 0; i < len(list); i++ {

		items.Btccoin[i] = BalanceCrypto{}
		items.Btccoin[i].Balance = list[i].Balance
		items.Btccoin[i].Currency = list[i].Currency
		items.Btccoin[i].CotacaoAtual = list[i].CotacaoAtual
		items.Btccoin[i].ValueInCashAUD = list[i].ValueInCashAUD
		items.Btccoin[i].DateTime = list[i].DateTime

	}

	t.Execute(httpwriter, items)

	return RetBtccoin
}

// RecordTick is xxx
func RecordTick(balcrypto []BalanceCrypto, redisclient *redis.Client) {

	for i := 0; i < len(balcrypto); i++ {

		bcc := BalanceCrypto{}
		bcc.Balance = balcrypto[i].Balance
		bcc.Currency = balcrypto[i].Currency
		bcc.CotacaoAtual = balcrypto[i].CotacaoAtual
		bcc.ValueInCashAUD = balcrypto[i].ValueInCashAUD

		APIcallAdd(redisclient, bcc)
	}

	return
}

// Lpad is left pad
func Lpad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}
