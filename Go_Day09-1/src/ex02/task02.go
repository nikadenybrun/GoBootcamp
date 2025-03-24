package ex02

import (
	"sync"
)

func Multiplex(channels ...chan interface{}) chan interface{} {
	out := make(chan interface{})

	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(len(channels))
		for _, ch := range channels {
			go func(ch chan interface{}) {
				for msg := range ch {
					out <- msg
				}
				wg.Done()
			}(ch)
		}
		wg.Wait()
		close(out)
	}()

	return out
}
