package xmlreader

import (
	"Golang/Go_Day00-1/day01/models"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
)

type XML models.Recipes

func (f XML) DBReader(filename string) ([]byte, *models.Recipes, error) {
	var s models.Recipes
	fileCont, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, &s, err
	}
	if err := xml.Unmarshal(fileCont, &s); err != nil {
		return nil, &s, err
	}

	data1, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return nil, &s, err
	}
	return data1, &s, nil
}
