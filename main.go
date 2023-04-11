package main

import (
	"fmt"
	"log"
	"miniProject/entities"
	"miniProject/routers"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = os.Getenv("PGHOST")
	user     = os.Getenv("PGUSER")
	password = os.Getenv("PGPASSWORD")
	dbPort   = os.Getenv("PGPORT")
	dbName   = os.Getenv("PGDATABASE")
	db       *gorm.DB
	err      error
)

func init() {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbName, dbPort)

	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database :", err)
	}

	db.Debug().AutoMigrate(entities.User{}, entities.Post{}, entities.Follow{})
}

func main() {
	var PORT = os.Getenv("PORT")
	routers.StartService(db).Run(":" + PORT)
}
