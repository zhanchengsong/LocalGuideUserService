package postgres

import (
	"fmt"
	"net/http"

	guuid "github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.io/zhanchengsong/LocalGuideUserService/model"
	"golang.org/x/crypto/bcrypt"
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
	PG_ERROR_CONNECT   = "cannot connect to db"
	PG_ERROR_HASH      = "cannot hash"
	PG_ERROR_CREATE    = "cannot create user"
	PG_SUCCESS         = "success"
	PG_PASSWORD        = "wrong password"
)

type DatabaseStatus struct {
	Code    int
	Message string
	Reason  string
}

func getLogger() *log.Entry {
	pg_log := log.WithFields(log.Fields{"source": "Postgress Persistence"})
	return pg_log
}

func SaveUser(user model.User) DatabaseStatus {
	db, err := ConnectDB()
	if err != nil {
		getLogger().Error("Cannot connect to postgres")
		return DatabaseStatus{Code: http.StatusInternalServerError, Message: "Cannot connec to db", Reason: PG_ERROR_CONNECT}
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		getLogger().Error("Cannot hash password")
		return DatabaseStatus{Code: http.StatusInternalServerError, Message: "Cannot hash password", Reason: PG_ERROR_HASH}
	}
	// Replace clear text password with the hashed value
	user.Password = string(hashedPass)
	// Replace userId with uuid generated
	user.UserId = guuid.NewString()
	saveErr := db.Create(&user).Error
	if saveErr != nil {
		getLogger().Error(fmt.Sprintf("Cannot create user %s", saveErr.Error()))
		return DatabaseStatus{Code: http.StatusConflict, Message: saveErr.Error(), Reason: PG_ERROR_CREATE}
	}
	return DatabaseStatus{Code: http.StatusCreated, Message: "Success", Reason: PG_SUCCESS}
}

func GetUserByUsernameAndPassword(username string, password string) (user model.User, status DatabaseStatus) {
	fetchedUser := model.User{}
	db, err := ConnectDB()
	if err != nil {
		getLogger().Error("Cannot connect to postgres")
		return fetchedUser, DatabaseStatus{Code: http.StatusInternalServerError, Message: "Cannot connect to db", Reason: PG_ERROR_CONNECT}
	}

	err = db.Where("username= ?", username).First(&fetchedUser).Error
	if err != nil {
		getLogger().Error(err.Error())
		return fetchedUser, DatabaseStatus{Code: http.StatusNotFound, Message: err.Error(), Reason: PG_ERROR_NO_RECORD}
	}

	// Verify password
	errf := bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte(password))
	fetchedUser.Password = ""
	if errf != nil {
		getLogger().Error(errf.Error())
		return fetchedUser, DatabaseStatus{Code: http.StatusForbidden, Message: errf.Error(), Reason: PG_PASSWORD}
	}
	return fetchedUser, DatabaseStatus{Code: http.StatusOK, Message: "success", Reason: PG_SUCCESS}
}

func GetUserByUserId(userId string) (model.User, DatabaseStatus) {
	fetchedUser := model.User{}
	db, err := ConnectDB()
	if err != nil {
		getLogger().Error("Cannot connect to postgres")
		return fetchedUser, DatabaseStatus{Code: http.StatusInternalServerError, Message: "Cannot connect to db", Reason: PG_ERROR_CONNECT}
	}

	err = db.Where("user_id= ?", userId).First(&fetchedUser).Error
	if err != nil {
		getLogger().Error(err.Error())
		return fetchedUser, DatabaseStatus{Code: http.StatusNotFound, Message: err.Error(), Reason: PG_ERROR_NO_RECORD}
	}
	fetchedUser.Password = ""
	return fetchedUser, DatabaseStatus{Code: http.StatusOK, Message: "success", Reason: PG_SUCCESS}
}
