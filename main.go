package main

import (
	"database/sql"
	handler "dinesh-cpu/handlers"
	"fmt"
	"log"
	"net/http"
)

func createUserTable(db *sql.DB) {
	Table, err := db.Prepare(`CREATE TABLE IF NOT EXISTS FINALDATA("name" TEXT  NOT NULL,"password" TEXT NOT NULL,"rollno" INTEGER UNSIGNED NOT NULL,"batch" INTEGER UNSIGNED NOT NULL,"tag" TEXT UNSIGNED NOT NULL,"events" INTEGER UNSIGNED NOT NULL,"coin" INTEGER UNSIGNED NOT NULL);`)
	if err != nil {
		panic(err)
	}
	Table.Exec()
	fmt.println("FINALDATA table created (if not existed) successfully!")
}
func createTransactionTable(db *sql.DB) {
	table, err := db.Prepare(`CREATE TABLE IF NOT EXISTS EVENTS("sender" INTEGER UNSIGNED NOT NULL,"reciver" INTEGER UNSIGNED NOT NULL,"amount" INTEGER UNSIGNED NOT NULL,"isreward" INTEGER UNSIGNED NOT NULL DEFAULT 0,"date" TEXT NOT NULL);`)
	if err != nil {
		panic(err)
	}

	table.Exec()
	fmt.println("EVENTS table created (if not existed) successfully!")
}

func createRedeemTable(db *sql.DB) {
	table, err := db.Prepare(`CREATE TABLE IF NOT EXISTS REDEEM("rollno" INTEGER UNSIGNED NOT NULL,"amount" INTEGER UNSIGNED NOT NULL,"item" TEXT NOT NULL,"status" TEXT NOT NULL,"date" TEXT NOT NULL,id INTEGER PRIMARY KEY AUTOINCREMENT);`)
	if err != nil {
		panic(err)
	}

	table.Exec()
	fmt.println("REDEEM table created (if not existed) successfully!")
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
	createRedeemTable(database)      //created the redeem table
	// table, _ := database.Prepare(`DROP TABLE REDEEM`)
	// table.Exec()

	http.HandleFunc("/login", handler.Signin)            // Login Route
	http.HandleFunc("/signup", handler.Signup)           // Signup Route
	http.HandleFunc("/getcoin", handler.GETCOIN)         // Route for checking available balance
	http.HandleFunc("/addcoin", handler.ADDCOIN)         // route   for admin to  add coin
	http.HandleFunc("/transfer", handler.TransferCOIN)   // Route for transfer coin
	http.HandleFunc("/redeemcoins", handler.RedeemCoin)  // Route for user to redeem coin
	http.HandleFunc("/logout", handler.Logout)           // Route for Logout
	http.HandleFunc("/pendingrequests", handler.PENDING) // Route for admin to check all pending requests for reedeem coin
	http.HandleFunc("/action", handler.Action)           // Route for action for admin on pending requests

	log.Fatal(http.ListenAndServe(":8080", nil)) //starting server on Port  8080

}
