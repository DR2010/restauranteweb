package security

import (
	"encoding/json"
	helper "festajuninaweb/areas/helper"
	"log"
	"mongodb/dishes"
	"net/http"
	"restauranteapi/security"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/template"
	"github.com/go-redis/redis"
)

// SignupPage is for the user to signup
func SignupPage(httpresponsewriter http.ResponseWriter, req *http.Request, redisclient *redis.Client) {

	type ControllerInfo struct {
		Name    string
		Message string
	}
	type DisplayTemplate struct {
		Info ControllerInfo
	}

	items := DisplayTemplate{}
	items.Info.Name = "Login Page"

	if req.Method != "POST" {

		// t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
		t, _ := template.ParseFiles("templates/security/signupheader.html", "templates/security/signupdetail.html")
		items.Info.Message = ""
		t.Execute(httpresponsewriter, items)

		// http.ServeFile(res, req, "templates/security/signup.html")
		return
	}

	usernamemix := req.FormValue("username")
	preferredname := req.FormValue("preferredname")
	password := req.FormValue("password")
	passwordvalidate := req.FormValue("passwordvalidate")
	// applicationid := req.FormValue("applicationid")
	applicationid := "Restaurante" // Festa Junina

	username := strings.ToUpper(usernamemix)

	if username == "" {
		// t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
		t, _ := template.ParseFiles("templates/security/signupheader.html", "templates/security/signupdetail.html")
		items.Info.Message = "Please enter Name."
		t.Execute(httpresponsewriter, items)
		return
	}

	if preferredname == "" {
		// t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
		t, _ := template.ParseFiles("templates/security/signupheader.html", "templates/security/signupdetail.html")
		items.Info.Message = "Please enter Preferred Name."
		t.Execute(httpresponsewriter, items)
		return
	}

	if password == "" {
		// t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
		t, _ := template.ParseFiles("templates/security/signupheader.html", "templates/security/signupdetail.html")
		items.Info.Message = "Please enter details."
		t.Execute(httpresponsewriter, items)
		return
	}

	if password == passwordvalidate {

		// Call API to check if user exists and create
		var resultado = SignUp(redisclient, username, preferredname, password, passwordvalidate, applicationid)
		if resultado.ErrorCode == "200 OK" {

		} else {
			// t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
			t, _ := template.ParseFiles("templates/security/signupheader.html", "templates/security/signupdetail.html")
			items.Info.Message = "Passwords mismatch."
			t.Execute(httpresponsewriter, items)
			return
		}

		http.Redirect(httpresponsewriter, req, "/", 303)
	} else {
		// t, _ := template.ParseFiles("templates/security/signup.html", "templates/security/loginmessagetemplate.html")
		t, _ := template.ParseFiles("templates/security/signupheader.html", "templates/security/signupdetail.html")
		items.Info.Message = "Passwords do not match."
		t.Execute(httpresponsewriter, items)
		return
	}

}

// LogoutPage is for the user to logout
func LogoutPage(httpresponsewriter http.ResponseWriter, httprequest *http.Request) {

	cookie, _ := httprequest.Cookie("DanBTCjwt")
	if cookie != nil {
		if cookie.Value != "Anonymous" {
			c := &http.Cookie{
				Name:     "DanBTCjwt",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
				HttpOnly: true,
			}
			http.SetCookie(httpresponsewriter, c)
		}
	}

	http.Redirect(httpresponsewriter, httprequest, "/", 303)
}

