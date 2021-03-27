package postgres

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/zhanchengsong/LocalGuideUserService/model"
)

var (
	HOST     = os.Getenv("POSTGRES_HOST")
	USERNAME = os.Getenv("POSTGRES_USERNAME")
	PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	DATABASE = os.Getenv("POSTGRES_DATABASE")
)

func ConnectDB(username string, password string, databaseName string, databaseHost string) *gorm.DB {
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", HOST, USERNAME, DATABASE, PASSWORD)
	db, err := gorm.Open("postgres", dbURI)
	if err != nil {
		log.Println(err)
		log.Fatal("DB connection failed")
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(
		&model.User{})
	log.Print("Succesfully connected to db")
	return db
}
