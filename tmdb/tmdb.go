package tmdb

import (
	"any-days.com/celebs/logger"
	"any-days.com/celebs/model"
	"encoding/json"
	"fmt"
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/evolidev/evoli/framework/filesystem"
	"os"
	"path/filepath"
	"time"
)

const ApiKey = "9298b7e03f223cc27836c6d8e23fd5e0"

var _client *tmdb.Client
var log = logger.AppLog

func GetClient() *tmdb.Client {
	if _client != nil {
		return _client
	}

	c, err := tmdb.Init(ApiKey)
	if err != nil {
		panic(err)
	}

	_client = c
	return _client
}

func FetchPeople(page, limit int) {
	log.Debug("Fetch people from TMDB (page: %d)", page)
	client := GetClient()

	i := 1
	totalPages := 500
	for i <= limit {
		log.Debug("Fetch page %d out of %d", i, totalPages)
		response, err := client.GetPersonPopular(map[string]string{
			"page": fmt.Sprintf("%d", i),
		})

		if err != nil {
			log.Error("Failed to fetch page %d: %s", i, err)
			time.Sleep(5 * time.Second)
			log.Debug("Retry page %d", i)
			continue
		}

		totalPages = int(response.TotalPages)
		SavePage(i, response, err)

		if i >= int(response.TotalPages) {
			log.Debug("No more pages")
			break
		}

		// sleep every 4 pages
		if i%10 == 0 {
			log.Debug("Sleeping for 2 seconds...")
			time.Sleep(2 * time.Second)
		}

		i++
	}

	// get all tmdb ids from db
	//var people []model.Person
	//db.Db().Find(&people)
	//
	//ids := map[int]bool{}
	//for _, person := range people {
	//	ids[int(person.ID)] = true
	//}

}

func SavePage(page int, response *tmdb.PersonPopular, err error) {
	actors := []*model.Person{}
	for _, p := range response.Results {
		//if _, ok := ids[int(p.ID)]; ok {
		//	log.Debug("Person already exists: %s (%d)", p.Name, p.ID)
		//	continue
		//}

		person := &model.Person{
			//TmdbID:      int(p.ID),
			ID:          uint(p.ID),
			Name:        p.Name,
			ProfilePath: p.ProfilePath,
			Popularity:  p.Popularity,
			Adult:       p.Adult,
		}

		actors = append(actors, person)

		//result := db.Db().Save(person)
		//if result.Error != nil {
		//	log.Error("Failed to save person %s (%d): %s", p.Name, p.ID, result.Error)
		//	continue
		//}

		//log.Debug("Person saved: %s (%d)", p.Name, p.ID)
	}

	jsonData, err := json.Marshal(actors)
	if err != nil {
		panic(err)
	}

	// file id with 3 digits
	fileId := fmt.Sprintf("%03d", page)
	// save to page-1.json
	file := fmt.Sprintf("data/page-%s.json", fileId)
	filesystem.Write(file, string(jsonData))
}

func FusePeople() {
	// people
	people := []map[string]any{}
	// get all the files within data folder
	filepath.Walk("data", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		data := filesystem.Read(path)
		var p []map[string]any
		json.Unmarshal([]byte(data), &p)
		people = append(people, p...)

		return nil
	})

	// save to people.json
	jsonData, err := json.Marshal(people)
	if err != nil {
		log.Error("Failed to marshal people: %s", err)
		return
	}

	filesystem.Write("api/people.json", string(jsonData))
}
