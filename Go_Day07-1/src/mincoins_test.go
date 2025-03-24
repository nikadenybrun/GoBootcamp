package main

import (
	"reflect"
	"src/min_coins"
	"testing"
)

func TestMinCoins(t *testing.T) {
	val := 6
	coins := []int{1, 3, 4}
	expected := []int{3, 3}
	actual := min_coins.MinCoins(val, coins)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("MinCoins(%d, %v) = %v, want %v", val, coins, actual, expected)
	}

	val = 10
	coins = []int{1, 5, 10}
	expected = []int{10}
	actual = min_coins.MinCoins(val, coins)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("MinCoins(%d, %v) = %v, want %v", val, coins, actual, expected)
	}

	val = 13
	coins = []int{1, 5, 10}
	expected = []int{10, 1, 1, 1}
	actual = min_coins.MinCoins(val, coins)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("MinCoins(%d, %v) = %v, want %v", val, coins, actual, expected)
	}

	val = 0
	coins = []int{1, 5, 10}
	expected = []int{}
	actual = min_coins.MinCoins(val, coins)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("MinCoins(%d, %v) = %v, want %v", val, coins, actual, expected)
	}
}
