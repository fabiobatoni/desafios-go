package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type LogRequest struct {
	method       string
	endpoint     string
	responseTime int
	timestamp    time.Time
}

// Função principal
func main() {
	rand.Seed(time.Now().UnixNano())

	requests := []string{
		"GET /users",
		"POST /users",
		"GET /invalid",
		"DELETE /users/1",
	}

	for i, request := range requests {
		// 1. Separar método e endpoint
		method, endpoint := parseRequest(request)

		// 2. Gerar status e tempo de resposta
		status := generateStatus()
		responseTime := generateResponseTime()

		// 3. Criar struct com os dados
		log := LogRequest{
			method:       method,
			endpoint:     endpoint,
			responseTime: responseTime,
			timestamp:    time.Now(),
		}

		// 4. Chamar função para imprimir o log
		logResponse(i+1, log, status)
	}
}

// ---------------------- Funções auxiliares ----------------------

func parseRequest(request string) (string, string) {
	parts := strings.Split(request, " ")
	return parts[0], parts[1]
}

func generateStatus() int {
	statusOptions := []int{200, 201, 404, 500}
	randomIndex := rand.Intn(len(statusOptions))
	return statusOptions[randomIndex]
}

func generateResponseTime() int {
	return rand.Intn(451) + 50 
}

func getStatusMessage(status int) string {
	messages := map[int]string{
		200: "OK",
		201: "Created",
		404: "Not Found",
		500: "Error",
	}
	return messages[status]
}

func logResponse(index int, log LogRequest, status int) {
	fmt.Printf("[%d] %s %s - %d %s (%dms) | %s\n",
		index,
		log.method,
		log.endpoint,
		status,
		getStatusMessage(status),
		log.responseTime,
		log.timestamp.Format("15:04:05"),
	)
}
