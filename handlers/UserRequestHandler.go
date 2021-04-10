package handlers

import (
	"encoding/json"
	"errors"
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

type tokenRequestBody struct {
	RefreshToken string `json:"refreshToken"`
}

type tokenResponseBody struct {
	JWTToken     string `json:"jwtToken"`
	RefreshToken string `json:"refreshToken"`
}

const (
	ERROR_EMPTY = "Empty Value"
)

func getLogger() *log.Entry {
	handler_log := log.WithFields(log.Fields{"source": "User Request Handler"})
	return handler_log
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	getLogger().Info("Handling create user")
	user := &model.User{}
	json.NewDecoder(r.Body).Decode(user)
	err := validateRegistrationRequest(*user)
	if err != nil {
		getLogger().Error(fmt.Sprintf("Error when creating user: %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlerError{Message: err.Error(), Reason: "validation"})
		return
	}
	status := postgres.SaveUser(*user)
	if status.Code != http.StatusCreated {
		getLogger().Error(fmt.Sprintf("Error when creating user: %s", status.Message))
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
	jwt, rt, err := TokenizeUser(fetchedUser)
	if err != nil {
		getLogger().Error(fmt.Sprintf("Cannot compute jwt token %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlerError{Message: "Cannot compute jwt token", Reason: err.Error()})
		return
	}

	fetchedUser.JWTToken = jwt
	fetchedUser.RefreshToken = rt
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedUser)
	elapsed := time.Since(start).Milliseconds()
	getLogger().Info(fmt.Sprintf("Request handled in %d ms", elapsed))

}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	getLogger().Info("Handling refresh token")
	tokenRequest := tokenRequestBody{}
	json.NewDecoder(r.Body).Decode(&tokenRequest)
	token, err := jwt.Parse(tokenRequest.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		getLogger().Error(err.Error())
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(handlerError{Message: "Error parsing refresh token", Reason: "Refresh Token"})
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["userId"].(string)
		user, status := postgres.GetUserByUserId(userId)
		if status.Code != http.StatusOK {
			getLogger().Error(status.Message)
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(handlerError{Message: "Error parsing refresh token", Reason: "Refresh Token"})
			return
		}
		jwtToken, rt, err := TokenizeUser(user)
		if err != nil {
			getLogger().Error(err.Error())
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(handlerError{Message: "Error parsing refresh token", Reason: "Refresh Token"})
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tokenResponseBody{JWTToken: jwtToken, RefreshToken: rt})
		elapsed := time.Since(start).Milliseconds()
		getLogger().Info(fmt.Sprintf("Request handled in %d ms", elapsed))
	} else {
		getLogger().Error("cannot parse jwt claims")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(handlerError{Message: "Error parsing claims", Reason: "Refresh Token"})
	}
}

func TokenizeUser(user model.User) (jwtToken string, refreshToken string, err error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.UserId
	claims["username"] = user.Username
	claims["displayName"] = user.DisplayName
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	tokenStr, err := token.SignedString(jwtSecret)
	rtoken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := rtoken.Claims.(jwt.MapClaims)
	rtClaims["userId"] = user.UserId
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	rtStr, err := rtoken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}
	return tokenStr, rtStr, nil
}

func validateRegistrationRequest(user model.User) error {
	if user.DisplayName == "" {
		return errors.New("empty displayName")
	}
	if user.Username == "" {
		return errors.New("empty username")
	}
	if user.Email == "" {
		return errors.New("empty email")
	}
	if user.Password == "" {
		return errors.New("empty password")
	}
	return nil
}
