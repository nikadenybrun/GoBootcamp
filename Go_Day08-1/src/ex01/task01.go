package main

import (
	"fmt"
	"reflect"
)

type UnknownPlant struct {
	FlowerType string
	LeafType   string
	Color      int `color_scheme:"rgb"`
}

type AnotherUnknownPlant struct {
	FlowerColor int
	LeafType    string
	Height      int `unit:"inches"`
}

type Plants interface {
	printPlant() interface{}
}

func (p UnknownPlant) printPlant() interface{} {
	return p
}
func (p AnotherUnknownPlant) printPlant() interface{} {
	return p
}

func describePlant(plant Plants) {
	p := plant.printPlant()
	rv := reflect.ValueOf(p)
	for i := 0; i < rv.NumField(); i++ {
		fmt.Printf("%s: %v\n", rv.Type().Field(i).Name, rv.Field(i).Interface())
	}
}
func main() {
	up := &UnknownPlant{
		FlowerType: "rose",
		LeafType:   "sircle",
		Color:      25500,
	}
	aup := &AnotherUnknownPlant{
		FlowerColor: 2550,
		LeafType:    "square",
		Height:      12,
	}
	describePlant(up)
	describePlant(aup)
}
