package myhttp

import (
	"Golang/Go_Day03-1/src/store"
	"Golang/Go_Day03-1/src/types"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type HTTP types.Place

func (ht HTTP) HttpShow(es1 *elasticsearch.Client) {
	var store *store.ElasticsearchStore
	store = store.NewElasticsearchStore(es1, "places")
	setSettings(es1)
	places, _, err := store.GetPlaces(16358, 0)

	if err != nil {
		log.Fatal(err)
	}
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}
	var amountOfPages int
	pageSize := 20
	if (len(places)/pageSize)%10 == 0 {
		amountOfPages = len(places) / pageSize
	} else {
		amountOfPages = len(places)/pageSize + 1
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pageString := r.URL.Query().Get("page")
		pageNum, err := strconv.Atoi(pageString)
		if err != nil || pageNum > amountOfPages {
			t, err := template.New("index").Parse(
				`
		<!doctype html>
		<html>
		<head>
			<meta charset="utf-8">
			<title>Error 404</title>
			<meta name="description" content="">
			<meta name="viewport" content="width=device-width, initial-scale=1">
		</head>
		
		<body>
		<h5>ERROR 404</h5>
		<ul>
				<li>
				<div>Invalid 'page' value: '{{ .Value }}'</div>
				</li>
		</ul>
		</body>
		</html>
		
		`)
			if err != nil {
				fmt.Printf("err %e", err)
				http.Error(w, "err 404", http.StatusInternalServerError)
				return
			}
			data := struct {
				Value string
			}{
				pageString,
			}
			err = t.Execute(w, data)
			if err != nil {
				fmt.Printf("err %e", err)
				http.Error(w, "HTTP 400 error", http.StatusInternalServerError)
				return
			}

		} else {
			var (
				pagePlaces [][]types.Place
			)
			for i := 0; i < amountOfPages; i++ {
				if (i+1)*pageSize > len(places) {
					pagePlaces = append(pagePlaces, places[i*pageSize:])
				} else {
					pagePlaces = append(pagePlaces, places[i*pageSize:(i+1)*pageSize])

				}
			}
			t, err := template.New("index").Funcs(funcMap).Parse(`
			<!doctype html>
			<html>
			<head>
				<meta charset="utf-8">
				<title>Places</title>
				<meta name="description" content="">
				<meta name="viewport" content="width=device-width, initial-scale=1">
			</head>
			
			<body>
			<h5>Total: {{len .Places}}</h5>
			<ul>
				{{ range (index .Places (sub .PageNum 1))}}
					<li>
					<div>{{ .Name }}</div>
					<div>{{ .Phone }}</div>
					<div>{{ .Address }}</div>
					</li>
				{{ end }}
			</ul>
			{{ if gt .PageNum 1 }}
				<a href="/?page={{sub .PageNum 1}}">Previous</a>
			{{ end }}
			{{ if lt .PageNum .AmountOfPages }}
				<a href="/?page={{add .PageNum 1}}">Next</a>
			{{ end }}
			<a href="/?page={{.AmountOfPages}}">Last</a>
			</body>
			</html>
			
			`)
			if err != nil {
				fmt.Printf("err %e", err)
				http.Error(w, "fuck0", http.StatusInternalServerError) //err.Error()
				return
			}

			data := struct {
				Places        [][]types.Place
				PageNum       int
				AmountOfPages int
			}{
				pagePlaces,
				pageNum,
				amountOfPages,
			}

			err = t.Execute(w, data)
			if err != nil {
				fmt.Printf("err %e", err)
				http.Error(w, "HTTP 400 error", http.StatusInternalServerError) //err.Error()
				return
			}
		}

	})
	log.Println(http.ListenAndServe(":8888", nil))
}

func setSettings(es *elasticsearch.Client) {

	settings := map[string]interface{}{
		"index": map[string]interface{}{
			"max_result_window": 20000,
		},
	}

	jsonBody, err := json.Marshal(settings)
	if err != nil {
		log.Fatal(err)
	}

	putSettingsReq := esapi.IndicesPutSettingsRequest{
		Index: []string{"places"},
		Body:  bytes.NewReader(jsonBody),
	}

	res, err := putSettingsReq.Do(context.Background(), es)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal("Failed to update index settings")
	}

}
