package main

import "fmt"

func main() {
	fmt.Println(HttpStatusValidator(200))
	fmt.Println(HttpStatusValidator(404))
	fmt.Println(HttpStatusValidator(500))
	fmt.Println(HttpStatusValidator(301))
	fmt.Println(HttpStatusValidator(999))
}

func HttpStatusValidator(s int) string {
	switch {
	case s <= 299:
		return "Success"
	case s <= 399:
		return "Redirect"
	case s <= 499:
		return "Client Error"
	case s <= 599:
		return "Server Error"
	default:
		return "Invalid"
	}
}
