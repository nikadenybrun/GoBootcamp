package mysearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"Golang/Go_Day03-1/src/ex03/store"
)

type Recommendation struct {
	Name   string        `json:"name"`
	Places []store.Place `json:"places"`
}
type SRCH store.Place

func (s SRCH) Search(index string) {

	http.HandleFunc("/api/recommend", func(w http.ResponseWriter, r *http.Request) {
		lat := r.URL.Query().Get("lat")
		lon := r.URL.Query().Get("lon")

		places, err := store.New().GetPlaces(3, lat, lon)
		if err != nil {
			fmt.Printf("err %e", err)
			http.Error(w, "Cannot gat places", http.StatusInternalServerError)
		}
		recommendation := Recommendation{
			Name:   "Recommendation",
			Places: make([]store.Place, len(places)),
		}

		for i := range len(places) {
			recommendation.Places[i] = places[i]
		}

		w.Header().Set("Content-Type", "application/json")
		marJSON, err := json.Marshal(recommendation)
		if err != nil {
			http.Error(w, "Cannot marshal json", http.StatusInternalServerError)
		}

		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, marJSON, "", "\t")
		if err != nil {
			http.Error(w, "Cannot marshal json", http.StatusInternalServerError)
		}
		w.Write(prettyJSON.Bytes())
	})

	http.ListenAndServe(":8888", nil)
}
