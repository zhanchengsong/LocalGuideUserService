package main

import (
	log "github.com/sirupsen/logrus"
	postgress "github.io/zhanchengsong/LocalGuideUserService/persistence"
)

func main() {
	_, err := postgress.ConnectDB()
	if err != nil {
		log.Error(err.Error())
	}
}
