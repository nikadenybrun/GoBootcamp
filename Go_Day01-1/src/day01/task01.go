package main

import (
	"Golang/Go_Day00-1/day01/jsonreader"
	"Golang/Go_Day00-1/day01/models"
	"Golang/Go_Day00-1/day01/xmlreader"
	"flag"
	"fmt"
	"log"
	"strings"
)

type Pair struct {
	count, unit string
}

func compare(obj1, obj2 models.Recipes) {
	oldRecipeNames := make(map[string]bool)
	newRecipeNames := make(map[string]bool)

	for _, recipe := range obj1.Cake {
		oldRecipeNames[recipe.Name] = true

	}
	for _, recipe := range obj2.Cake {
		newRecipeNames[recipe.Name] = true
	}
	for _, new_recipe := range obj2.Cake {
		if _, inMap := oldRecipeNames[new_recipe.Name]; inMap == false {
			fmt.Printf("ADDED cake \"%s\"\n", new_recipe.Name)
		}

	}
	for _, old_recipe := range obj1.Cake {
		if _, inMap := newRecipeNames[old_recipe.Name]; inMap == false {
			fmt.Printf("REMOVED cake \"%s\"\n", old_recipe.Name)
		}

	}
	for _, old_recipe := range obj1.Cake {
		for _, new_recipe := range obj2.Cake {
			if old_recipe.Name == new_recipe.Name {
				if old_recipe.Time != new_recipe.Time {
					fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", new_recipe.Name, new_recipe.Time, old_recipe.Time)
				}
				oldIngName := make(map[string]Pair) //для названия ингридиентов
				newIngName := make(map[string]Pair)
				for _, ing := range old_recipe.Ingridients {
					oldIngName[ing.Name] = Pair{ing.Count, ing.Unit}
				}
				for _, ing := range new_recipe.Ingridients {
					newIngName[ing.Name] = Pair{ing.Count, ing.Unit}
				}
				for _, new_ing := range new_recipe.Ingridients {
					if _, inMap := oldIngName[new_ing.Name]; inMap == false {
						fmt.Printf("ADDED ingredient \"%s\" for cake  \"%s\"\n", new_ing.Name, new_recipe.Name)
					} else {
						if new_ing.Unit != oldIngName[new_ing.Name].unit && new_ing.Unit != "" && oldIngName[new_ing.Name].unit != "" {
							fmt.Printf("CHANGED unit for ingredient \"%s\" for cake  \"%s\" - \"%s\" instead of \"%s\"\n", new_ing.Name, new_recipe.Name, new_ing.Unit, oldIngName[new_ing.Name].unit)

						} else if new_ing.Unit != oldIngName[new_ing.Name].unit && new_ing.Unit == "" && oldIngName[new_ing.Name].unit != "" {
							fmt.Printf("REMOVED unit \"%s\" for ingredient \"%s\" for cake \"%s\"\n", oldIngName[new_ing.Name].unit, new_ing.Name, new_recipe.Name)

						} else if new_ing.Unit != oldIngName[new_ing.Name].unit && new_ing.Unit != "" && oldIngName[new_ing.Name].unit == "" {
							fmt.Printf("ADDED unit \"%s\" for ingredient \"%s\" for cake \"%s\"\n", new_ing.Unit, new_ing.Name, new_recipe.Name)
						}
					}
				}

				for _, old_ing := range old_recipe.Ingridients {
					if _, inMap := newIngName[old_ing.Name]; inMap == false {
						fmt.Printf("REMOVED ingredient \"%s\" for cake  \"%s\"\n", old_ing.Name, old_recipe.Name)
					} else {
						if old_ing.Count != newIngName[old_ing.Name].count && old_ing.Count != "" && newIngName[old_ing.Name].count != "" {
							fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake  \"%s\" - \"%s\" instead of \"%s\"\n", old_ing.Name, new_recipe.Name, newIngName[old_ing.Name].count, old_ing.Count)

						} else if old_ing.Count != newIngName[old_ing.Name].count && old_ing.Count == "" && newIngName[old_ing.Name].count != "" {
							fmt.Printf("ADDED unit count \"%s\" for ingredient \"%s\" for cake  \"%s\" - \"%s\" instead of \"%s\"\n", newIngName[old_ing.Name].count, old_ing.Name, new_recipe.Name)

						} else if old_ing.Count != newIngName[old_ing.Name].count && old_ing.Count != "" && newIngName[old_ing.Name].count == "" {
							fmt.Printf("REMOVED unit count \"%s\" for ingredient \"%s\" for cake  \"%s\"\n", newIngName[old_ing.Name].count, old_ing.Name, new_recipe.Name)
						}
					}
				}

			}
		}
	}
}

func main() {
	old_filename := flag.String("old", "", "input file name")
	new_filename := flag.String("new", "", "input file name")
	flag.Parse()
	if *old_filename == "" || *new_filename == "" {
		log.Fatal("no input file specified")
	}
	idx1 := strings.Index(*old_filename, ".")
	idx2 := strings.Index(*new_filename, ".")
	rs1 := (*old_filename)[idx1+1:]
	rs2 := (*new_filename)[idx2+1:]
	xml := xmlreader.XML{}
	json := jsonreader.JSON{}
	if rs1 == "xml" && rs2 == "json" {
		_, s, err := xml.DBReader(*old_filename)
		if err != nil {
			log.Fatal(err)
		}
		_, s1, err := json.DBReader(*new_filename)
		if err != nil {
			log.Fatal(err)
		}
		compare(*s, *s1)
		// fmt.Printf("%s", data)

	} else if rs1 == "json" && rs2 == "xml" {
		_, s, err := json.DBReader(*old_filename)
		if err != nil {
			log.Fatal(err)
		}
		_, s1, err := xml.DBReader(*new_filename)
		if err != nil {
			log.Fatal(err)
		}
		compare(*s, *s1)
		// fmt.Printf("%s", data)

	} else {
		fmt.Println("invalid filename")
	}
}
