package services

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"

	"vac_informer_tgbot/pkg/database"
	tg_methods "vac_informer_tgbot/pkg/database/telegram_db_methods"

	"github.com/PuerkitoBio/goquery"
)

func Jooble(tag string) {
	url := fmt.Sprintf("https://ua.jooble.org/SearchResult?ukw=%s", tag)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errorMessage := fmt.Sprintf("Status code error(Jooble): %d %s", resp.StatusCode, resp.Status)
		//database.ErrorLogger(errorMessage)
		tg_methods.ErrorLogger(errorMessage)
		fmt.Println(errorMessage)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var vacancyTitle string
	var vacancyLink string
	var company string
	var location string

	doc.Find("main ._18fg4O article").Each(func(i int, s *goquery.Selection) {
		vacancyTitle = s.Find("a").Text()
		vacancyLink, _ = s.Find("a").Attr("href")
		company = s.Find(".Ya0gV9").Text()
		location = s.Find("._2_Ab4T").Text()
		fmt.Printf("Job %d: %s - %s - %s. %s\n", i+1, vacancyTitle, company, location, vacancyLink)
	})

	linkMDHash := md5.Sum([]byte(vacancyLink))
	hash := fmt.Sprintf("%x", linkMDHash)

	text := fmt.Sprintf("🔵Jooble - %s🔵:\n• %s;\n• %s;\n• %s;\n• %s",
		tag, vacancyTitle, company, location, vacancyLink)
	fmt.Println(text)

	database.CheckVacancy("Jooble.org", vacancyLink, vacancyTitle, location, company, hash, text)
}
