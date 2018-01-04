package security

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"mongodb/helper"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

// DishupdateAPI is
// func DishupdateAPI(redisclient *redis.Client, dishUpdate Dish) helper.Resultado {

func LoginUser(redisclient *redis.Client, userid string, password string) helper.Resultado {

	mongodbvar := new(helper.DatabaseX)

	mongodbvar.APIServer, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	apiURL := mongodbvar.APIServer
	resource := "/securitylogin"

	data := url.Values{}
	data.Add("userid", userid)
	data.Add("password", password)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(data.Encode())
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	fmt.Println("resp2.Status:" + resp2.Status)

	var emptydisplay helper.Resultado
	emptydisplay.ErrorCode = resp2.Status

	defer resp2.Body.Close()

	var response string

	if err := json.NewDecoder(resp2.Body).Decode(&response); err != nil {
		log.Println(err)
	}

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
		emptydisplay.ErrorCode = "200 OK"
		emptydisplay.ErrorDescription = "200 OK"
		emptydisplay.ReturnedValue = response

		// Find out what to do to get a value back
		//
		tokenreturned := response

		// Store Token in Cache
		var jwttoken = tokenreturned
		year, month, day := time.Now().Date()

		var key = userid + strconv.Itoa(int(year)) + strconv.Itoa(int(month)) + strconv.Itoa(int(day))

		_ = redisclient.Set(key, jwttoken, 0).Err()

	} else {
		emptydisplay.IsSuccessful = "N"
		emptydisplay.ErrorCode = "404 Error"
		emptydisplay.ErrorDescription = "404 Shit happens!... and it happened!"
	}

	return emptydisplay
}

func SignUp(redisclient *redis.Client, userid string, password string, passwordvalidate string) helper.Resultado {

	mongodbvar := new(helper.DatabaseX)
	mongodbvar.APIServer, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	var emptydisplay helper.Resultado

	apiURL := mongodbvar.APIServer
	resource := "/securitysignup"

	if password != passwordvalidate {
		emptydisplay.ErrorCode = "404 Error"
		emptydisplay.ErrorDescription = "Password mismatch"
		return emptydisplay
	}

	var passwordhashed = Hashstring(password)
	var passwordvalidatehashed = Hashstring(passwordvalidate)

	// passwordhashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// passwordvalidatehashed, _ := bcrypt.GenerateFromPassword([]byte(passwordvalidate), bcrypt.DefaultCost)

	// passwordhasheds := string(passwordhashed)
	// passwordvalidatehasheds := string(passwordvalidatehashed)

	data := url.Values{}
	data.Add("userid", userid)
	data.Add("password", passwordhashed)
	data.Add("passwordvalidate", passwordvalidatehashed)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	body := strings.NewReader(data.Encode())

	// Call method here
	resp2, _ := http.Post(urlStr, "application/x-www-form-urlencoded", body)

	fmt.Println("resp2.Status:" + resp2.Status)

	emptydisplay.ErrorCode = resp2.Status

	if resp2.Status == "200 OK" {
		emptydisplay.IsSuccessful = "Y"
	}

	return emptydisplay

}

// ValidateToken is half way
func ValidateToken(redisclient *redis.Client, userid string, token string) string {

	//Get  Token in Cache
	year, month, day := time.Now().Date()

	var key = userid + strconv.Itoa(int(year)) + strconv.Itoa(int(month)) + strconv.Itoa(int(day))

	tokenstored, _ := redisclient.Get(key).Result()

	var ret = "NotOkToLogin"
	if tokenstored == token {
		ret = "OkToLogin"
	}

	return ret
}

// this is just a reference key
// the roles, date and user will be stored at the server
func Hashstring(str string) string {

	s := str
	h := sha1.New()
	h.Write([]byte(s))

	sha1hash := hex.EncodeToString(h.Sum(nil))

	return sha1hash
}
