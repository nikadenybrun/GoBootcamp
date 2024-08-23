package store

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/olivere/elastic"
	// elastic "gopkg.in/olivere/elastic.v7"
)

type SortJSON struct {
	Sort []Sort `json:"sort"`
}

type Sort struct {
	GeoDistance GeoDistance `json:"_geo_distance"`
}

type GeoDistance struct {
	Location       elastic.GeoPoint `json:"location"`
	Order          string           `json:"order"`
	Unit           string           `json:"unit"`
	Mode           string           `json:"mode"`
	DistanceType   string           `json:"distance_type"`
	IgnoreUnmapped bool             `json:"ignore_unmapped"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Place struct {
	ID       int              `json:"id"`
	Name     string           `json:"name"`
	Address  string           `json:"address"`
	Phone    string           `json:"phone"`
	Location elastic.GeoPoint `json:"location"`
}

type Types struct {
	Places []Place `json:"_source"`
}

func New() *Types {
	return &Types{}
}

func (ty *Types) GetPlaces(limit int, lat, lon string) ([]Place, error) {
	var buf bytes.Buffer
	if lat == "" || lon == "" {
		return nil, errors.New("empty lat and lon")
	}
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}

	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return nil, err
	}

	lonFloat, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return nil, err
	}

	sort := SortJSON{Sort: []Sort{{
		GeoDistance{
			Location:       elastic.GeoPoint(Location{Lon: lonFloat, Lat: latFloat}),
			Order:          "asc",
			Unit:           "m",
			Mode:           "max",
			DistanceType:   "arc",
			IgnoreUnmapped: true,
		},
	}}}

	if err = json.NewEncoder(&buf).Encode(sort); err != nil {
		return nil, err
	}

	req := esapi.SearchRequest{
		Index:      []string{"places"},
		Size:       &limit,
		Body:       &buf,
		SearchType: "query_then_fetch",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(ty); err != nil {
		return nil, err
	}

	return ty.Places, nil
}

func (p *Types) UnmarshalJSON(data []byte) error {
	p.Places = p.Places[:0]
	tmpl := struct {
		Hits struct {
			Hits []struct {
				Source struct {
					ID       int              `json:"id"`
					Name     string           `json:"name"`
					Address  string           `json:"address"`
					Phone    string           `json:"phone"`
					Location elastic.GeoPoint `json:"location"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}{}

	if err := json.Unmarshal(data, &tmpl); err != nil {
		return err
	}
	for _, v := range tmpl.Hits.Hits {
		p.Places = append(p.Places, v.Source)
	}

	return nil
}
