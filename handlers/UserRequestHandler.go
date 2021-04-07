package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	}
	elapsed := time.Since(start).Milliseconds()
	getLogger().Info(fmt.Sprintf("Request handled in %d ms", elapsed))
}
