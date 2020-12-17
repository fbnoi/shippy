package main

import (
	"fmt"
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CreateSession(host, port, database, charset,username, password string) (*gorm.DB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", username, password, host, port, database, charset)
	if db, err := gorm.Open(mysql.Open(dns), &gorm.Config{}); err != nil {
		log.Fatalln(dns)
		log.Fatalln(err)
		return nil, err
	} else {
		return db, nil
	}
}