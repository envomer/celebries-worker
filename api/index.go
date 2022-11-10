package handler

import (
	"any-days.com/celebs/db"
	"any-days.com/celebs/logger"
	"any-days.com/celebs/model"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

//go:embed db.db
var content embed.FS

var log = logger.WebLog

func Handler(w http.ResponseWriter, r *http.Request) {

	//exec ls -lah
	cmd := exec.Command("ls", "-lah")
	out, err := cmd.Output()
	if err != nil {
		log.Error("Failed to execute command: %s", err)
	}

	fmt.Fprintf(w, string(out))

	fmt.Fprintf(w, "<h1>Hello from Go!</h1>")

	GetRandomPeople(w, r)
}

func GetRandomPeople(w http.ResponseWriter, r *http.Request) {

	var people []*model.Person
	db.Db().Order("RANDOM()").Limit(10).Find(&people)

	toJson(w, people)
}

func toJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(data)
	if err != nil {
		log.Error("Failed to encode json: %s", err)
	}

	w.WriteHeader(http.StatusOK)

}
