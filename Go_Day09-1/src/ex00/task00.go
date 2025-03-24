package ex00

import (
	"time"
)

func SleepSort(slice []int) <-chan int {
	out := make(chan int)
	for _, el := range slice {
		go func(n int) {
			time.Sleep(time.Duration(n) * time.Second)
			out <- n
		}(el)
	}
	return out
}
