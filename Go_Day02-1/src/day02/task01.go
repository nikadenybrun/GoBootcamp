package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

type Info struct {
	flags string
	files []string
}

func (f *Info) parseFlag() *Info {
	if len(os.Args) > 1 {
		var i int = 1
		if strings.HasPrefix(os.Args[1], "-") {
			f.flags = os.Args[1]
			i = 2
		} else {
			f.flags = "-w"
		}
		for ; i < len(os.Args); i++ {
			f.files = append(f.files, os.Args[i])
		}
	} else {
		fmt.Println("No path provided")
	}
	return f
}

func main() {
	info := Info{}
	info.parseFlag()
	type pair struct {
		filename string
		info     int
	}
	res := make([]pair, len(info.files))
	var wg sync.WaitGroup

	wg.Add(len(info.files))
	for i, filename := range info.files {

		res[i].filename = filename
		go func() {
			defer wg.Done()
			file, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
				return
			}

			scanner := bufio.NewScanner(file)
			if info.flags == "-l" {
				scanner.Split(bufio.ScanLines)
			} else if info.flags == "-m" {
				scanner.Split(bufio.ScanRunes)
			} else if info.flags == "-w" {
				scanner.Split(bufio.ScanWords)
			} else {
				log.Fatal("incorrect fuck")
			}

			for scanner.Scan() {
				res[i].info++
			}

		}()

	}
	wg.Wait()
	for _, elems := range res {
		fmt.Printf("%d\t%s\n", elems.info, elems.filename)
	}
}
