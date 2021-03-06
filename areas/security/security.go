package security

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	helper "festajuninaweb/areas/helper"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"restauranteapi/security"
	"strings"

	"github.com/go-redis/redis"
)

// LoginUserV2 but not
func LoginUserV2(redisclient *redis.Client, userid string, password string) security.Credentials {

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

	var response security.Credentials

	if err := json.NewDecoder(resp2.Body).Decode(&response); err != nil {
		log.Println(err)
	}

	if resp2.Status == "200 OK" {
		return response
	}

	response.ApplicationID = "None"
	response.JWT = "Error"

	return response
}

// LoginUser something
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

	} else {
		emptydisplay.IsSuccessful = "N"
		emptydisplay.ErrorCode = "404 Error"
		emptydisplay.ErrorDescription = "404 Shit happens!... and it happened!"

	}

	return emptydisplay
}

// SignUp function
func SignUp(redisclient *redis.Client, userid string, preferredname string, password string, passwordvalidate string, applicationid string) helper.Resultado {

	mongodbvar := new(helper.DatabaseX)
	mongodbvar.APIServer, _ = redisclient.Get("Web.APIServer.IPAddress").Result()

	var emptydisplay helper.Resultado

	apiURL := mongodbvar.APIServer
	resource := "/securitysignup"

	if userid == "" {
		emptydisplay.ErrorCode = "404 Error"
		emptydisplay.ErrorDescription = "user id not suppplied"
		return emptydisplay
	}

	if password == "" {
		emptydisplay.ErrorCode = "404 Error"
		emptydisplay.ErrorDescription = "password not suppplied"
		return emptydisplay
	}

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
	data.Add("preferredname", preferredname)
	data.Add("password", passwordhashed)
	data.Add("passwordvalidate", passwordvalidatehashed)
	data.Add("applicationid", applicationid)

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
func ValidateToken(redisclient *redis.Client, httprequest *http.Request) string {

	var credtemp helper.Credentials

	cookie, _ := httprequest.Cookie("DanBTCjwt")
	if cookie == nil {
		return "NotOkToLogin"
	}

	cookieinbytes := []byte(cookie.Value)
	_ = json.Unmarshal(cookieinbytes, &credtemp)

	var key = credtemp.KeyJWT

	tokenstored, _ := redisclient.Get(key).Result()

	var ret = "NotOkToLogin"
	if tokenstored == credtemp.JWT {
		ret = "OkToLogin"
	}

	return ret
}

// ValidateTokenV2 will get info from cache
func ValidateTokenV2(redisclient *redis.Client, httprequest *http.Request) (string, helper.Credentials) {
	var credentials helper.Credentials

	credentials.ApplicationID = "Restaurante"
	credentials.UserID = "Anonymous"
	credentials.UserName = "Anonymous"
	credentials.IsAdmin = "No"
	credentials.IsAnonymous = "Yes"

	// The system will store an object in cache and the key must be the used ID
	// The same user can logon in 2 places, I think
	// Users can't be mixed, I can't trust the variables since it is completely stateless - each request is stateless

	// Machine credentials
	//
	clientsecret := httprequest.FormValue("macdantoken")

	if clientsecret != "" {
		// Issue keys - should be stored in the database API Key or Secret I think
		//
		if clientsecret == "BypassSecurity" {
			var credentialsmachine helper.Credentials
			credentialsmachine.ApplicationID = "Restaurante"
			credentialsmachine.UserID = "Machine"
			credentialsmachine.JWT = clientsecret
			return "OkToLogin", credentialsmachine
		}
	}

	credentials.JWT = "Error"

	jwtincookie := ""
	useridincookie := ""

	cookiekeyJWT := "DanBTCjwt"
	cookiekeyUSERID := "DanBTCuserid"

	cookieJWT, err := httprequest.Cookie(cookiekeyJWT)

	if err != nil {
		log.Println(err)
		log.Println("Not found Cookie: " + cookiekeyJWT)
	}

	if cookieJWT == nil {
		return "NotOkToLogin", credentials
	}

	cookieUSERID, err2 := httprequest.Cookie(cookiekeyUSERID)
	if err2 != nil {
		log.Println(err2)
		log.Println("Not found Cookie: " + cookiekeyUSERID)
	}

	if cookieUSERID == nil {
		return "NotOkToLogin", credentials
	}

	jwtincookie = cookieJWT.Value
	useridincookie = cookieUSERID.Value

	var keyredis = cookiekeyJWT + useridincookie

	tokenstored, _ := redisclient.Get(keyredis).Result()
	tokenstoredbytes := []byte(tokenstored)

	_ = json.Unmarshal(tokenstoredbytes, &credentials)

	var ret = "NotOkToLogin"
	if credentials.JWT == jwtincookie {
		credentials.IsAnonymous = "No"
		ret = "OkToLogin"
	} else {
		credentials.ApplicationID = "Restaurante"
		credentials.UserID = "Anonymous"
		credentials.UserName = "Anonymous"
		credentials.IsAdmin = "No"
	}

	return ret, credentials
}

// Hashstring is just for hashing - only reference key
func Hashstring(str string) string {

	s := str
	h := sha1.New()
	h.Write([]byte(s))

	sha1hash := hex.EncodeToString(h.Sum(nil))

	return sha1hash
}
