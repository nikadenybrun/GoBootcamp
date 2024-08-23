package main

import (
	"Golang/Go_Day00-1/day01/jsonreader"
	"Golang/Go_Day00-1/day01/xmlreader"
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {

	filename := flag.String("f", "", "input file name")
	flag.Parse()
	if *filename == "" {
		log.Fatal("no input file specified")
	}
	idx := strings.Index(*filename, ".")
	rs := (*filename)[idx+1:]
	if rs == "xml" {
		xml := xmlreader.XML{}
		data, _, err := xml.DBReader(*filename)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", data)

	} else if rs == "json" {
		json := jsonreader.JSON{}
		data, _, err := json.DBReader(*filename)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", data)

	} else {
		fmt.Println("invalid filename")
	}

}
