package main

import (
	"Golang/Go_Day03-1/src/myapi"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}
	a := myapi.API{}
	a.ApiShow(es, 13649)

}
