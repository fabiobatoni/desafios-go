package main

import "fmt"

func main() {

	methods := []string{"GET", "POST", "GET", "DELETE", "GET", "POST", "PUT"}

	countMethods := make(map[string]int)

	for _ , method := range methods {
		countMethods[method] ++ 
	} 

	fmt.Println("Contagem de Metodos:", countMethods)

	fmt.Println("\nMetodos repeditos:")

	for method , count := range countMethods {
		if count > 1 {
			fmt.Printf("- %s: %d vezes\n", method, count)
		}
	}

}

