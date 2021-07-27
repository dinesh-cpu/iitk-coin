package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	// "sync"
	functions "dinesh-cpu/functions"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	// helpers "github.com/dinesh-cpuiitk-coin/Helpers"
)

var jwtKey = []byte("Put jwt key here")

//Schema for user Login
type Signcred struct {
	Rollno   int    `json:"rollno"`
	Password string `json:"password"`
}

// Schema for user Signup
type Credentials struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	Rollno   int    `json:"rollno"`
	Batch    string `json:"batch"`
}

// Schema for addition of coin to user account
type Addcoin struct {
	Rollno int `json:"rollno"`
	Coin   int `json:"coin"`
}

// Schema for Redeem
type Redeem struct {
	Coin int    `json:"coin"`
	Item string `json:"item"`
}

// Schema for Transfer coin
type Transfercoin struct {
	Rollno2 int `json:"rollno1"`
	Coin    int `json:"coin"`
}

// Schema for JWT (json web tokens)
type Claims struct {
	Rollno int `json:"rollno"`
	jwt.StandardClaims
}

// Schema for Action on pending requests
type PendingAction struct {
	Id     int `json:"id"`
	Action int `json:"action"`
}

/******************SIGNIN***********************/
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
			var a bool = ((roll == creds.Rollno) && (functions.ComparePasswords(pass, []byte(creds.Password))))

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
				fmt.Fprintf(w, "Welcome, You logged In")
				return

			}
		}
		if flag {
			fmt.Fprintf(w, "invalid username or password ")
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

/************************SIGNUP****************/
func Signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Please create a POST request")
	case "POST":
		var creds Credentials
		json.NewDecoder(r.Body).Decode(&creds)
		username := creds.Name
		password := creds.Password
		hash, _ := functions.HashPassword(password)
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
			functions.Adddata(database, username, hash, creds.Rollno, batch)
			fmt.Fprintf(w, "wlecome %s,you succesfully signedup", username)

		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

/*********************GETCOIN*****************/

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

/***********************ADDCOIN***************/
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
			ctx := context.Background()
			tx, err := database.BeginTx(ctx, nil)
			if err != nil {
				return
			}
			res, err := tx.ExecContext(ctx, "UPDATE FINALDATA SET coin = coin + ? WHERE rollno=? ", cred.Coin, cred.Rollno)
			if err != nil {
				tx.Rollback()
				panic(err)
			}
			rows_affected, err := res.RowsAffected()
			if err != nil {
				panic(err)
			}

			if rows_affected != 1 {

				tx.Rollback()
				panic(err)
			}
			// err = tx.Commit()
			// if err != nil {
			// 	panic(err)
			// }

			_, err = tx.ExecContext(ctx, "UPDATE FINALDATA SET coin = coin - ? WHERE rollno = ? AND coin - ? >= 0", cred.Coin, claims.Rollno, cred.Coin)
			if err != nil {
				tx.Rollback()
				panic(err)
			}
			// rows_effected, err := Res.RowsAffected()
			// if err != nil {
			// 	panic(err)
			// }

			// if rows_effected != 2 {
			// 	tx.Rollback()
			// 	panic(err)
			// }
			err = tx.Commit()
			if err != nil {
				panic(err)
			}
			addtrans, err := database.Prepare(`INSERT INTO EVENTS(sender,reciver,amount,isreward,date) VALUES(?, ?,?,?,?)`)

			if err != nil {
				panic(err)
			}
			addtrans.Exec(claims.Rollno, cred.Rollno, cred.Coin, 1, time.Now().String())
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

/*************************TARANSFERCOIN*******************/

func TransferCOIN(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/transfer" {
	// 	http.NotFound(w, r)
	// 	return
	// }

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

		rows, err := database.Query("SELECT rollno,batch,coin FROM FINALDATA")
		if err != nil {
			panic(err)
		}
		var ROLL int
		var COIN int
		var Bth int
		var bth1 int
		var bth2 int

		var b bool
		b = false
		for rows.Next() {
			rows.Scan(&ROLL, &Bth, &COIN)

			if ROLL == cred.Rollno2 {
				b = true
				bth1 = Bth
			}
			if ROLL == claims.Rollno {

				bth2 = Bth
			}

		}
		var tax bool
		tax = (bth1 == bth2)

		if b {
			ctx := context.Background()
			tx, err := database.BeginTx(ctx, nil)
			if err != nil {
				return
			}
			rollno := claims.Rollno

			if rollno == 190304 || cred.Rollno2 == 190304 {
				fmt.Fprintf(w, "Admin cannot involve in Transfer ")

			} else {

				if tax {
					res, err := tx.ExecContext(ctx, "UPDATE FINALDATA SET coin = coin - ? WHERE rollno = ? AND coin - ? >= 0", int64(float64(cred.Coin)*1.02), rollno, int64(float64(cred.Coin)*1.02))
					if err != nil {
						tx.Rollback()
						panic(err)
					}
					rows_affected, err := res.RowsAffected()
					if err != nil {
						panic(err)
					}

					if rows_affected != 1 {

						tx.Rollback()
						panic(err)
					}
					// err = tx.Commit()
					// if err != nil {
					// 	panic(err)
					// }

					_, err = tx.ExecContext(ctx, "UPDATE FINALDATA SET coin = coin + ? WHERE rollno=? ", cred.Coin, cred.Rollno2)
					if err != nil {
						tx.Rollback()
						panic(err)
					}
					// rows_effected, err := Res.RowsAffected()
					// if err != nil {
					// 	panic(err)
					// }

					// if rows_effected != 2 {
					// 	tx.Rollback()
					// 	panic(err)
					// }
					err = tx.Commit()
					if err != nil {
						panic(err)
					}
					addtrans, err := database.Prepare(`INSERT INTO EVENTS(sender,reciver,amount,isreward,date) VALUES(?, ?,?,?,?)`)

					if err != nil {
						panic(err)
					}
					addtrans.Exec(rollno, cred.Rollno2, cred.Coin, 0, time.Now().String())
					fmt.Fprintf(w, "coin transferd ")
				} else {
					res, err := tx.ExecContext(ctx, "UPDATE FINALDATA SET coin = coin - ? WHERE rollno = ? AND coin - ? >= 0", int64(float64(cred.Coin)*1.35), rollno, int64(float64(cred.Coin)*1.35))
					if err != nil {
						tx.Rollback()
						panic(err)
					}
					rows_affected, err := res.RowsAffected()
					if err != nil {
						panic(err)
					}

					if rows_affected != 1 {

						tx.Rollback()
						panic(err)
					}
					// err = tx.Commit()
					// if err != nil {
					// 	panic(err)
					// }

					_, err = tx.ExecContext(ctx, "UPDATE FINALDATA SET coin = coin + ? WHERE rollno=? ", cred.Coin, cred.Rollno2)
					if err != nil {
						tx.Rollback()
						panic(err)
					}
					// rows_effected, err := Res.RowsAffected()
					// if err != nil {
					// 	panic(err)
					// }

					// if rows_effected != 2 {
					// 	tx.Rollback()
					// 	panic(err)
					// }
					err = tx.Commit()
					if err != nil {
						panic(err)
					}
					addtrans, err := database.Prepare(`INSERT INTO EVENTS(sender,reciver,amount,isreward,date) VALUES(?, ?,?,?,?)`)

					if err != nil {
						panic(err)
					}
					addtrans.Exec(rollno, cred.Rollno2, cred.Coin, 0, time.Now().String())
					fmt.Fprintf(w, "coin transfered")
				}

			}

			// tx.Commit()
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

/*****************LOGOUT****************/
func Logout(w http.ResponseWriter, r *http.Request) {
	expirationTime := time.Now().Add(-1 * time.Minute)
	claims := &Claims{
		Rollno: 00,
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
}

/*******************REDEEM********/
func RedeemCoin(w http.ResponseWriter, r *http.Request) {
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
		var creds Redeem

		json.NewDecoder(r.Body).Decode(&creds)

		// w.Write([]byte(fmt.Sprintf("Welcome %s!,still loggedin", claims.Rollno)))
		database, err := sql.Open("sqlite3", "./coindatabase.db")
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Trying to redeem coin")
		}

		// _, err = database.Exec("UPDATE REDEEM SET amount = ? WHERE rollno = ?", creds.Coin, claims.Rollno)
		// if err != nil {
		// 	fmt.Fprintf(w, "Server error or not sufficient balance")
		// 	panic(err)

		// }
		addredeem, err := database.Prepare(`INSERT INTO REDEEM(rollno,amount,item,status,date) VALUES(?, ?,?,?,?)`)

		if err != nil {
			panic(err)
		}
		addredeem.Exec(claims.Rollno, creds.Coin, creds.Item, "pending", time.Now().String())

		fmt.Fprintf(w, "Your request is created for admin aprovel for %d coins", creds.Coin)
		return

	case "GET":

		fmt.Fprintf(w, "Please create POST request")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

/*********************************pendingrequests*****************************/
func PENDING(w http.ResponseWriter, r *http.Request) {
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
		if claims.Rollno != 190304 {
			fmt.Fprintf(w, "This route is only for admin")
			return
		}
		database, err := sql.Open("sqlite3", "./coindatabase.db")
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Trying to get coin")
		}

		rows, err := database.Query("SELECT rollno ,amount,item,status,id FROM REDEEM")
		if err != nil {
			panic(err)
		}

		var ROLL int
		var COIN int
		var STAT string
		var ID int
		var ITEM string

		for rows.Next() {
			rows.Scan(&ROLL, &COIN, &ITEM, &STAT, &ID)
			var a bool = (STAT == "pending")
			if a {

				fmt.Fprintf(w, "Request for %s is created by %d for %d COIN with id :%d \n", ITEM, ROLL, COIN, ID)

			}
		}
		return

	case "POST":

		fmt.Fprintf(w, "Please create GET request")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

/********************************action*******************************/
func Action(w http.ResponseWriter, r *http.Request) {
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
			fmt.Println("Trying to approve or reject request")
		}

		var cred PendingAction
		json.NewDecoder(r.Body).Decode(&cred)
		if claims.Rollno != 190304 {
			fmt.Fprintf(w, "This route is only for Admin use's")
			return
		}

		// rows, err := database.Query("SELECT rollno,amount,status,id FROM REDEEM")
		// if err != nil {
		// 	panic(err)
		// }
		// var ROLL int
		// var COIN int

		// var STAT string
		// var ID int

		// for rows.Next() {
		// 	rows.Scan(&ROLL, &COIN, &STAT, &ID)

		ctx := context.Background()
		tx, err := database.BeginTx(ctx, nil)
		if err != nil {
			return
		}

		if cred.Action == 0 {
			// fmt.Fprintf(w, "In this block")
			// string acti = "rejected"
			res, err := tx.ExecContext(ctx, "DELETE FROM REDEEM WHERE id = ?", cred.Id)
			if err != nil {
				tx.Rollback()
				panic(err)
			}
			rows_affected, err := res.RowsAffected()
			if err != nil {
				panic(err)
			}

			if rows_affected != 1 {

				tx.Rollback()
				panic(err)
			}
			err = tx.Commit()
			if err != nil {
				panic(err)
			}

			fmt.Fprintf(w, "action taken for rejection")
			return
		} else if cred.Action == 1 {
			var rollno int
			var amount int
			row := database.QueryRow("SELECT rollno, amount FROM REDEEM WHERE id= ?", cred.Id)
			row.Scan(&rollno, &amount)
			_, err = database.Exec("DELETE FROM REDEEM WHERE id = ?", cred.Id)
			if err != nil {
				fmt.Fprintf(w, "Server error")
				panic(err)

			}

			_, err = database.Exec("UPDATE FINALDATA SET coin = coin - ? WHERE rollno = ? AND coin - ? >= 0", amount, rollno, amount)
			if err != nil {
				fmt.Fprintf(w, "Action for rollno : %d is not completed due to low balance or internal server error", rollno)
				panic(err)
			}

			addtrans, err := database.Prepare(`INSERT INTO EVENTS(sender,reciver,amount,isreward,date) VALUES(?, ?,?,?,?)`)

			fmt.Fprintf(w, "action is completed")
			if err != nil {
				panic(err)
			}
			addtrans.Exec(rollno, 190304, amount, 0, time.Now().String())
			return

		}

		fmt.Fprintf(w, "Invalid id ")

	case "GET":
		fmt.Fprintf(w, "Please create POST request")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
