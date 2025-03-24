package main

import (
	"reflect"
	ss "src/ex02"
	"strconv"
	"testing"
)

func TestMultiplex(t *testing.T) {

	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})

	out := ss.Multiplex(ch1, ch2, ch3)
	go func() {
		ch1 <- "Message 1"
		ch1 <- "Message 2"
		close(ch1)
	}()

	go func() {
		ch2 <- "Message 3"
		ch2 <- "Message 4"
		close(ch2)
	}()

	go func() {
		ch3 <- "Message 5"
		ch3 <- "Message 6"
		close(ch3)
	}()
	actual := make(map[interface{}]int)
	expected := make(map[interface{}]int)
	for i := 0; i < 6; i++ {
		expected["Message "+strconv.Itoa(i+1)]++
	}
	for msg := range out {
		actual[msg]++
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("error")
	}

}
