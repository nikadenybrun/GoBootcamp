package parser

import (
	"Golang/Go_Day03-1/src/types"
	"os"
	"strconv"

	"github.com/grailbio/base/tsv"
	"github.com/olivere/elastic"
)

type Store interface {
	GetPlaces(limit int, offset int) ([]types.Place, int, error)
}

type PRS types.Place

func (p PRS) readDataSet(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := tsv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p PRS) createPlace(data []string) (types.Place, error) {
	var place types.Place
	id, err := strconv.Atoi(data[0])
	if err != nil {

		return place, err
	}
	lon, err := strconv.ParseFloat(data[4], 64)
	if err != nil {
		return place, err
	}
	lat, err := strconv.ParseFloat(data[5], 64)
	if err != nil {
		return place, err
	}
	place = types.Place{
		ID:       id,
		Name:     data[1],
		Address:  data[2],
		Phone:    data[3],
		Location: elastic.GeoPoint{Lat: lat, Lon: lon},
	}
	return place, err

}
func (p PRS) ParseCSV(places []*types.Place) ([]*types.Place, error) {

	str_places, err := p.readDataSet("../materials/data.csv")
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(str_places); i++ {
		new_place, err := p.createPlace(str_places[i])
		if err != nil {
			return nil, err
		}
		places = append(places, &new_place)
	}
	return places, nil
}
