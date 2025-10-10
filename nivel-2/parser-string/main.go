package main

import (
	"fmt"
	"strings"
)

func main() {
	text := "name=John&age=30&city=NYC"

	parts := strings.Split(text, "&")

	for _, result := range parts {
		fmt.Println(strings.Replace(result, "=", " : ", 1))
	}
}
