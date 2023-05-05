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

func Dou(tag string) {
	url := fmt.Sprintf("https://jobs.dou.ua/vacancies/?search=%s", tag)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errorMessage := fmt.Sprintf("Status code error(Dou): %d %s", resp.StatusCode, resp.Status)
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

	doc.Find(".lt li .vacancy .title").Each(func(i int, s *goquery.Selection) {
		vacancyTitle = s.Find(".vt").Text()
		vacancyLink, _ = s.Find(".vt").Attr("href")
		company = s.Find(".company").Text()
		location = s.Find(".cities").Text()
		fmt.Printf("Job %d: %s - %s - %s. %s\n", i+1, vacancyTitle, company, location, vacancyLink)
	})

	linkMDHash := md5.Sum([]byte(vacancyLink))
	hash := fmt.Sprintf("%x", linkMDHash)

	text := fmt.Sprintf("🔵Dou - %s🔵:\n• %s;\n• %s;\n• %s;\n• %s",
		tag, vacancyTitle, company, location, vacancyLink)
	fmt.Println(text)

	database.CheckVacancy("DOU.ua", vacancyLink, vacancyTitle, location, company, hash, text)
}
