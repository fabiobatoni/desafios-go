package main

import "fmt"

func main() {
    // Teste 1: Com parâmetros
    params1 := map[string]string{
        "id": "123",
        "active": "true",
    }
    fmt.Println(BuildURL("https://api.example.com", "/users", params1))
    
    // Teste 2: Sem parâmetros
    params2 := map[string]string{}
    fmt.Println(BuildURL("https://api.example.com", "/users", params2))
    
    // Teste 3: Um só parâmetro
    params3 := map[string]string{"page": "1"}
    fmt.Println(BuildURL("https://api.example.com", "/posts", params3))
}

func BuildURL(base string, endpoint string, params map[string]string) string {

	primeiro := true

    url := base + endpoint
    
    if len(params) == 0 {
        return url
    }
    
    for key, value := range params {
		if primeiro {
			url += "?"
			primeiro = false
		} else {
			url += "&"
		}
		url += key + "=" + value
	}
    
    return url
}