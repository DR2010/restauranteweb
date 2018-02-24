package security

import (
	"encoding/json"
	"net/http"
	helper "restauranteweb/areas/helper"
	"strconv"
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

	username := req.FormValue("username")
	preferredname := req.FormValue("preferredname")
	password := req.FormValue("password")
	passwordvalidate := req.FormValue("passwordvalidate")
	applicationid := req.FormValue("applicationid")

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

	userid := httprequest.FormValue("userid")
	password := httprequest.FormValue("password")

	if userid == "" {
		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "Please enter details."
		t.Execute(httpresponsewriter, items)
		return
	}

	if password == "" {
		t, _ := template.ParseFiles("templates/security/login.html", "templates/security/loginmessagetemplate.html")
		items.Info.Message = "Please enter details."
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

	jsonval, _ := json.Marshal(credentials)
	jsonstring := string(jsonval)

	_ = redisclient.Set(rediskey, jsonstring, 0).Err()

	// store in cookie
	expiration := time.Now().Add(1 * 2 * time.Hour)

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
