package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"github.io/zhanchengsong/LocalGuideUserService/model"
	postgres "github.io/zhanchengsong/LocalGuideUserService/persistence"
)

type handlerError struct {
	Message string
	Reason  string
}

func getLogger() *log.Entry {
	handler_log := log.WithFields(log.Fields{"source": "User Request Handler"})
	return handler_log
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	getLogger().Info("Handling create user")
	user := &model.User{}
	json.NewDecoder(r.Body).Decode(user)
	status := postgres.SaveUser(*user)
	if status.Code != http.StatusCreated {
		getLogger().Error(fmt.Sprintf("Error when creating user %s", status.Message))
		w.WriteHeader(status.Code)
		json.NewEncoder(w).Encode(status)
		return
	}
	user.Password = ""
	w.WriteHeader(status.Code)
	json.NewEncoder(w).Encode(user)
	elapsed := time.Since(start).Milliseconds()
	getLogger().Info(fmt.Sprintf("Request handled in %d ms", elapsed))
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	getLogger().Info("Handling login user")
	user := &model.User{}
	json.NewDecoder(r.Body).Decode(user)
	fetchedUser, status := postgres.GetUserByUsernameAndPassword(user.Username, user.Password)
	if status.Code != http.StatusOK {
		getLogger().Error(fmt.Sprintf("Error when logining user %s", status.Message))
		w.WriteHeader(status.Code)
		json.NewEncoder(w).Encode(status)
		return
	}
	// Calculate token
	jwt, err := TokenizeUser(fetchedUser)
	if err != nil {
		getLogger().Error(fmt.Sprintf("Cannot compute jwt token %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlerError{Message: "Cannot compute jwt token", Reason: err.Error()})
		return
	}
	fetchedUser.JWTToken = jwt
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedUser)
	elapsed := time.Since(start).Milliseconds()
	getLogger().Info(fmt.Sprintf("Request handled in %d ms", elapsed))

}

func TokenizeUser(user model.User) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.UserId
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
