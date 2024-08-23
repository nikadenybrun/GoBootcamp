package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	old_filename := flag.String("old", "", "input file name")
	new_filename := flag.String("new", "", "input file name")
	flag.Parse()
	if *old_filename == "" || *new_filename == "" {
		log.Fatal("no input file specified")
	}

	file1, err1 := os.Open(*old_filename)
	file2, err2 := os.Open(*new_filename)
	if err1 != nil || err2 != nil {
		log.Fatal(err1)
	}
	defer file1.Close()
	defer file2.Close()

	rd1 := bufio.NewReader(file1)
	rd2 := bufio.NewReader(file2)

	oldLines := make(map[string]bool)
	newLines := make(map[string]bool)

	for {
		line1, err1 := rd1.ReadString('\n')
		if err1 != nil && err1 != io.EOF {
			log.Fatal(err1)
		}
		if err1 == io.EOF {
			break
		}
		oldLines[line1] = true
	}

	for {
		line2, err2 := rd2.ReadString('\n')
		if err2 != nil && err2 != io.EOF {
			log.Fatal(err2)
		}
		if err2 == io.EOF {
			break
		}
		newLines[line2] = true
	}

	for line := range oldLines {
		if !newLines[line] {
			fmt.Printf("REMOVED %s", line)
		}
	}

	for line := range newLines {
		if !oldLines[line] {
			fmt.Printf("ADDED %s", line)
		}
	}
}
