package main

import (
	"reflect"
	ss "src/ex00"
	"sync"
	"testing"
)

func TestSleepSort(t *testing.T) {
	slice := []int{2, 3, 6, 1}
	expected := []int{1, 2, 3, 6}
	actual := make([]int, len(slice))
	c := ss.SleepSort(slice)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < len(slice); i++ {
			actual[i] = <-c
		}
		wg.Done()
	}()
	wg.Wait()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("sleepSort(%v) = %v, want %v", slice, actual, expected)
	}
}
