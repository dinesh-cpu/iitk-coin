package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func adddata(database *sql.DB, rollno int, name string) {

	Adddata, err := database.Prepare(`INSERT INTO User(rollno, name) VALUES(?, ?)`)

	if err != nil {
		panic(err)
	}
	Adddata.Exec(rollno, name)

}
func main() {
	database, err := sql.Open("sqlite3", "./studentinfo.db")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connection succesfull")
	}

	Table, err := database.Prepare(`CREATE TABLE IF NOT EXISTS User("rollno" INTEGER NOT NULL,"name" TEXT NOT NULL);`)
	if err != nil {
		panic(err)
	}
	Table.Exec()
	adddata(database, 190355, "maink")

}
