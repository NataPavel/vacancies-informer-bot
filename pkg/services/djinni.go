package services

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"strings"

	"vac_informer_tgbot/pkg/database"
	tg_methods "vac_informer_tgbot/pkg/database/telegram_db_methods"

	"github.com/PuerkitoBio/goquery"
)

func Djinni(tag string) {
	url := fmt.Sprintf("https://djinni.co/jobs/?keywords=%s&all-keywords=&any-of-keywords=&exclude-keywords=", tag)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errorMessage := fmt.Sprintf("Status code error(Djinni): %d %s", resp.StatusCode, resp.Status)
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

	doc.Find("body .list-jobs li").Each(func(i int, s *goquery.Selection) {
		vacancyTitle = s.Find(".profile").Find("span").Text()
		vacancyLink, _ = s.Find(".profile").Attr("href")
		company = strings.TrimSpace(s.Find(".list-jobs__details__info > a").First().Text())
		location = strings.TrimSpace(s.Find(".list-jobs__details__info").Find(".location-text").Text())

		company = strings.ReplaceAll(company, " ", "")
		location = strings.ReplaceAll(location, " ", "")

		fmt.Printf("Job %d: %s - %s - %s. https://djinni.co%s\n", i+1, vacancyTitle, company, location, vacancyLink)
	})

	linkMDHash := md5.Sum([]byte(vacancyLink))
	hash := fmt.Sprintf("%x", linkMDHash)

	text := fmt.Sprintf("🔵Djinni - %s🔵:\n• %s;\n• %s;\n• %s;\n• https://djinni.co%s",
		tag, vacancyTitle, company, location, vacancyLink)
	fmt.Println(text)

	database.CheckVacancy("Djinni.co", vacancyLink, vacancyTitle, location, company, hash, text)
}
