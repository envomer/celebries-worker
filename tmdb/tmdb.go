package tmdb

import (
	"any-days.com/celebs/logger"
	"encoding/json"
	"fmt"
	"github.com/evolidev/evoli/framework/filesystem"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const ApiKey = "9298b7e03f223cc27836c6d8e23fd5e0"

var log = logger.AppLog

type Person struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	ProfilePath        string  `json:"profile_path"`
	Popularity         float64 `json:"popularity"`
	Adult              bool    `json:"adult"`
	KnownForDepartment string  `json:"known_for_department"`
}

type PeopleResponse struct {
	Page         int       `json:"page"`
	TotalResults int       `json:"total_results"`
	TotalPages   int       `json:"total_pages"`
	Results      []*Person `json:"results"`
}

func FetchPeople(page, limit int) {
	log.Debug("Fetch people from TMDB (page: %d)", page)

	i := 1
	totalPages := 500
	for i <= limit {
		log.Debug("Fetch page %d out of %d", i, totalPages)
		response, err := GetPopularPeople(i)

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

func GetPopularPeople(page int) (*PeopleResponse, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/person/popular?api_key=%s&language=en-US&page=%d", ApiKey, page)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var data *PeopleResponse
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func SavePage(page int, response *PeopleResponse, err error) {
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		panic(err)
	}

	// file id with 3 digits
	fileId := fmt.Sprintf("%03d", page)
	file := fmt.Sprintf("data/page-%s.json", fileId)

	filesystem.Write(file, string(jsonData))
}

func FusePeople() {
	// people
	var people []*Person
	// get all the files within data folder
	filepath.Walk("data", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		data := filesystem.Read(path)
		var p *PeopleResponse
		json.Unmarshal([]byte(data), &p)
		people = append(people, p.Results...)

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

func DownloadAllPeople() {
	base := "https://files.tmdb.org/p/exports/person_ids_%02d_%02d_%04d.json.gz"
	currentDay := time.Now().Day() - 1
	currentMonth := time.Now().Month()
	currentYear := time.Now().Year()
	url := fmt.Sprintf(base, currentDay, currentMonth, currentYear)

	log.Debug("Downloading people from %s", url)
}
