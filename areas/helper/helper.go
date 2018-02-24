package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"restauranteapi/security"
	"strings"
	"time"
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
	ReturnedValue    string
}

// Credentials is a struct
// ----------------------------------------------------
type Credentials struct {
	UserID        string // error code
	UserName      string // description
	KeyJWT        string
	JWT           string
	Expiry        string
	Roles         []string         // Y or N
	ClaimSet      []security.Claim // Y or N
	ApplicationID string           //
	IsAdmin       string           //
}

// Claim is
type Claim struct {
	Type  string
	Value string
}

func add() {
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// RestEnvVariables = restaurante environment variables
//
type RestEnvVariables struct {
	APIMongoDBLocation    string // location of the database localhost, something.com, etc
	APIMongoDBDatabase    string // database name
	APIAPIServerPort      string // collection name
	APIAPIServerIPAddress string // apiserver name
	WEBDebug              string // debug
	RecordCurrencyTick    string // debug
	RunningFromServer     string // debug
	WEBServerPort         string // collection name
	ConfigFileFound       string // collection name
	ApplicationID         string // collection name
	UserID                string // collection name
	AppFestaJuninaEnabled string
	AppBelnorthEnabled    string
	AppBitcoinEnabled     string
}

// Readfileintostruct is
func Readfileintostruct() RestEnvVariables {
	dat, err := ioutil.ReadFile("restaurante.ini")
	check(err)
	fmt.Print(string(dat))

	var restenv RestEnvVariables

	json.Unmarshal(dat, &restenv)

	if restenv.APIAPIServerIPAddress == "" {
		restenv.APIAPIServerIPAddress = "localhost"
		restenv.APIAPIServerPort = "1520"
		restenv.WEBServerPort = ":1510"
		restenv.RunningFromServer = "Ubuntu"
		restenv.WEBDebug = "Y"
		restenv.ConfigFileFound = "Not found - hardcoded values"
	}

	return restenv
}

func keyfortheday(day int) string {

	var key = "De tudo, ao meu amor serei atento antes" +
		"E com tal zelo, e sempre, e tanto" +
		"Que mesmo em face do maior encanto" +
		"Dele se encante mais meu pensamento" +
		"Quero vivê-lo em cada vão momento" +
		"E em seu louvor hei de espalhar meu canto" +
		"E rir meu riso e derramar meu pranto" +
		"Ao seu pesar ou seu contentamento" +
		"E assim quando mais tarde me procure" +
		"Quem sabe a morte, angústia de quem vive" +
		"Quem sabe a solidão, fim de quem ama" +
		"Eu possa lhe dizer do amor (que tive):" +
		"Que não seja imortal, posto que é chama" +
		"Mas que seja infinito enquanto dure"

	stringSlice := strings.Split(key, " ")
	var stringSliceFinal []string

	x := 0
	for i := 0; i < len(stringSlice); i++ {
		if len(stringSlice[0]) > 3 {
			stringSliceFinal[x] = stringSlice[i]
			x++
		}
	}

	return stringSliceFinal[day]
}

func getjwtfortoday() string {

	_, _, day := time.Now().Date()

	s := keyfortheday(day)
	h := sha1.New()
	h.Write([]byte(s))

	sha1hash := hex.EncodeToString(h.Sum(nil))

	return sha1hash
}

// Encrypt string to base64 crypto using AES
func Encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// Decrypt from base64 to decrypted string
func Decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
