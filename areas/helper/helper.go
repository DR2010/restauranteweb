package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// DatabaseX is a struct
type DatabaseX struct {
	Location   string // location of the database localhost, something.com, etc
	Database   string // database name
	Collection string // collection name
	APIServer  string // apiserver name
}

// Resultado is a struct
type Resultado struct {
	ErrorCode        string // error code
	ErrorDescription string // description
	IsSuccessful     string // Y or N
}

func add() {
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// DatabaseX is a struct
type RestEnvVariables struct {
	APIMongoDBLocation    string // location of the database localhost, something.com, etc
	APIMongoDBDatabase    string // database name
	APIAPIServerPort      string // collection name
	APIAPIServerIPAddress string // apiserver name
	WEBDebug              string // debug
}

// Readfileintostruct is
func Readfileintostruct() RestEnvVariables {
	dat, err := ioutil.ReadFile("restaurante.ini")
	check(err)
	fmt.Print(string(dat))

	var restenv RestEnvVariables

	json.Unmarshal(dat, &restenv)

	return restenv
}
