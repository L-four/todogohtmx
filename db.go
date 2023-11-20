package main

import (
	"database/sql"
	"io/fs"
	"log"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func DbConnection() *sql.DB {
	clientConnStr := "db.sqlite3"
	if db == nil {
		newDb, err := sql.Open("sqlite", clientConnStr)
		if err != nil {
			log.Fatal(err)
		}
		db = newDb
	}
	return db
}

func CheckAndSetupDb(fs fs.FS) {
	db := DbConnection()
	var name string
	rows, err := db.Query("SELECT name from sqlite_master")
	if err != nil {
		log.Fatal(err)
		return
	}
	val := rows.Next()
	if val {
		err = rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
			return
		}
		// we don't need to do anything.
		return
	}

	file, err := fs.Open("db/schema.sql")
	defer file.Close()

	if err != nil {
		log.Fatal(err)
		return
	}

	data, err := ReadFile(file)

	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = db.Query(string(data))

	if err != nil {
		log.Fatal(err)
		return
	}

	file2, err := fs.Open("db/seed_data.sql")

	data, err = ReadFile(file2)

	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = db.Query(string(data))

	if err != nil {
		log.Fatal(err)
		return
	}
}
