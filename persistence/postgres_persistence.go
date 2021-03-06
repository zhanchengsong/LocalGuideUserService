package postgres

import (
	"errors"
	"fmt"
	"net/http"

	guuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.io/zhanchengsong/LocalGuideUserService/model"
	"github.io/zhanchengsong/LocalGuideUserService/transferObject"
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
	PG_ERROR_GENERIC   = "generic error"
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

func SaveUser(user model.User) (model.User, DatabaseStatus) {
	db, err := ConnectDB()
	if err != nil {
		getLogger().Error("Cannot connect to postgres")
		return model.User{}, DatabaseStatus{Code: http.StatusInternalServerError, Message: "Cannot connec to db", Reason: PG_ERROR_CONNECT}
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		getLogger().Error("Cannot hash password")
		return model.User{}, DatabaseStatus{Code: http.StatusInternalServerError, Message: "Cannot hash password", Reason: PG_ERROR_HASH}
	}
	// Replace clear text password with the hashed value
	user.Password = string(hashedPass)
	// Replace userId with uuid generated
	user.UserId = guuid.NewString()
	saveErr := db.Create(&user).Error
	if saveErr != nil {
		getLogger().Error(fmt.Sprintf("Cannot create user %s", saveErr.Error()))
		return model.User{}, DatabaseStatus{Code: http.StatusConflict, Message: saveErr.Error(), Reason: PG_ERROR_CREATE}
	}
	user.Password = ""
	return user, DatabaseStatus{Code: http.StatusCreated, Message: "Success", Reason: PG_SUCCESS}
}

func UpdateUser(updateUser transferObject.UserUpdateBody, userId string) (transferObject.UserUpdateBody, DatabaseStatus) {
	db, err := ConnectDB()
	if err != nil {
		getLogger().Error("Cannot connect to postgres")
		return transferObject.UserUpdateBody{}, DatabaseStatus{Code: http.StatusInternalServerError, Message: "Cannot connec to db", Reason: PG_ERROR_CONNECT}
	}
	updateResult := db.Model(&model.User{}).Where("user_id= ?", userId).Updates(updateUser)
	updateError := updateResult.Error
	if updateError != nil {
		getLogger().Error(fmt.Sprintf("Cannot update user %s", updateError.Error()))
		return transferObject.UserUpdateBody{}, DatabaseStatus{Code: http.StatusConflict, Message: updateError.Error(), Reason: PG_ERROR_GENERIC}
	}
	updateCount := updateResult.RowsAffected
	if updateCount == 0 {
		getLogger().Error("Cannot update user: No user found")
		return transferObject.UserUpdateBody{}, DatabaseStatus{Code: http.StatusNotFound, Message: "No record found", Reason: PG_ERROR_GENERIC}
	}
	return updateUser, DatabaseStatus{Code: http.StatusOK, Message: "Success", Reason: PG_SUCCESS}
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fetchedUser, DatabaseStatus{Code: http.StatusNotFound, Message: err.Error(), Reason: PG_ERROR_NO_RECORD}
		}
		return fetchedUser, DatabaseStatus{Code: http.StatusInternalServerError, Message: err.Error(), Reason: PG_ERROR_GENERIC}
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fetchedUser, DatabaseStatus{Code: http.StatusNotFound, Message: err.Error(), Reason: PG_ERROR_GENERIC}
		}
		return fetchedUser, DatabaseStatus{Code: http.StatusInternalServerError, Message: err.Error(), Reason: PG_ERROR_GENERIC}
	}
	fetchedUser.Password = ""
	return fetchedUser, DatabaseStatus{Code: http.StatusOK, Message: "success", Reason: PG_SUCCESS}
}

func SearchUsersByUsername(username string) ([]model.User, DatabaseStatus) {
	var users []model.User
	db, err := ConnectDB()
	if err != nil {
		getLogger().Error("Cannot connect to postgres")
		return users, DatabaseStatus{Code: http.StatusInternalServerError, Message: "Cannot connect to db", Reason: PG_ERROR_CONNECT}
	}

	err = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username)).Find(&users).Error
	if err != nil {

		getLogger().Error(err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return users, DatabaseStatus{Code: http.StatusNotFound, Message: err.Error(), Reason: PG_ERROR_NO_RECORD}
		}
		return users, DatabaseStatus{Code: http.StatusInternalServerError, Message: err.Error(), Reason: PG_ERROR_GENERIC}
	}
	return users, DatabaseStatus{Code: http.StatusOK, Message: "success", Reason: PG_SUCCESS}
}

func CountByUsername(username string) (int64, DatabaseStatus) {
	var testUser = model.User{}
	db, err := ConnectDB()
	if err != nil {
		getLogger().Error("Cannot connect to postgres")
	}
	dbResult := db.Where("username= ?", username).First(&testUser)
	err = dbResult.Error
	if err != nil {
		getLogger().Error(err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, DatabaseStatus{Code: http.StatusOK, Message: "success", Reason: PG_SUCCESS}
		}
		return 0, DatabaseStatus{Code: http.StatusInternalServerError, Message: err.Error(), Reason: PG_ERROR_GENERIC}
	}
	count := dbResult.RowsAffected
	return count, DatabaseStatus{Code: http.StatusOK, Message: "success", Reason: PG_SUCCESS}
}

func CountByDisplayName(displayName string) (int64, DatabaseStatus) {
	var testUser = model.User{}
	db, err := ConnectDB()
	if err != nil {
		getLogger().Error("Cannot connect to postgres")
	}
	dbResult := db.Where("display_name= ?", displayName).First(&testUser)
	err = dbResult.Error
	if err != nil {
		getLogger().Error(err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, DatabaseStatus{Code: http.StatusOK, Message: "success", Reason: PG_SUCCESS}
		}
		return 0, DatabaseStatus{Code: http.StatusInternalServerError, Message: err.Error(), Reason: PG_ERROR_GENERIC}
	}
	count := dbResult.RowsAffected
	return count, DatabaseStatus{Code: http.StatusOK, Message: "success", Reason: PG_SUCCESS}
}
