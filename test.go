package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// Estruturas para a requisição e resposta
type Route struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Seats   int    `json:"seats"`
	Company string `json:"company"`
}

type ResponseBuy struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func main() {
	// Configuração dos servidores
	serverURLs := []string{
		"https://localhost:8080/passages/buy",
		"https://localhost:8081/passages/buy",
		"https://localhost:8082/passages/buy",
	}

	// Exibe as rotas disponíveis (substitua isso pela lógica real de obtenção das rotas)
	routes := []Route{
		{From: "RECIFE", To: "SALVADOR", Seats: 1, Company: "A"},
		// Adicione mais rotas conforme necessário
	}

	var wg sync.WaitGroup
	numGoroutines := 20 // Número de goroutines para simular compras concorrentes
	var mu sync.Mutex
	var maxDuration time.Duration

	// Inicia as goroutines para compras concorrentes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ { // Número de compras por goroutine
				for _, url := range serverURLs {
					// Simula a compra de uma passagem
					route := routes[j%len(routes)] // Seleciona uma rota
					requestBody, err := json.Marshal(route)
					if err != nil {
						fmt.Printf("Erro ao criar corpo da requisição: %v\n", err)
						continue
					}

					start := time.Now()
					resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
					if err != nil {
						fmt.Printf("Erro ao fazer a requisição: %v\n", err)
						continue
					}
					defer resp.Body.Close()

					duration := time.Since(start)
					mu.Lock()
					if duration > maxDuration {
						maxDuration = duration
					}
					mu.Unlock()

					responseBody, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Printf("Erro ao ler a resposta: %v\n", err)
						continue
					}

					var responseBuy ResponseBuy
					if err := json.Unmarshal(responseBody, &responseBuy); err != nil {
						fmt.Printf("Erro ao decodificar a resposta: %v\n", err)
						continue
					}
					fmt.Print(responseBuy.Status, " ")
				}
			}
		}()
	}

	wg.Wait() // Espera todas as goroutines terminarem
	fmt.Printf("Todas as compras foram processadas. Maior tempo de resposta: %v\n", maxDuration)
}
