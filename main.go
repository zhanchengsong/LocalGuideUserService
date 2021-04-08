package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.io/zhanchengsong/LocalGuideUserService/handlers"
	postgress "github.io/zhanchengsong/LocalGuideUserService/persistence"
)

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	db, err := postgress.ConnectDB()

	if err != nil {
		log.Error(err.Error())
	}
	defer db.Close()

	port := os.Getenv("SERVICE_PORT")
	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)
	router.HandleFunc("/user", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	log.Info(fmt.Sprintf("Service is up and running on port %s", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
