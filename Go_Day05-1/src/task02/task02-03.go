package task02

import (
	"container/heap"
)

type Present struct {
	Value int
	Size  int
}

type Presents []Present

func (h Presents) Len() int { return len(h) }
func (h Presents) Less(i, j int) bool {
	if h[i].Value > h[j].Value {
		return true
	} else if h[i].Value == h[j].Value {
		if h[i].Size < h[j].Size {
			return true
		}
	}
	return false
}
func (h Presents) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *Presents) Push(x interface{}) {

	*h = append(*h, x.(Present))
}

func (h *Presents) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *Presents) GetNCoolestPresents(n int) []Present {
	var res []Present
	for i := 0; i < n; i++ {
		res = append(res, heap.Pop(h).(Present))
	}
	return res
}

func Knapsack(h []Present, capacity int) []Present {
	n := len(h)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}
	mapa := make(map[Present]int)
	selectedPresents := make([]Present, 0)
	for i := 1; i <= n; i++ {
		for w := 1; w <= capacity; w++ {
			if h[i-1].Size > w {
				dp[i][w] = dp[i-1][w]
			} else {
				if dp[i-1][w] < dp[i-1][w-h[i-1].Size]+h[i-1].Value {
					mapa[h[i-1]]++
				} else if dp[i-1][w] > dp[i-1][w-h[i-1].Size]+h[i-1].Value {

					mapa[h[i-1]] = 0
				}
				dp[i][w] = max(dp[i-1][w], dp[i-1][w-h[i-1].Size]+h[i-1].Value)
			}
		}
	}
	w := capacity
	for i := n; i > 0; i-- {
		if dp[i][w] != dp[i-1][w] {
			selectedPresents = append(selectedPresents, h[i-1])
			w -= h[i-1].Size
		}
	}
	return selectedPresents
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
