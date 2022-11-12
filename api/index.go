package handler

import (
	"any-days.com/celebs/logger"
	"embed"
	"encoding/json"
	"github.com/evolidev/evoli/framework/use"
	"io"
	"math/rand"
	"net/http"
)

//go:embed people.json
var content embed.FS

var log = logger.WebLog

func Handler(w http.ResponseWriter, r *http.Request) {
	// read request url
	url := r.URL.Path

	// handle request
	if url == "/popular" {
		popular(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func popular(w http.ResponseWriter, r *http.Request) {
	var exclude []int

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to read request body: %s", err)
		return
	}

	use.JsonDecodeStruct(string(body), &exclude)

	GetRandomPeople(w, r, exclude)
}

func GetRandomPeople(w http.ResponseWriter, r *http.Request, exclude []int) {
	people := GetRandomPeopleMap(25, exclude)

	toJson(w, people)
}

func GetRandomPeopleMap(limit int, exclude []int) []map[string]any {
	if limit > 100 {
		limit = 100
	}

	filePath := "people.json"

	data, err := content.ReadFile(filePath)
	if err != nil {
		log.Error("Failed to read file: %s", err)
		return nil
	}

	people := make([]map[string]any, 0)
	err = json.Unmarshal([]byte(data), &people)
	if err != nil {
		log.Error("Failed to unmarshal json: %s", err)
	}

	// randomize
	for i := range people {
		j := rand.Intn(i + 1)
		people[i], people[j] = people[j], people[i]
	}

	excludeIds := make(map[int]bool)
	for _, id := range exclude {
		excludeIds[id] = true
	}

	var filtered []map[string]any
	for i := 0; i < len(people); i++ {
		id := int(people[i]["id"].(float64))

		if _, ok := excludeIds[id]; ok {
			continue
		}

		if len(filtered) >= limit {
			break
		}

		filtered = append(filtered, people[i])
	}

	return filtered
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
