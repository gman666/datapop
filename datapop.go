package main

import (
	"fmt"
	"github.com/icrowley/fake"
	"database/sql"
	_ "github.com/lib/pq"
	"strings"
	"log"
	"github.com/lib/pq"
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
		err error
		db *sql.DB
		conString = "user=postgres password=x dbname=test sslmode=disable"
	)

	fmt.Printf("starting..to write %d users\n", numUsers)

	if db, err = sql.Open("postgres", conString); err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS user_test (id serial PRIMARY KEY, " +
		"firstname VARCHAR(100) DEFAULT NULL, " +
		"lastname VARCHAR(100) DEFAULT NULL, " +
		"email varchar(100) DEFAULT NULL, " +
		"state varchar(20) DEFAULT NULL, " +
		"postcode varchar(10) DEFAULT NULL)"); err != nil {
			log.Fatalln(err)
	}

	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn("user_test", "firstname", "lastname", "email", "state", "postcode"))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < numUsers; i++ {
		auser := generateUser(i)
		insertUsers(stmt, err, auser)
	}

	// flush buffer
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("inserted %d users\n", numUsers)
}

func insertUsers(stmt *sql.Stmt, err error, auser user) {
	if _, err = stmt.Exec(auser.firstname, auser.lastname, auser.email, auser.state, auser.postcode); err != nil {
		log.Fatalln(err)
	}

	//fmt.Printf("inserted: %s %s %s %s %s\n", auser.firstname, auser.lastname, auser.state, auser.postcode, auser.email)
}

func generateUser(i int) user {
	return user {
		id: i,
		firstname : fake.FirstName(),
		lastname: fake.LastName(),
		email: strings.ToLower(fake.EmailAddress()),
		state: fake.StateAbbrev(),
		postcode: fake.Zip(),
	}
}
