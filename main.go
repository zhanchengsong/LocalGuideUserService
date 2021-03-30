package main

import (
	log "github.com/sirupsen/logrus"
	postgress "github.io/zhanchengsong/LocalGuideUserService/persistence"
)

func main() {
	db, err := postgress.ConnectDB()

	if err != nil {
		log.Error(err.Error())
	}
	_, err = postgress.CheckUserByUsername("TestUser")

	if err != nil {
		log.Error(err.Error())
	}
	defer db.Close()
}
