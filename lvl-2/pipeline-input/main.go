package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {

	testCases := []string{
		"  Hello World! @#$  ",
		"TESTE123!@#",
	}

	for _, s := range testCases {
		fmt.Println(pipeInputValidator(s))
	}
}

func pipeInputValidator(s string) string {
	reg, _ := regexp.Compile(`[^a-zA-Z0-9]+`)

	if strings.HasPrefix(s, " ") || strings.HasSuffix(s, " ") {
		s = strings.Replace(s, " ", "", 1)
	}
	s = reg.ReplaceAllString(s, "")
	s = strings.ToLower(s)

	if len(s) >= 50 {
		s = s[:50]
	}

	return s
}
