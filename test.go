package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Estruturas para a requisição e resposta
type Route struct {
	From    string `json:"From"`
	To      string `json:"To"`
	Seats   int    `json:"Seats"`
	Company string `json:"Company"`
}

type requestBuy struct {
	Routes []Route `json:"Routes"`
}

type ResponseBuy struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func craftGetURL(baseURL string, origem string, destino string) string {
	// Configuração dos parâmetros de consulta
	params := url.Values{}
	params.Add("origem", origem)
	params.Add("destino", destino)

	// Construção da URL completa com os parâmetros de consulta
	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

func getRoutes(url string) ([][]Route, error) {
	var routes [][]Route
	resp, err := http.Get(url)
	if err != nil {
		return routes, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&routes)
	if err != nil {
		return routes, err
	}
	return routes, nil
}

func main() {
	urlGet := craftGetURL("http://localhost:8080/passages/routes", "ARACAJU", "PORTO VELHO") // ARACAJU PORTO VELHO // BELO HORIZONTE SALVADOR
	serverURLs := []string{
		"http://localhost:8080/passages/buy",
		"http://localhost:8081/passages/buy",
		"http://localhost:8082/passages/buy",
	}

	// Buscar as rotas via http e pegar a priemira rota
	routes, err := getRoutes(urlGet)
	if err != nil {
		fmt.Println("Erro ao buscar as rotas:", err)
		return
	}
	fmt.Println("Rotas disponíveis:", routes[0])
	// Construir request do buy para enviar no post com a primeira rota
	requestBuy := requestBuy{
		Routes: routes[0],
	}
	// Serializar a rota para JSON
	jsonRoute, err := json.Marshal(requestBuy)
	if err != nil {
		fmt.Println("Erro ao serializar a rota:", err)
		return
	}

	var wg sync.WaitGroup
	numGoroutines := 50 // Número de goroutines para simular compras concorrentes
	numTry := 50       // Número de tentativas de compra por goroutine
	var mu sync.Mutex
	var maxDuration time.Duration

	// Cada servidor vai 20 de cada goroutine
	for _, urlBuy := range serverURLs {
		wg.Add(1)
		go func(urlBuy string) {
			defer wg.Done()
			for i := 0; i < numGoroutines; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < numTry; j++ { // Número de compras por goroutine
						start := time.Now()
						resp, err := http.Post(urlBuy, "application/json", bytes.NewReader(jsonRoute))
						if err != nil {
							fmt.Println("Erro ao comprar passagem:", err)
							return
						}
						defer resp.Body.Close()
						fmt.Print(resp.StatusCode, " ")
						var responseBuy bool
						err = json.NewDecoder(resp.Body).Decode(&responseBuy)
						if err != nil {
							fmt.Println("Erro ao desserializar a resposta:", err)
							return
						}
						// Atualiza o maior tempo de resposta
						mu.Lock()
						if time.Since(start) > maxDuration {
							maxDuration = time.Since(start)
						}
						mu.Unlock()
					}
				}()
			}
		}(urlBuy)
	}

	wg.Wait() // Espera todas as goroutines terminarem
	fmt.Printf("\nTodas as compras foram processadas. Maior tempo de resposta: %v\n", maxDuration)
}
