package tmdb

import (
	"any-days.com/celebs/db"
	"any-days.com/celebs/logger"
	"any-days.com/celebs/model"
	tmdb "github.com/cyruzin/golang-tmdb"
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

func FetchPeople(page int) {
	log.Debug("Fetch people from TMDB (page: %d)", page)
	client := GetClient()

	response, err := client.GetPersonPopular(map[string]string{
		"page": "1",
	})

	if err != nil {
		panic(err)
	}

	log.Debug("Total pages: %d", response.TotalPages)

	// get all tmdb ids from db
	var people []model.Person
	db.Db().Find(&people)

	ids := map[int]bool{}
	for _, person := range people {
		ids[int(person.ID)] = true
	}

	for _, p := range response.Results {
		if _, ok := ids[int(p.ID)]; ok {
			log.Debug("Person already exists: %s (%d)", p.Name, p.ID)
			continue
		}

		person := &model.Person{
			TmdbID:      int(p.ID),
			Name:        p.Name,
			ProfilePath: p.ProfilePath,
			Popularity:  p.Popularity,
			Adult:       p.Adult,
		}

		result := db.Db().Save(person)
		if result.Error != nil {
			log.Error("Failed to save person %s (%d): %s", p.Name, p.ID, result.Error)
			continue
		}

		log.Debug("Person saved: %s (%d)", p.Name, p.ID)
	}
}
