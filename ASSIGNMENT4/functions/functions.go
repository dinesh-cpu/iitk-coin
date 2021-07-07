package functions

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {

	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func Adddata(database *sql.DB, username string, password string, rollno int, batch string) {

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
