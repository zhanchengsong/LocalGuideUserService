package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.io/zhanchengsong/LocalGuideUserService/transferObject"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"github.io/zhanchengsong/LocalGuideUserService/model"
	postgres "github.io/zhanchengsong/LocalGuideUserService/persistence"
)

type handlerError struct {
	Message string
	Reason  string
}

type TokenRequestBody struct {
	RefreshToken string `json:"refreshToken"`
}

type TokenResponseBody struct {
	JWTToken     string `json:"jwtToken"`
	RefreshToken string `json:"refreshToken"`
}

type UsersResponseBody struct {
	Users []model.User `json:"users"`
}

type UserCountResponseBody struct {
	Count int64 `json:"count"`
}

func getLogger() *log.Entry {
	handlerLog := log.WithFields(log.Fields{"source": "User Request Handler"})
	return handlerLog
}

// CreateUser godoc
// @Summary Create a user
// @Description Create a user profile
// @Tags Create a user
// @Accept  json
// @Produce  json
// @Param user body transferObject.UserRegisterBody true "Create user"
// @Success 201 {object} transferObject.UserResponseBody
// @Failure 409 {object} handlerError
// @Failure 500 {object} handlerError
// @Router /user [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	getLogger().Info("Handling create user")
	user := &model.User{}
	json.NewDecoder(r.Body).Decode(user)
	err := validateRegistrationRequest(*user)
	if err != nil {
		getLogger().Error(fmt.Sprintf("Error when creating user: %s", err.Error()))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlerError{Message: err.Error(), Reason: "validation"})
		return
	}
	savedUser, status := postgres.SaveUser(*user)
	if status.Code != http.StatusCreated {
		getLogger().Error(fmt.Sprintf("Error when creating user: %s", status.Message))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status.Code)
		json.NewEncoder(w).Encode(handlerError{Message: status.Message, Reason: status.Reason})
		return
	}
	savedUser.Password = ""
	// Create the token at the same time
	jwt, rt, err := TokenizeUser(savedUser)
	if err != nil {
		getLogger().Error(fmt.Sprintf("Cannot compute jwt token %s", err.Error()))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlerError{Message: "Cannot compute jwt token", Reason: err.Error()})
		return
	}

	savedUser.JWTToken = jwt
	savedUser.RefreshToken = rt
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status.Code)
	json.NewEncoder(w).Encode(savedUser)
	elapsed := time.Since(start).Milliseconds()
	getLogger().Info(fmt.Sprintf("Request handled in %d ms", elapsed))
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update an existing user profile with provided body
// @Tags Update user
// @Accept  json
// @Produce  json
// @Param userId query string true "User id to update profile"
// @Param user body transferObject.UserUpdateBody true "Body to update"
// @Success 200 {object} transferObject.UserUpdateBody
// @Failure 409 {object} handlerError
// @Failure 500 {object} handlerError
// @Router /user [patch]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	getLogger().Info("Handling update user")
	userId := r.URL.Query()["userId"][0]
	if userId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlerError{Message: "Missing userId in path", Reason: "parameter"})
		return
	}
	updateUser := &transferObject.UserUpdateBody{}
	json.NewDecoder(r.Body).Decode(updateUser)
	updateResult, status := postgres.UpdateUser(*updateUser, userId)
	if status.Code != http.StatusOK {
		getLogger().Error(fmt.Sprintf("Error when updating user: %s", status.Message))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status.Code)
		json.NewEncoder(w).Encode(handlerError{Message: status.Message, Reason: status.Reason})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updateResult)
	elapsed := time.Since(start).Milliseconds()
	getLogger().Info(fmt.Sprintf("Request handled in %d ms", elapsed))
}

// LoginUser godoc
// @Summary Login a user and obtain jwtToken/refreshToken
// @Description Takes in username and password to assign token
// @Tags Log in a user
// @Accept  json
// @Produce  json
// @Param user body transferObject.UserLoginBody true "Login user"
// @Success 200 {object} transferObject.UserResponseBody
// @Failure 404 {object} handlerError
// @Failure 500 {object} handlerError
// @Failure 403 {object} handlerError
// @Router /login [post]
func LoginUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	getLogger().Info("Handling login user")
	login := &transferObject.UserLoginBody{}
	json.NewDecoder(r.Body).Decode(login)
	fetchedUser, status := postgres.GetUserByUsernameAndPassword(login.Username, login.Password)
	if status.Code != http.StatusOK {
		getLogger().Error(fmt.Sprintf("Error when logining user %s", status.Message))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status.Code)

		json.NewEncoder(w).Encode(handlerError{Message: status.Message, Reason: status.Reason})
		return
	}
	// Calculate token
	jwt, rt, err := TokenizeUser(fetchedUser)
	if err != nil {
		getLogger().Error(fmt.Sprintf("Cannot compute jwt token %s", err.Error()))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlerError{Message: "Cannot compute jwt token", Reason: err.Error()})
		return
	}

	fetchedUser.JWTToken = jwt
	fetchedUser.RefreshToken = rt

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedUser)
	elapsed := time.Since(start).Milliseconds()
	getLogger().Info(fmt.Sprintf("Request handled in %d ms", elapsed))

}

