package controller

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Word struct {
	Tr       string `json:"tr"`
	En       string `json:"en"`
	Category string `json:"kategori"`
	Type     string `json:"tur"`
}

type Words []Word

var words Words

func FetchFromTureng(query string) (Words, error) {
	res, err := http.Get("https://tureng.com/tr/turkce-ingilizce/" + query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, errors.New("status code error: " + res.Status)
	}

	// Load the HTML document from tureng.com
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Find the the all rows. Every row equals a word on tureng.com
	tr := doc.Find("#englishResultsTable").First().Find("tbody tr")
	tr.Each(func(i int, str *goquery.Selection) {
		var word Word
		tableColumns := str.Find("td")

		wordType := convertType(tableColumns.Eq(2).Find("i").Text())
		word.Category = tableColumns.Eq(1).Text()
		word.Type = wordType
		word.En = tableColumns.Eq(2).Find("a").Text()
		word.Tr = tableColumns.Eq(3).Find("a").Text()

		if word.En != "" || word.Tr != "" {
			words = append(words, word)
		}
	})

	defer func() {
		words = []Word{}
	}()

	return words, nil
}

func convertType(wordType string) string {
	wordType = strings.Trim(wordType, " ")
	switch wordType {
	case "i.":
		return "isim"
	case "f.":
		return "fiil"
	case "zf.":
		return "zarf"
	case "ünl.":
		return "ünlem"
	default:
		return ""
	}
}