// LoginPage is for login
func LoginPage(httpresponsewriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client) {

	type ControllerInfo struct {
		Name    string
		Message string
	}
	type DisplayTemplate struct {
		Info ControllerInfo
	}

	items := DisplayTemplate{}
	items.Info.Name = "Login Page"

	if httprequest.Method != "POST" {

		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = ""
		t.Execute(httpresponsewriter, items)

		// http.ServeFile(httpresponsewriter, httprequest, "templates/security/login.html")
		return
	}

	usernamemix := httprequest.FormValue("userid")
	password := httprequest.FormValue("password")

	userid := strings.ToUpper(usernamemix)

	if userid == "" {
		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "Enter email address and password."
		t.Execute(httpresponsewriter, items)
		return
	}

	if password == "" {
		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "Enter email address and password."
		t.Execute(httpresponsewriter, items)
		return
	}

	cookiekeyJWT := "DanBTCjwt"
	cookiekeyUSERID := "DanBTCuserid"

	cookieJWT, _ := httprequest.Cookie(cookiekeyJWT)
	cookieUSERID, _ := httprequest.Cookie(cookiekeyUSERID)

	if cookieJWT != nil {
		cokJWT := &http.Cookie{
			Name:     cookiekeyJWT,
			Value:    "X",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(httpresponsewriter, cokJWT)
	}

	if cookieUSERID != nil {
		cokUSERID := &http.Cookie{
			Name:     cookiekeyUSERID,
			Value:    "X",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(httpresponsewriter, cokUSERID)
	}

	// Check if the user is valid and issue reference token
	//
	var resultado = LoginUserV2(redisclient, userid, password)

	if resultado.JWT == "Error" {
		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "Login error. Try again."
		t.Execute(httpresponsewriter, items)
		return
	}

	if resultado.ApplicationID != "Restaurante" {
		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "User is invalid."
		t.Execute(httpresponsewriter, items)
		return
	}

	// Store Token in Cache
	var jwttoken = resultado.JWT
	year, month, day := time.Now().Date()
	var expiry = strconv.Itoa(int(year)) + strconv.Itoa(int(month)) + strconv.Itoa(int(day))

	rediskey := "DanBTCjwt" + userid

	var credentials helper.Credentials
	credentials.UserID = userid
	credentials.KeyJWT = rediskey
	credentials.JWT = jwttoken
	credentials.Expiry = expiry
	credentials.ClaimSet = resultado.ClaimSet
	credentials.ApplicationID = resultado.ApplicationID
	credentials.UserName = resultado.Name
	credentials.IsAdmin = resultado.IsAdmin
	credentials.CentroID = resultado.CentroID

	jsonval, _ := json.Marshal(credentials)
	jsonstring := string(jsonval)

	_ = redisclient.Set(rediskey, jsonstring, 0).Err()

	// store in cookie
	// 2 hours ==> 4 hours
	expiration := time.Now().Add(4 * time.Hour)

	cokJWT := &http.Cookie{
		Name:     cookiekeyJWT,
		Value:    jwttoken,
		Path:     "/",
		Expires:  expiration,
		MaxAge:   0,
		HttpOnly: true,
	}

	http.SetCookie(httpresponsewriter, cokJWT)

	cokUSERID := &http.Cookie{
		Name:     cookiekeyUSERID,
		Value:    userid,
		Path:     "/",
		Expires:  expiration,
		MaxAge:   0,
		HttpOnly: true,
	}

	http.SetCookie(httpresponsewriter, cokUSERID)

	http.Redirect(httpresponsewriter, httprequest, "/", 303)

	return
}

// AnonymousLogin is for login
func AnonymousLogin(httpresponsewriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client, useridin string, username string) {
	log.Println("AnonymousLogin Called " + useridin)

	userid := strings.ToUpper(useridin)
	log.Println("AnonymousLogin - User ID: " + userid)

	cookiekeyJWT := "DanBTCjwt"
	cookiekeyUSERID := "DanBTCuserid"
	cookiealreadystored := "No"

	cookieJWT, _ := httprequest.Cookie(cookiekeyJWT)
	cookieUSERID, _ := httprequest.Cookie(cookiekeyUSERID)

	if cookieJWT != nil {
		//-------------------------------------------------------
		// Neste caso apenas retorne o cookie value, nao apague
		//-------------------------------------------------------

		// ??? Este e' o cookie que armazena a JWT
	}

	if cookieUSERID != nil {
		//-------------------------------------------------------
		// Neste caso apenas retorne o cookie value, nao apague
		//-------------------------------------------------------

		userid = cookieUSERID.Value
		cookiealreadystored = "Yes"
	}

	resultado := security.Credentials{}
	resultado.JWT = "Anonymous"
	resultado.ApplicationID = "Restaurante"

	// Store Token in Cache
	var jwttoken = resultado.JWT
	year, month, day := time.Now().Date()
	var expiry = strconv.Itoa(int(year)) + strconv.Itoa(int(month)) + strconv.Itoa(int(day))

	rediskey := "DanBTCjwt" + userid

	var credentials helper.Credentials
	credentials.UserID = userid
	credentials.KeyJWT = rediskey
	credentials.JWT = jwttoken
	credentials.Expiry = expiry
	credentials.ClaimSet = resultado.ClaimSet
	credentials.ApplicationID = resultado.ApplicationID
	credentials.UserName = username
	credentials.IsAdmin = resultado.IsAdmin
	credentials.CentroID = resultado.CentroID

	jsonval, _ := json.Marshal(credentials)
	jsonstring := string(jsonval)

	// ---------------------------------------
	//         Store in cache
	// ---------------------------------------
	_ = redisclient.Set(rediskey, jsonstring, 0).Err()

	if cookiealreadystored == "No" {

		// store in cookie
		// 1 month
		expiration := time.Now().Add(720 * time.Hour)
		// expiration := time.Now().Add(1 * time.Hour)

		cokJWT := &http.Cookie{
			Name:     cookiekeyJWT,
			Value:    jwttoken,
			Path:     "/",
			Expires:  expiration,
			MaxAge:   0,
			HttpOnly: true,
		}

		http.SetCookie(httpresponsewriter, cokJWT)
		log.Println("Storing Cookie: " + cookiekeyJWT)

		cokUSERID := &http.Cookie{
			Name:     cookiekeyUSERID,
			Value:    userid,
			Path:     "/",
			Expires:  expiration,
			MaxAge:   0,
			HttpOnly: true,
		}

		http.SetCookie(httpresponsewriter, cokUSERID)
		log.Println("Storing Cookie: " + cookiekeyUSERID)
	} else {

		log.Println("Reusing Cookie ! ")

	}

	// http.Redirect(httpresponsewriter, httprequest, "/", 303)

	return
}

// ControllerInfo is
type ControllerInfo struct {
	Name          string
	Message       string
	UserID        string
	UserName      string
	ApplicationID string //
	IsAdmin       string //
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
	Pratos     []dishes.Dish
}

// Instructions is for login
func Instructions(httpresponsewriter http.ResponseWriter, httprequest *http.Request, redisclient *redis.Client) {

	// create new template
	t, error := template.ParseFiles("html/homepage.html", "templates/main/instructions.html")

	if error != nil {
		panic(error)
	}

	// Assemble the display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Instructions"

	t.Execute(httpresponsewriter, items)

}
