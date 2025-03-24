// min_coins provides 2 functions that will count the minimal amount of coins.
package min_coins

import "math"

// minCoins is the given function seraching for the minimal amount of coins.
func MinCoins(val int, coins []int) []int {
	res := make([]int, 0)
	i := len(coins) - 1
	for i >= 0 {
		for val >= coins[i] {
			val -= coins[i]
			res = append(res, coins[i])
		}
		i -= 1
	}
	return res
}

// minCoins2 is a new version of the minCoins function.
func MinCoins2(val int, coins []int) []int {
	dp := make([]int, val+1)
	for i := 1; i <= val; i++ {
		dp[i] = math.MaxInt32
	}
	dp[0] = 0

	for _, coin := range coins {
		for i := coin; i <= val; i++ {
			if dp[i-coin]+1 < dp[i] {
				dp[i] = dp[i-coin] + 1
			}
		}
	}

	if dp[val] == math.MaxInt32 {
		return []int{}
	}

	res := make([]int, 0)
	i := val
	for _, coin := range coins {
		for i >= coin && dp[i] == dp[i-coin]+1 {
			res = append(res, coin)
			i -= coin
		}
	}

	return BubbleSort(res)
}

// bublesort function for sorting the slice
func BubbleSort(res []int) []int {
	n := len(res)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if res[j] < res[j+1] {
				res[j], res[j+1] = res[j+1], res[j]
			}
		}
	}
	return res
}
