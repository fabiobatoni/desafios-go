package main

import (
	"fmt"
	"strings"
)

func main() {
	result := DetectiveJSON(`{"name":"John","age":30}`)
	result2 := DetectiveJSON(`{name:John}`)
	result3 := DetectiveJSON(`{"name":"John"`)
	fmt.Println("O JSON está: ", result)
	fmt.Println("O JSON está: ", result2)
	fmt.Println("O JSON está: ", result3)
}

func DetectiveJSON(jsonString string) bool {
	if len(jsonString) < 2 {
		return false
	}

	if jsonString[0] != '{' ||  jsonString[len(jsonString)-1] != '}' {
		return false
	}

	if !strings.Contains(jsonString, ":") {
    	return false
	}

	if !strings.Contains(jsonString, ",") {
		return false
	}

	return true
}
