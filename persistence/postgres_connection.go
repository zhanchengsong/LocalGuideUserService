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
	HOST     = os.Getenv("PG_HOST")
	USERNAME = os.Getenv("PG_USERNAME")
	PASSWORD = os.Getenv("PG_PASSWORD")
	DATABASE = os.Getenv("PG_DBNAME")
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
		dbURI := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", HOST, USERNAME, DATABASE, PASSWORD)
		db, err := gorm.Open("postgres", dbURI)
		db.LogMode(true)
		if err != nil {
			log.Error("DB connection failed")
			clientError = err
			return
		}
		clientInstance = db
		db.AutoMigrate(
			&model.User{})
		log.Info("Successfully connected to db")
	})

	return clientInstance, clientError
}
