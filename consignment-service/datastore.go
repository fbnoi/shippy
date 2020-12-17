package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"
)

func CreateSession(host, port, database, charset,username, password string) (*sql.DB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", username, password, host, port, database, charset)
	log.Print(dns)
	if db,err := sql.Open("mysql",dns); err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}
		return db, nil
	}
}
