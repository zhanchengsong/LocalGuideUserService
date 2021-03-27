package postgres

import (
	"fmt"
	"os"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"github.io/zhanchengsong/LocalGuideUserService/model"
)

var (
	HOST     = os.Getenv("POSTGRES_HOST")
	USERNAME = os.Getenv("POSTGRES_USERNAME")
	PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	DATABASE = os.Getenv("POSTGRES_DATABASE")
)

//Used to execute client creation procedure only once.
var dbOnce sync.Once

var clientInstance *gorm.DB
var clientError error
func ConnectDB() (*gorm.DB, error) {
	dbOnce.Do(func() {
		log.WithFields(log.Fields{
			"source": "Postgres Gorm",
		}).Info(fmt.Sprintf("Connecting to postgres at %s", HOST))
		dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", HOST, USERNAME, DATABASE, PASSWORD)
		db, err := gorm.Open("postgres", dbURI)
		if err != nil {
			log.Error("DB connection failed")
			clientError = err
		}
		clientInstance = db
		defer db.Close()
		db.AutoMigrate(
			&model.User{})
		log.Info("Successfully connected to db")
	})

	return clientInstance, clientError
}
