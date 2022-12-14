package handler

import (
	"any-days.com/celebs/logger"
	"embed"
	"encoding/json"
	"github.com/evolidev/evoli/framework/use"
	"github.com/spf13/cast"
	"io"
	"math"
	"net/http"
	"strings"
)

//go:embed people.json
var content embed.FS

var log = logger.WebLog

func Handler(w http.ResponseWriter, r *http.Request) {
	// get param page from url
	page := r.URL.Query().Get("page")

	// handle request
	if page == "popular" {
		popular(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func popular(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Exclude string `json:"exclude"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to read request body: %s", err)
		return
	}

	limit := cast.ToInt(r.URL.Query().Get("limit"))
	if limit == 0 || limit > 100 {
		limit = 100
	}

	use.JsonDecodeStruct(string(body), &request)

	// explode exclude string
	excludeString := strings.Split(request.Exclude, ",")

	people := GetRandomPeopleMap(limit, excludeString)

	toJson(w, people)
}

func GetRandomPeopleMap(limit int, exclude []string) map[string]any {
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
	//for i := range people {
	//	j := rand.Intn(i + 1)
	//	people[i], people[j] = people[j], people[i]
	//}

	excludeIds := make(map[int]bool)
	for _, id := range exclude {
		excludeIds[cast.ToInt(id)] = true
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

	totalResults := len(people) - len(exclude)
	return map[string]any{
		"results":       filtered,
		"total_results": totalResults,
		"page":          1,
		"total_pages":   math.Ceil(float64(totalResults / limit)),
	}
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