// RefreshToken godoc
// @Summary Referesh JWT Token using the refresh token
// @Description Use referesh token to obtain new jwt token
// @Tags Refresh JWT Token
// @Accept  json
// @Produce  json
// @Param user body TokenRequestBody true "RefreshToken"
// @Success 200 {object} TokenResponseBody
// @Failure 404 {object} handlerError
// @Failure 500 {object} handlerError
// @Failure 403 {object} handlerError
// @Router /token [post]
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	getLogger().Info("Handling refresh token")
	tokenRequest := TokenRequestBody{}
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(handlerError{Message: "Error parsing refresh token", Reason: "Refresh Token"})
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["userId"].(string)
		user, status := postgres.GetUserByUserId(userId)
		if status.Code != http.StatusOK {
			getLogger().Error(status.Message)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(handlerError{Message: "Error parsing refresh token", Reason: "Refresh Token"})
			return
		}
		jwtToken, rt, err := TokenizeUser(user)
		if err != nil {
			getLogger().Error(err.Error())

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(handlerError{Message: "Error parsing refresh token", Reason: "Refresh Token"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(TokenResponseBody{JWTToken: jwtToken, RefreshToken: rt})
		elapsed := time.Since(start).Milliseconds()
		getLogger().Info(fmt.Sprintf("Request handled in %d ms", elapsed))
	} else {
		getLogger().Error("cannot parse jwt claims")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(handlerError{Message: "Error parsing claims", Reason: "Refresh Token"})
	}
}

// FindUserByUsername godoc
// @Summary Get a list of users with username that contains the string provided
// @Description Search for users with username that is partial matching the string
// @Tags Find user by partial matching
// @Produce json
// @Param username query string true "Term for partial matching username"
// @Success 200 {object} transferObject.UsersResponseBody
// @Failure 404 {object} handlerError
// @Failure 500 {object} handlerError
// @Router /users [GET]
func FindUserByUsername(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	getLogger().Info("Handling find user by matching username")
	username := r.URL.Query()["username"][0]
	users, status := postgres.SearchUsersByUsername(username)

	if status.Code != http.StatusOK {
		getLogger().Error(fmt.Sprintf("Error when finding user %s", status.Message))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status.Code)
		json.NewEncoder(w).Encode(handlerError{Message: status.Message, Reason: status.Reason})
		return
	}
	for i := range users {
		users[i].Password = ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UsersResponseBody{Users: users})
	elapsed := time.Since(start).Milliseconds()
	getLogger().Info(fmt.Sprintf("Request handled in %d ms", elapsed))
}

// CountByUsername godoc
// @Summary Get count of users in the database by username
// @Description Return total count of users with the given username
// @Tags Count user by username
// @Produce json
// @Param username query string true "exact username string to search"
// @Success 200 {object} UserCountResponseBody
// @Failure 404 {object} handlerError
// @Failure 500 {object} handlerError
// @Router /countByUsername [GET]
func CountByUsername(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	getLogger().Info("Counting user by username")
	username := r.URL.Query()["username"][0]
	if len(username) == 0 {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlerError{Message: "displayName not in param", Reason: "MISSING_PARAM"})
	}
	username = strings.TrimSpace(username)
	count, status := postgres.CountByUsername(username)

	if status.Code != http.StatusOK {
		getLogger().Error(fmt.Sprintf("Error when searching username %s", status.Message))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status.Code)
		json.NewEncoder(w).Encode(handlerError{Message: status.Message, Reason: status.Reason})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserCountResponseBody{Count: count})
	elapsed := time.Since(start).Milliseconds()
	getLogger().Info(fmt.Sprintf("Request handled in %d ms", elapsed))
}

// CountByDisplayName godoc
// @Summary Get count of users in the database by displayName
// @Description Return total count of users with the given display name
// @Tags Count user by display name
// @Produce json
// @Param displayName query string true "exact display name string to search"
// @Success 200 {object} UserCountResponseBody
// @Failure 404 {object} handlerError
// @Failure 500 {object} handlerError
// @Router /countByDisplayName [GET]
func CountByDisplayName(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	getLogger().Info("Counting user by displayName")
	dname := r.URL.Query()["displayName"][0]
	if len(dname) == 0 {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlerError{Message: "displayName not in param", Reason: "MISSING_PARAM"})
		return
	}
	dname = strings.TrimSpace(dname)
	count, status := postgres.CountByDisplayName(dname)

	if status.Code != http.StatusOK {
		getLogger().Error(fmt.Sprintf("Error when searching displayName %s", status.Message))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status.Code)
		json.NewEncoder(w).Encode(handlerError{Message: status.Message, Reason: status.Reason})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserCountResponseBody{Count: count})
	elapsed := time.Since(start).Milliseconds()
	getLogger().Info(fmt.Sprintf("Request handled in %d ms", elapsed))
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
	if err != nil {
		return "", "", err
	}
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
