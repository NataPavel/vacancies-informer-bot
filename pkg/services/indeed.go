package services

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"vac_informer_tgbot/pkg/database"
	tg_methods "vac_informer_tgbot/pkg/database/telegram_db_methods"
)

type SearchResult struct {
	TotalResults int `json:"totalResults"`
	Results      []struct {
		JobTitle          string `json:"jobtitle"`
		Company           string `json:"company"`
		City              string `json:"city"`
		State             string `json:"state"`
		Country           string `json:"country"`
		Date              string `json:"date"`
		JobKey            string `json:"jobkey"`
		Sponsored         bool   `json:"sponsored"`
		IndeedLink        string `json:"indeedApply"`
		FormattedLocation string `json:"formattedLocation"`
	} `json:"results"`
}

func Indeed(tag string) {
	apiKey := os.Getenv("Indeed_Access_Token")

	// use the API key to make a search request
	searchURL := fmt.Sprintf("https://api.indeed.com/ads/apisearch?publisher=%s&q=%s&l=&sort=date&limit=10&format=json&v=2", apiKey, tag)
	resp, err := http.Get(searchURL)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		errorMessage := fmt.Sprintf("Status code error(Indeed): %d %s", resp.StatusCode, resp.Status)
		tg_methods.ErrorLogger(errorMessage)
		fmt.Println(errorMessage)
	}
	defer resp.Body.Close()

	// decode the JSON response into a SearchResult struct
	var result SearchResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic(err)
	}

	if len(result.Results) == 0 {
		return
	}

	var vacancyTitle string
	var vacancyLink string
	var company string
	var location string

	// iterate over the search results and print information about each job
	for _, job := range result.Results {
		vacancyTitle = fmt.Sprintf("%v", job.JobTitle)
		vacancyLink = fmt.Sprintf("%v", job.IndeedLink)
		company = fmt.Sprintf("%v", job.Company)
		location = fmt.Sprintf("%v", job.FormattedLocation)

		fmt.Printf("Job %s: - %s - %s. %s\n", vacancyTitle, company, location, vacancyLink)
	}

	linkMDHash := md5.Sum([]byte(vacancyLink))
	hash := fmt.Sprintf("%x", linkMDHash)

	text := fmt.Sprintf("🔵Indeed.com - %s🔵:\n• %s;\n• %s;\n• %s;\n• %s",
		tag, vacancyTitle, company, location, vacancyLink)
	fmt.Println(text)

	database.CheckVacancy("Indeed.com", vacancyLink, vacancyTitle, location, company, hash, text)
}
