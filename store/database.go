package store

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
)

func DbTest() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	db.Exec("CREATE TABLE sample (key VARCHAR, value VARCHAR)")
	db.Exec("INSERT INTO sample VALUES ('key1', 'value1')")
	db.Exec("INSERT INTO sample VALUES ('key2', 'value2')")

	rows, err := db.Query("SELECT key, value FROM sample")
	for rows.Next() {
		var key, value string
		err = rows.Scan(&key, &value)

		fmt.Println(key, value)
	}
}
