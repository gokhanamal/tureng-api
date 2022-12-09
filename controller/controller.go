package controller

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Phrase struct {
	Source   string `json:"source"`
	Target   string `json:"target"`
	Category string `json:"category"`
	Type     string `json:"type"`
}

type Phrases []Phrase

var phrases Phrases

func FetchFromTureng(query string) (Phrases, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MaxVersion: tls.VersionTLS12,
		},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", "https://tureng.com/tr/turkce-ingilizce/"+query, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPad; CPU OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148")
	res, err := client.Do(req)
	//res, err := http.Get("https://tureng.com/tr/turkce-ingilizce/" + query)
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

	// Find the the all rows. Every row equals a phrase on tureng.com
	source := doc.Find("#englishResultsTable").Find("tbody tr")
	source.Each(func(i int, str *goquery.Selection) {
		var phrase Phrase
		tableColumns := str.Find("td")

		secondColumnTypeText := tableColumns.Eq(2).Find("i").Text()

		if secondColumnTypeText != "" {
			phrase.Type = convertType(secondColumnTypeText)
		} else {
			phrase.Type = convertType(tableColumns.Eq(3).Find("i").Text())
		}

		phrase.Category = tableColumns.Eq(1).Text()
		phrase.Target = tableColumns.Eq(3).Find("a").Text()
		phrase.Source = tableColumns.Eq(2).Find("a").Text()

		if phrase.Target != "" || phrase.Source != "" {
			phrases = append(phrases, phrase)
		}
	})

	defer func() {
		phrases = []Phrase{}
	}()

	return phrases, nil
}

func convertType(phraseType string) string {
	phraseType = strings.Trim(phraseType, " ")
	switch phraseType {
	case "i.":
		return "isim"
	case "f.":
		return "fiil"
	case "zf.":
		return "zarf"
	case "ünl.":
		return "ünlem"
	case "s.":
		return "sıfat"
	default:
		return "unknown"
	}
}
