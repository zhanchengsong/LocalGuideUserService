package main

import (
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.io/zhanchengsong/LocalGuideUserService/docs"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.io/zhanchengsong/LocalGuideUserService/handlers"
	postgress "github.io/zhanchengsong/LocalGuideUserService/persistence"
)

// @title User API
// @version 1.0.0
// @description Service for user registration, login and jwtToken refresh
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host auth.zhancheng.dev
// @BasePath /
func main() {
	db, err := postgress.ConnectDB()

	if err != nil {
		log.Error(err.Error())
	}
	defer db.Close()

	port := "8100"
	router := mux.NewRouter()
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/user", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/users", handlers.FindUserByUsername).Methods("GET")
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	router.HandleFunc("/token", handlers.RefreshToken).Methods("POST")

	log.Info(fmt.Sprintf("Service is up and running on port %s", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
