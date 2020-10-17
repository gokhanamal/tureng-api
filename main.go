package main

import (
	"encoding/json"
	"fmt"
	"github.com/gokhanamal/tureng-api/controller"
	"log"
	"net/http"
	"strings"
)

type Error struct {
	Message string `json:"message, omitempty"`
}

type Response struct {
	Count int `json:"count"`
	Phrases []controller.Phrase `json:"phrases"`
}


func writeJSON(w http.ResponseWriter, v interface{}) {
	jsonData, err := json.Marshal(v)

	if err != nil {
		log.Printf("json write error %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func main() {
	http.HandleFunc("/",  func(w http.ResponseWriter, r *http.Request) {
		err := Error{
			"Please use /translate",
		}
		fmt.Print(err)
		w.WriteHeader(http.StatusNotFound)
		writeJSON(w, err)
	})

	http.HandleFunc("/translate",  func(w http.ResponseWriter, req *http.Request) {
		queries := req.URL.Query()
		phrase := queries["phrase"]

		if phrase == nil {
			err := Error{
				"Missing query, you should add to your request phrase that want to translate.",
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			writeJSON(w, err)
			return
		}

		response, err := controller.FetchFromTureng(strings.Join(phrase, ""))
		if err != nil {
			err := Error{
				"Someting went wrong while fetching the phrases from Tureng.",
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			writeJSON(w, err)
		}

		w.WriteHeader(http.StatusOK)
		writeJSON(w, Response{Count: len(response), Phrases: response})
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("%s", err)
	}
}
