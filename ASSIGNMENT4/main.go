package main

import (
	"database/sql"

	"fmt"
	"log"

	handler "dinesh-cpu/handlers"

	"net/http"
)

func createUserTable(db *sql.DB) {
	Table, err := db.Prepare(`CREATE TABLE IF NOT EXISTS FINALDATA("name" TEXT  NOT NULL,"password" TEXT NOT NULL,"rollno" INTEGER UNSIGNED NOT NULL,"batch" INTEGER UNSIGNED NOT NULL,"tag" TEXT UNSIGNED NOT NULL,"events" INTEGER UNSIGNED NOT NULL,"coin" INTEGER UNSIGNED NOT NULL);`)
	if err != nil {
		panic(err)
	}
	Table.Exec()
}
func createTransactionTable(db *sql.DB) {
	table, err := db.Prepare(`CREATE TABLE IF NOT EXISTS EVENTS("sender" INTEGER UNSIGNED NOT NULL,"reciver" INTEGER UNSIGNED NOT NULL,"amount" INTEGER UNSIGNED NOT NULL,"isreward" INTEGER UNSIGNED NOT NULL DEFAULT 0,"date" TEXT NOT NULL);`)
	if err != nil {
		panic(err)
	}

	table.Exec()
}

func main() {
	database, err := sql.Open("sqlite3", "./coindatabase.db")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connection succesfull")
	}

	createUserTable(database)        //created the table for user credentials
	createTransactionTable(database) //created the transuction table

	http.HandleFunc("/login", handler.Signin)

	http.HandleFunc("/signup", handler.Signup)
	http.HandleFunc("/getcoin", handler.GETCOIN)
	http.HandleFunc("/addcoin", handler.ADDCOIN)
	http.HandleFunc("/transfer", handler.TransferCOIN)
	http.HandleFunc("/redeemcoins", handler.RedeemCoin)
	http.HandleFunc("/logout", handler.Logout)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
