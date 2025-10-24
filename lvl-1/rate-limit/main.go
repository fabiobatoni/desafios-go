package main

import "fmt"

const totalAllowed = 1000

func main() {
	var requests = 1200

	if requests <= totalAllowed {
		requestsRemaining := requests - totalAllowed
		percentagem := (float64(requestsRemaining) / float64(totalAllowed)) * 100
		fmt.Println("Requisições feitas: ", requests)
		fmt.Println("Restantes: ", requestsRemaining*-1)
		fmt.Printf("Porcentagem: %.2f%%\n", percentagem*-1)
	} else {
		moreRequests := totalAllowed - requests
		percentagem := (float64(requests) / float64(totalAllowed)) * 100
		fmt.Println("Requisições feitas: ", requests)
		fmt.Println("Requisições a mais: ", moreRequests*-1)
		fmt.Printf("Porcentagem: %.2f%%\n", percentagem)
	}
}
