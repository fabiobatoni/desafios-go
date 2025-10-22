package main

import (
	"fmt"
	"issue/github"
	"log"
	"os"
	"time"
)

// Issues exibe tablea de problemas do github que correspondem aos termos de pesquisa

func main() {

	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)

	lessThanMonth, lessThanYear, moreThanYear := categorizeIssues(result.Items)


	fmt.Printf("\nIssues com menos de 1 mês:\n")
	for _, item := range lessThanMonth {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}

	fmt.Printf("\nIssues com menos de 1 ano:\n")
	for _, item := range lessThanYear {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}

	fmt.Printf("\nIssues com mais de 1 ano:\n")
	for _, item := range moreThanYear {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}

}

func categorizeIssues(issues []*github.Issue) (month, year, older []*github.Issue) {
	for _, item := range issues {
		age := time.Since(item.CreatedAt)
		switch {
		case age < 30*24*time.Hour:
			month = append(month, item)
		case age < 365*24*time.Hour:
			year = append(year, item)
		default:
			older = append(older, item)
		}
	}
	return
}
