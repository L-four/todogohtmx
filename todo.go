package main

import (
	"database/sql"
	"errors"
)

type Todo struct {
	Id        uint64
	Title     string
	Completed bool
	ListId    uint64
}

type TodoSlice []Todo

func (T *Todo) Scan(rows *sql.Rows) error {
	err := rows.Scan(&T.Id, &T.Title, &T.Completed, &T.ListId)
	return err
}

func (T *Todo) Load() error {
	if T.Id == 0 {
		return errors.New(".Id is not set. Cannot load id 0")
	}
	db := DbConnection()
	rows, err := db.Query(`SELECT id, title, completed, list_id FROM todos WHERE id = $1`, &T.Id)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	rows.Next()
	err = T.Scan(rows)

	if err != nil {
		return err
	}

	return nil
}

func (T *Todo) UpdateCompleted() error {
	if T.Id == 0 {
		return errors.New(".Id is not set. Cannot update todo with id 0")
	}
	_, err := db.Query(`UPDATE todos SET completed=$2 WHERE id=$1`, &T.Id, &T.Completed)
	if err != nil {
		panic(err)
	}
	return nil
}

func (T *Todo) Insert() error {
	db := DbConnection()
	rows, err := db.Query(`INSERT INTO todos (title, list_id, completed) VALUES ($1, $2, $3) RETURNING id`, &T.Title, &T.ListId, &T.Completed)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&T.Id)
	}
	return nil
}

func (T *Todo) Delete() error {
	if T.Id == 0 {
		return errors.New(".Id is not set. Cannot update todo with id 0")
	}
	_, err := db.Query(`DELETE FROM todos WHERE id=$1`, &T.Id)
	if err != nil {
		panic(err)
	}
	return nil
}

func (T *TodoSlice) LoadAll() uint64 {
	limit := len(*T)
	db := DbConnection()
	rows, err := db.Query(`SELECT id, title, completed, list_id FROM todos LIMIT 0, $1`, limit)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var count uint64 = 0
	for i := range *T {
		if !rows.Next() {
			break
		}
		err = (*T)[i].Scan(rows)
		count++
	}
	return count
}
