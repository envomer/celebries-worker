package handler

import (
	"any-days.com/celebs/logger"
	"encoding/json"
	"github.com/evolidev/evoli/framework/filesystem"
	"math/rand"
	"net/http"
)

////go:embed db.db
//var content embed.FS

var log = logger.WebLog

//func main() {
//	http.HandleFunc("/", Handler)
//	http.ListenAndServe(":8080", nil)
//}

func Handler(w http.ResponseWriter, r *http.Request) {
	GetRandomPeople(w, r)
}

func GetRandomPeople(w http.ResponseWriter, r *http.Request) {
	people := GetRandomPeopleMap(25, []int{})

	toJson(w, people)
}

func GetRandomPeopleMap(limit int, exclude []int) []map[string]any {
	if limit > 100 {
		limit = 100
	}

	filePath := "people.json"
	if !filesystem.Exists(filePath) {
		log.Error("File %s not found", filePath)
		return nil
	}

	data := filesystem.Read(filePath)

	people := make([]map[string]any, 0)
	err := json.Unmarshal([]byte(data), &people)
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

	filtered := []map[string]any{}
	for i := 0; i < len(people); i++ {
		id := int(people[i]["tmdb_id"].(float64))

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
