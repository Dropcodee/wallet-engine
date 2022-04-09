package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB
)

func ConnectDb() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	db_username := "root"
	db_password := "Coder__99"
	db_name := "wallet_db"
	db_driver := "mysql"

	connection, err := gorm.Open(db_driver, db_username+":"+db_password+"@/"+db_name+"?charset=utf8&parseTime=true&loc=Local")
	if err != nil {
		panic(err)
	}

	db = connection
}

func GetDB() *gorm.DB {
	return db
}
