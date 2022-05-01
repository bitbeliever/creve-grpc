package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var db *gorm.DB

func init() {
	var err error
	//db, err = gorm.Open(mysql.Open("hs:hs@tcp(172.22.183.60:3306)/db?charset=utf8mb4&parseTime=True&loc=Local"))
	db, err = gorm.Open(mysql.Open("hs:hs@tcp(mysql-service:3306)/db?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}

	d, err := db.DB()
	if err != nil {
		panic(err)
	}
	if err := d.Ping(); err != nil {
		panic(err)
	}

	go keepalive()
}

func keepalive() {
	ticker := time.NewTicker(time.Minute * 30)
	for {
		select {
		case <-ticker.C:
			d, err := db.DB()
			if err != nil {
				log.Println(err)
			}
			if err := d.Ping(); err != nil {
				log.Println(err)
			}
		}
	}
}
