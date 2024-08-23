package main

import (
	"Golang/Go_Day03-1/src/myhttp"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}
	hp := myhttp.HTTP{}
	hp.HttpShow(es)

}
