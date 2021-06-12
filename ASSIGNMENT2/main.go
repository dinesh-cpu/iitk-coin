package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("dinesh_is_god")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Rollno   int    `json:"rollno"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {

	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func adddata(database *sql.DB, username string, password string, rollno int) {

	Adddata, err := database.Prepare(`INSERT INTO USRDATA(username, password,rollno) VALUES(?, ?,?)`)

	if err != nil {
		panic(err)
	}
	Adddata.Exec(username, password, rollno)

}

//

func Signin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "please create a POST request")
	case "POST":
		database, err := sql.Open("sqlite3", "./userdata.db")
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Trying to login")
		}

		var creds Credentials

		json.NewDecoder(r.Body).Decode(&creds)

		rows, err := database.Query("SELECT username, password,rollno  FROM USRDATA")
		if err != nil {
			panic(err)
		}
		var usr string
		var pass string
		var roll int
		var flag bool = true
		for rows.Next() {
			rows.Scan(&usr, &pass, &roll)
			var a bool = ((usr == creds.Username) && (comparePasswords(pass, []byte(creds.Password))) && roll == creds.Rollno)
			if a {

				flag = false
				expirationTime := time.Now().Add(5 * time.Minute)
				claims := &Claims{
					Username: usr,
					StandardClaims: jwt.StandardClaims{

						ExpiresAt: expirationTime.Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				tokenString, err := token.SignedString(jwtKey)
				if err != nil {

					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				cookie := http.Cookie{
					Name:    "Tok",
					Value:   tokenString,
					Expires: expirationTime,
				}
				http.SetCookie(w, &cookie)
				fmt.Fprintf(w, "Welcome :%s", creds.Username)
				return

			}
		}
		if flag {
			fmt.Fprintf(w, "Oh No! invalid username or password ")
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
func Secretpage(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("Tok")
	if err != nil {
		if err == http.ErrNoCookie {

			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Please login to acess the page")

		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome %s!,still loggedin", claims.Username)))
}
func Signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Please create a POST request")
	case "POST":
		var creds Credentials
		json.NewDecoder(r.Body).Decode(&creds)

		username := creds.Username
		password := creds.Password
		hash, _ := HashPassword(password)

		database, err := sql.Open("sqlite3", "./userdata.db")
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Connection succesfull")
		}
		rows, err := database.Query("SELECT username, password,rollno  FROM USRDATA")
		if err != nil {
			panic(err)
		}
		var usr string
		var pass string
		var roll int
		var flag bool = true
		var b bool = false
		var c bool = false
		for rows.Next() {
			rows.Scan(&usr, &pass, &roll)

			if usr == creds.Username {
				b = true
			}

			if roll == creds.Rollno {
				c = true
			}
			if b || c {
				flag = false

				fmt.Fprintf(w, "User alredy exists")
				return

			}
		}
		if flag {
			adddata(database, username, hash, creds.Rollno)
			fmt.Fprintf(w, "wlecome %s,you succesfully signedup", username)

		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

func main() {
	database, err := sql.Open("sqlite3", "./userdata.db")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connection succesfull")
	}

	Table, err := database.Prepare(`CREATE TABLE IF NOT EXISTS USRDATA("username" TEXT  NOT NULL,"password" TEXT NOT NULL,"rollno" INTEGER UNSIGNED NOT NULL);`)
	if err != nil {
		panic(err)
	}
	Table.Exec()

	// "Signin" and "Welcome" are the handlers that we will implement
	http.HandleFunc("/login", Signin)
	http.HandleFunc("/secretpage", Secretpage)
	http.HandleFunc("/signup", Signup)
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8080", nil))

}
