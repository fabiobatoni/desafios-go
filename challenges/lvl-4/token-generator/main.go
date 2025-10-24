package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println(TokenGenerator())
}

func TokenGenerator() string {

	const caracteres = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	 token := ""

	for i := 0; i < 16; i ++ {
		indiceAleatorio := rand.Intn(len(caracteres))
		caracterEscolhido := caracteres[indiceAleatorio]
		token += string(caracterEscolhido)
	}

	return token
}