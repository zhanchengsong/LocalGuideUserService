package postgres

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.io/zhanchengsong/LocalGuideUserService/model"
)

const (
	DUPLICATE_USERNAME  = "DUPLICATE_USERNAME"
	DUPLICATE_EMAIL     = "DUPLICATE_EMAIL"
	INVALID_PASSWORD    = "INVALID_PASSWORD"
	INVALID_USERNAME    = "INVALID_USERNAME"
	USERNAME_NOT_EXISTS = "NO_SUCH_USERNAME"
	CONNECTION          = "DB_CONNECTION"
	OTHER               = "OTHER"
)

const (
	PG_ERROR_NO_RECORD = "record not found"
)

type DatabaseError struct {
	error
	Code    int
	Message string
	Reason  string
}

func getLogger() *log.Entry {
	pg_log := log.WithFields(log.Fields{"source": "Postgress Persistence"})
	return pg_log
}

func (dbError DatabaseError) Error() string {
	return dbError.Message
}

func CheckUserByUsername(username string) (bool, error) {
	getLogger().Debug(fmt.Sprintf("Checking user existence with username: %s", username))
	resultUser := model.User{}
	db, err := ConnectDB()
	if err != nil {
		getLogger().Error("Failed to get postgress connection")
		return false, DatabaseError{Code: http.StatusInternalServerError, Message: err.Error(), Reason: CONNECTION}
	}
	findResult := db.Where("username = ?", username).Find(&resultUser)
	if findResult.Error != nil {
		getLogger().Error(findResult.Error)
		if findResult.Error.Error() == PG_ERROR_NO_RECORD {
			return false, DatabaseError{Code: http.StatusNotFound, Message: "User not found", Reason: USERNAME_NOT_EXISTS}
		}
		return false, DatabaseError{Code: http.StatusInternalServerError, Message: findResult.Error.Error(), Reason: OTHER}
	}
	return false, nil

}
