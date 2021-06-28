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

type Signcred struct {
	Rollno   int    `json:"rollno"`
	Password string `json:"password"`
}

type Credentials struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	Rollno   int    `json:"rollno"`
	Batch    string `json:"batch"`
}

// Tag      string `json:"tag"`
// 	Coin     int    `json:"coin"`
// 	Events   int    `json:"events"`

type Credaddcoin struct {
	Rollno int `json:"rollno"`
}

type Addcoin struct {
	Rollno int `json:"rollno"`
	Coin   int `json:"coin"`
}

type Transfercoin struct {
	Rollno2 int `json:"rollno1"`
	Coin    int `json:"coin"`
}

type Claims struct {
	Rollno int `json:"rollno"`
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

func adddata(database *sql.DB, username string, password string, rollno int, batch string) {

	Adddata, err := database.Prepare(`INSERT INTO FINALDATA(name, password,rollno,batch,tag,events,coin) VALUES(?, ?,?,?,?,?,?)`)

	if err != nil {
		panic(err)
	}
	if rollno == 190304 {
		Adddata.Exec(username, password, rollno, batch, "admin", 0, 0)
	} else {
		Adddata.Exec(username, password, rollno, batch, "user", 0, 0)
	}

}

//

func Signin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "please create a POST request")
	case "POST":
		database, err := sql.Open("sqlite3", "./coindatabase.db")
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Trying to login")
		}

		var creds Signcred

		json.NewDecoder(r.Body).Decode(&creds)

		rows, err := database.Query("SELECT password,rollno  FROM FINALDATA")
		if err != nil {
			panic(err)
		}

		var pass string
		var roll int
		var flag bool = false
		for rows.Next() {
			rows.Scan(&pass, &roll)
			var a bool = ((roll == creds.Rollno) && (comparePasswords(pass, []byte(creds.Password))))

			if a {

				flag = false
				expirationTime := time.Now().Add(20 * time.Minute)
				claims := &Claims{
					Rollno: roll,
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
				fmt.Fprintf(w, "Welcome to the page")
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

	w.Write([]byte(fmt.Sprintf("Welcome %d!,still loggedin", claims.Rollno)))
}
func Signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Please create a POST request")
	case "POST":
		var creds Credentials
		json.NewDecoder(r.Body).Decode(&creds)

		username := creds.Name
		password := creds.Password
		hash, _ := HashPassword(password)
		batch := creds.Batch

		database, err := sql.Open("sqlite3", "./coindatabase.db")
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Connection succesfull")
		}
		rows, err := database.Query("SELECT name,rollno  FROM FINALDATA")
		if err != nil {
			panic(err)
		}
		var usr string

		var roll int
		var flag bool = true
		var b bool = false
		var c bool = false
		for rows.Next() {
			rows.Scan(&usr, &roll)

			if usr == creds.Name {
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
			adddata(database, username, hash, creds.Rollno, batch)
			fmt.Fprintf(w, "wlecome %s,you succesfully signedup", username)

		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

func GETCOIN(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
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

		// w.Write([]byte(fmt.Sprintf("Welcome %s!,still loggedin", claims.Rollno)))
		database, err := sql.Open("sqlite3", "./coindatabase.db")
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Trying to get coin")
		}

		rows, err := database.Query("SELECT rollno ,coin FROM FINALDATA")
		if err != nil {
			panic(err)
		}

		var ROLL int
		var COIN int

		for rows.Next() {
			rows.Scan(&ROLL, &COIN)
			var a bool = (ROLL == claims.Rollno)
			if a {

				fmt.Fprintf(w, "Coins in your wallet: %d", COIN)
				return

			}
		}

	case "POST":

		fmt.Fprintf(w, "Please create GET request")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func ADDCOIN(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
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

		// w.Write([]byte(fmt.Sprintf("Welcome %s!,still loggedin", claims.Rollno)))
		database, err := sql.Open("sqlite3", "./coindatabase.db")
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Trying to add coin")
		}

		var cred Addcoin
		json.NewDecoder(r.Body).Decode(&cred)
		if cred.Rollno == 190304 {
			fmt.Fprintf(w, "You can not add money to your account")
			return
		}

		rows, err := database.Query("SELECT rollno ,coin FROM FINALDATA")
		if err != nil {
			panic(err)
		}
		var ROLL int
		var COIN int
		var a bool
		var b bool = claims.Rollno == 190304

		a = false
		for rows.Next() {
			rows.Scan(&ROLL, &COIN)

			if ROLL == cred.Rollno {
				a = true

			}

		}

		if a && b {

			_, err = database.Exec("UPDATE FINALDATA SET coin = coin + ? WHERE rollno=? ", cred.Coin, cred.Rollno)
			if err != nil {

				panic(err)
			}
			fmt.Fprintf(w, "%d coin added to to the %d", cred.Coin, cred.Rollno)

			return

		} else {
			fmt.Fprintf(w, "Invalid user credentials or this opretion is only valid for admin")
		}

	case "GET":

		fmt.Fprintf(w, "Please create POST request")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func TransferCOIN(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
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
		database, err := sql.Open("sqlite3", "./coindatabase.db")
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Trying to transfer the coin")
		}

		var cred Transfercoin
		json.NewDecoder(r.Body).Decode(&cred)

		rows, err := database.Query("SELECT rollno ,coin FROM FINALDATA")
		if err != nil {
			panic(err)
		}
		var ROLL int
		var COIN int

		var b bool
		b = false
		for rows.Next() {
			rows.Scan(&ROLL, &COIN)

			if ROLL == cred.Rollno2 {
				b = true

			}

		}

		if b {
			tx, err := database.Begin()
			if err != nil {
				return
			}
			rollno := claims.Rollno
			if rollno == 190304 {

				_, err = database.Exec("UPDATE FINALDATA SET coin = coin - ? WHERE rollno = ? AND coin - ? >= 0", cred.Coin, rollno, cred.Coin)
				if err != nil {
					tx.Rollback()
					panic(err)
				}
				_, err = database.Exec("UPDATE FINALDATA SET coin = coin + ? WHERE rollno=? ", cred.Coin, cred.Rollno2)
				if err != nil {
					tx.Rollback()
					panic(err)
				}
				addtrans, err := database.Prepare(`INSERT INTO EVENTS(sender,reciver,amount,isreward,date) VALUES(?, ?,?,?,?)`)

				if err != nil {
					panic(err)
				}
				addtrans.Exec(rollno, cred.Rollno2, cred.Coin, 1, time.Now().String())
				// if rollno == 190304 {
				// 	addtrans.Exec(rollno, cred.Rollno2, cred.Coin, 1, time.Now().String())
				// } else {
				// 	addtrans.Exec(rollno, cred.Rollno2, cred.Coin, 0, time.Now().String())
				// }
			} else {

				_, err = database.Exec("UPDATE FINALDATA SET coin = coin - ? WHERE rollno = ? AND coin - ? >= 0", cred.Coin*(102/100), rollno, cred.Coin)
				if err != nil {
					tx.Rollback()
					panic(err)
				}
				_, err = database.Exec("UPDATE FINALDATA SET coin = coin + ? WHERE rollno=? ", cred.Coin, cred.Rollno2)
				if err != nil {
					tx.Rollback()
					panic(err)
				}
				addtrans, err := database.Prepare(`INSERT INTO EVENTS(sender,reciver,amount,isreward,date) VALUES(?, ?,?,?,?)`)

				if err != nil {
					panic(err)
				}
				addtrans.Exec(rollno, cred.Rollno2, cred.Coin, 0, time.Now().String())
			}

			tx.Commit()
			return
		} else {
			fmt.Fprintf(w, "Invalid  user credentials")
		}

	case "GET":

		fmt.Fprintf(w, "Please create POST request")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

func main() {
	database, err := sql.Open("sqlite3", "./coindatabase.db")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connection succesfull")
	}

	Table, err := database.Prepare(`CREATE TABLE IF NOT EXISTS FINALDATA("name" TEXT  NOT NULL,"password" TEXT NOT NULL,"rollno" INTEGER UNSIGNED NOT NULL,"batch" INTEGER UNSIGNED NOT NULL,"tag" TEXT UNSIGNED NOT NULL,"events" INTEGER UNSIGNED NOT NULL,"coin" INTEGER UNSIGNED NOT NULL);`)
	if err != nil {
		panic(err)
	}
	Table.Exec()
	table, err := database.Prepare(`CREATE TABLE IF NOT EXISTS EVENTS("sender" INTEGER UNSIGNED NOT NULL,"reciver" INTEGER UNSIGNED NOT NULL,"amount" INTEGER UNSIGNED NOT NULL,"isreward" INTEGER UNSIGNED NOT NULL DEFAULT 0,"date" "date" TEXT NOT NULL);`)
	if err != nil {
		panic(err)
	}

	table.Exec()

	// Table, err = database.Prepare(`ALTER TABLE COIN ADD COLUMN coin INTEGER NOT NULL DEFAULT 0;`)
	// if err != nil {
	// 	panic(err)
	// }

	// Adddata, err := database.Prepare(`UPDATE FINALDATA SET coin = ? WHERE rollno = ?;`)
	// if err != nil {
	// 	panic(err)
	// }
	// Adddata.Exec(120, 190304)

	http.HandleFunc("/login", Signin)
	http.HandleFunc("/secretpage", Secretpage)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/getcoin", GETCOIN)
	http.HandleFunc("/addcoin", ADDCOIN)
	http.HandleFunc("/transfer", TransferCOIN)
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8080", nil))

}
