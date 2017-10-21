package main

import (
	"fmt"
	"github.com/icrowley/fake"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rmulley/go-fast-sql"
	"strings"
	"log"
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
	var (
		numUsers = 1000000
		flush uint = 13000
		err error
		db *fastsql.DB
	)

	fmt.Printf("starting..\n")

	if db, err = fastsql.Open("mysql", "root:@tcp(localhost:3306)/test", flush); err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	for i := 0; i < numUsers; i++ {
		auser := generateUser(i)
		insertUsers(db, err, auser)
	}
}

func insertUsers(db *fastsql.DB, err error, auser user) {
	if err = db.BatchInsert("INSERT INTO users (firstname, lastname, email, state, postcode) VALUES (?, ?, ?, ?, ?)",
		auser.firstname, auser.lastname, auser.email, auser.state, auser.postcode); err != nil {
			log.Fatalln(err)
	}

	//fmt.Printf("inserted: %s %s %s %s %s\n", auser.firstname, auser.lastname, auser.state, auser.postcode, auser.email)
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
