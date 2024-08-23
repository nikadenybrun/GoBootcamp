package jsonreader

import (
	"Golang/Go_Day00-1/day01/models"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
)

type JSON models.Recipes

func (f JSON) DBReader(filename string) ([]byte, *models.Recipes, error) {
	var x models.Recipes
	fileCont1, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, &x, err
	}
	if err := json.Unmarshal(fileCont1, &x); err != nil {
		return nil, &x, err
	}
	data, err := xml.MarshalIndent(x, "", "    ")
	if err != nil {
		return nil, &x, err
	}

	return data, &x, nil
}
