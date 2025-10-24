package main

import (
	"fmt"
	"sort"
)

func main() {
	responseTimes := []float64{120.5, 340.2, 89.7, 450.1, 200.3}
	avg, min, max, fastCount := CalculateTime(responseTimes)

	fmt.Printf("Tempo médio: %.2f ms\n", avg)
	fmt.Printf("Menor tempo: %.2f ms\n", min)
	fmt.Printf("Maior tempo: %.2f ms\n", max)
	fmt.Printf("Requests < 200ms: %d\n", fastCount)
}

func CalculateTime(r []float64) (float64, float64, float64, int) {
	sort.Float64s(r)

	min := r[0]
	max := r[len(r)-1]

	var sum float64
	fastCount := 0

	for _, value := range r {
		sum += value
		if value < 200 {
			fastCount++
		}
	}

	avg := sum / float64(len(r))
	return avg, min, max, fastCount
}
