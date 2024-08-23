package main

import (
	"Golang/Go_Day03-1/src/mydb"
	"Golang/Go_Day03-1/src/parser"
	"Golang/Go_Day03-1/src/types"
	"log"
)

func main() {
	log.SetFlags(0)
	db := mydb.MYDB{}
	p := parser.PRS{}
	var (
		places []*types.Place
	)
	places, err := p.ParseCSV(places)
	if err != nil {
		log.Fatal(err)
	}
	db.CreateESC(places)

}
