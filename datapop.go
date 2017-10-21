package main

import (
	"fmt"
	"github.com/icrowley/fake"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"log"
	"time"
)

type user struct {
	id int
	firstname string
	lastname string
	email string
	state string
	postcode string
}

// populates a database with test data
func main() {
	numUsers := 50000

	fmt.Printf("starting..\n")

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/test")

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < numUsers; i++ {
		auser := generateUser(i)
		time.Sleep(1 * time.Millisecond)
		go insertUser(db, err, auser)
	}
}

func generateUser(i int) user {
	return user {
		id: i,
		email: strings.ToLower(fake.EmailAddress()),
		firstname : fake.FirstName(),
		lastname: fake.LastName(),
		state: fake.StateAbbrev(),
		postcode: fake.Zip(),
	}
}

func insertUser(db *sql.DB, err error, auser user) {

	insert, err := db.Prepare("INSERT INTO users (firstname, lastname, email, state, postcode) VALUES (?, ?, ?, ?, ?)")
	defer insert.Close()

	_, err = insert.Exec(auser.firstname, auser.lastname, auser.email, auser.state, auser.postcode)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("inserted: %s %s %s %s %s\n", auser.firstname, auser.lastname, auser.state, auser.postcode, auser.email)
}