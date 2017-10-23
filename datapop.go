package main

import (
	"fmt"
	"github.com/icrowley/fake"
	"database/sql"
	_ "github.com/lib/pq"
	"strings"
	"log"
	"github.com/lib/pq"
	"time"
	"flag"
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
		txn *sql.Tx
		err error
		db *sql.DB
		stmt *sql.Stmt;
		conString = "user=postgres password=x dbname=test sslmode=disable"
	)

	var numRecords int
	flag.IntVar(&numRecords, "num_records", 1000, "create this many records in the database")
	flag.Parse()

	fmt.Printf("starting..to write %d users\n", numRecords)

	if db, err = sql.Open("postgres", conString); err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	createTable(db)

	if txn, err = db.Begin(); err != nil {
		log.Fatalln(err)
	}

	if stmt, err = txn.Prepare(pq.CopyIn("user_test", "firstname", "lastname", "email", "state", "postcode")); err != nil {
		log.Fatalln(err)
	}

	startTime := time.Now().UnixNano()

	for i := 0; i < numRecords; i++ {
		auser := generateUser()
		insertUsers(stmt, err, auser)
	}

	// flush buffer
	if _, err = stmt.Exec(); err != nil {
		log.Fatalln(err)
	}

	if err = stmt.Close(); err != nil {
		log.Fatalln(err)
	}

	if err = txn.Commit(); err != nil {
		log.Fatalln(err)
	}

	endTime := time.Now().UnixNano()
	fmt.Printf("inserted %d users in %d ms\n", numRecords, (endTime - startTime) / 1000000)
}

func createTable(db * sql.DB) {
	var err error
	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS user_test (id serial PRIMARY KEY, " +
		"firstname VARCHAR(100) DEFAULT NULL, " +
		"lastname VARCHAR(100) DEFAULT NULL, " +
		"email varchar(100) DEFAULT NULL, " +
		"state varchar(20) DEFAULT NULL, " +
		"postcode varchar(10) DEFAULT NULL)"); err != nil {
		log.Fatalln(err)
	}
}

func insertUsers(stmt *sql.Stmt, err error, auser user) {
	if _, err = stmt.Exec(auser.firstname, auser.lastname, auser.email, auser.state, auser.postcode); err != nil {
		log.Fatalln(err)
	}
	//fmt.Printf("inserted: %s %s %s %s %s\n", auser.firstname, auser.lastname, auser.state, auser.postcode, auser.email)
}

func generateUser() user {
	return user {
		firstname : fake.FirstName(),
		lastname: fake.LastName(),
		email: strings.ToLower(fake.EmailAddress()),
		state: fake.StateAbbrev(),
		postcode: fake.Zip(),
	}
}
