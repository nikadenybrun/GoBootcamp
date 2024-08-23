package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

type Anscombe struct {
	Mean   float64
	Median float64
	Mode   int
	SD     float64
}

func main() {
	mean := flag.Int("mean", 0, "Mean")
	median := flag.Int("median", 0, "Median")
	mode := flag.Int("mode", 0, "Mode")
	sd := flag.Int("sd", 0, "Sd")
	flag.Parse()
	any_flag := *mean + *median + *mode + *sd
	if any_flag == 0 {
		*mean = 1
		*median = 1
		*mode = 1
		*sd = 1
	}
	scanner := bufio.NewScanner(os.Stdin)
	rs := make([]int, 0)
	for scanner.Scan() {
		res := scanner.Text()
		if res == "" {
			fmt.Print("you entered a space")
			break
		} else if res != "q" {
			val, err := strconv.Atoi(res)
			if err == nil && val > -100000 && val < 100000 {
				rs = append(rs, val)
			} else if err == nil && (val < -100000 || val > 100000) {
				fmt.Print("you entered too big/small integer")
				break
			} else {
				fmt.Print("you entered something wierd")
				break
			}
		}
		if res == "q" && len(rs) > 0 {
			a := Anscombe{}
			a.Mean = find_mean(rs)
			a.Median = find_median(rs)
			a.Mode = find_mode(rs)
			a.SD = find_stdDev(rs)
			if *mean == 1 {
				fmt.Printf("Mean: %.2f\n", a.Mean)
			}
			if *median == 1 {
				fmt.Printf("Median: %.2f\n", a.Median)
			}
			if *mode == 1 {
				fmt.Println("Mode:", a.Mode)
			}
			if *sd == 1 {
				fmt.Printf("SD: %.2f\n", a.SD)
			}
			break
		} else if res == "q" && len(rs) == 0 {
			fmt.Print("you entered something wierd")
			break
		}
	}

}

func sum(numbers []int) float64 {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return float64(sum)
}

func find_mean(rs []int) float64 {
	mean := float64(sum(rs) / float64(len(rs)))
	return mean
}

func find_median(rs []int) float64 {
	sort.Ints(rs)
	if len(rs)%2 == 1 {
		return float64(rs[len(rs)/2])
	} else {
		return float64((float64(rs[len(rs)/2] + rs[len(rs)/2-1])) / 2)
	}
}

func find_mode(rs []int) int {
	m := make(map[int]int)
	for _, elem := range rs {
		if value, inMap := m[elem]; inMap {
			m[elem] = value + 1
		} else {
			m[elem] = 1
		}
	}
	maxx_count := m[0]
	maxx_i := 0
	for i, elem := range m {
		if elem > maxx_count {
			maxx_count = elem
			maxx_i = i
		}
	}
	return maxx_i
}

func find_stdDev(numbers []int) float64 {
	mean := find_mean(numbers)
	var sd float64
	for _, num := range numbers {
		sd += math.Pow(float64(num)-mean, 2)
	}
	return math.Sqrt(sd / float64(len(numbers)))
}
