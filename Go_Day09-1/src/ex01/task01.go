package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func SigTerm(done chan struct{}) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(done)
	}()
}

// crawlWeb возвращает канал для результатов парсинга веб-страниц.
func crawlWeb(in chan string, done chan struct{}) chan *string {
	out := make(chan *string)
	sem := make(chan bool, 8)
	go func() {
		var wg sync.WaitGroup
		for url := range in {
			sem <- true
			wg.Add(1)
			go func() {
				getData(done, &wg, out, url)
				<-sem
			}()
		}
		wg.Wait()
		close(out)
	}()

	return out
}

func getData(done chan struct{}, wg *sync.WaitGroup, out chan *string, url string) {

	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		out <- nil
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		out <- nil
		return
	}
	ptr := (string(body))
	select {
	case out <- &ptr:
	case <-done:
		close(out)
		os.Exit(0)
	}

}

func putData(in chan string, done chan struct{}) {
	for i := 11; i < 20; i++ {
		time.Sleep(1 * time.Second)
		url := "https://the-one-api.dev/character/" + fmt.Sprintf("%d", i)
		select {
		case in <- url:
		case <-done:
			close(in)
			os.Exit(0)
		}
	}

	close(in)
}
func main() {
	in := make(chan string)
	done := make(chan struct{})

	SigTerm(done)
	go putData(in, done)
	out := crawlWeb(in, done)
	for result := range out {
		if result != nil {
			fmt.Println(*result)
		} else {
			fmt.Println("Error")
		}
	}
}
