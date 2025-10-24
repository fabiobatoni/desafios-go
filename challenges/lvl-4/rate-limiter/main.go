package main

import "fmt"

func main() {
	requests := []int{1,2,3,4,5,6,7,8}
	RateLimiter(requests)
}

func RateLimiter(requests []int) {
    limite := 5
    contador := 0

	for _, req := range requests {
		if contador < limite {
			fmt.Printf("Request %d %s\n", req, "- Aceita")
			contador ++
		} else {
			fmt.Printf("Request %d %s\n", req, "- Recusada")
		}
	}
}