package main

import (
	"reflect"
	"src/min_coins"
	"testing"
)

func TestMinCoins2(t *testing.T) {
	// Test case 1: value = 6, coins = [1, 3, 4]
	val := 6
	coins := []int{1, 3, 4}
	expected := []int{3, 3}
	actual := min_coins.MinCoins2(val, coins)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("MinCoins2(%d, %v) = %v, want %v", val, coins, actual, expected)
	}

	// Test case 2: value = 10, coins = [1, 5, 10]
	val = 10
	coins = []int{1, 5, 10}
	expected = []int{10}
	actual = min_coins.MinCoins2(val, coins)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("MinCoins2(%d, %v) = %v, want %v", val, coins, actual, expected)
	}

	// Test case 3: value = 13, coins = [1, 5, 10]
	val = 13
	coins = []int{1, 5, 10}
	expected = []int{10, 1, 1, 1}
	actual = min_coins.MinCoins2(val, coins)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("MinCoins2(%d, %v) = %v, want %v", val, coins, actual, expected)
	}

	// Test case 4: value = 0, coins = [1, 5, 10]
	val = 0
	coins = []int{1, 5, 10}
	expected = []int{}
	actual = min_coins.MinCoins2(val, coins)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("MinCoins2(%d, %v) = %v, want %v", val, coins, actual, expected)
	}
}
