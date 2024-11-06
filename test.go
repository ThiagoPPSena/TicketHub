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

	// Buscar as rotas via http e pegar a primeira rota
	routes, err := getRoutes(urlGet)
	if err != nil {
		fmt.Println("Erro ao buscar as rotas:", err)
		return
	}
	// Construir request do buy para enviar no post com a primeira rota
	// requestBuy := requestBuy{
	// 	Routes: routes[0],
	// }
	requestBuyOne := requestBuy{
		Routes: routes[0][1:3],
	}
	requestBuyTwo := requestBuy{
		Routes: routes[0][2:4],
	}
	requestBuyThree := requestBuy{
		Routes: routes[0][3:5],
	}
	fmt.Println("Rotas disponíveis:", requestBuyOne.Routes)
	fmt.Println("Rotas disponíveis:", requestBuyTwo.Routes)
	fmt.Println("Rotas disponíveis:", requestBuyThree.Routes)
	jsonRouteOne, err := json.Marshal(requestBuyOne)
	if err != nil {
		fmt.Println("Erro ao serializar a rota:", err)
		return
	}
	jsonRouteTwo, err := json.Marshal(requestBuyTwo)
	if err != nil {
		fmt.Println("Erro ao serializar a rota:", err)
		return
	}
	jsonRouteThree, err := json.Marshal(requestBuyThree)
	if err != nil {
		fmt.Println("Erro ao serializar a rota:", err)
		return
	}

	allJsonRoutes := [][]byte{jsonRouteOne, jsonRouteTwo, jsonRouteThree}

	var wg sync.WaitGroup
	var startWg sync.WaitGroup // WaitGroup adicional para sincronizar o início das goroutines
	numGoroutines := 10        // Número de goroutines para simular compras concorrentes
	numTry := 15               // Número de tentativas de compra por goroutine
	var mu sync.Mutex
	var maxDuration time.Duration
	count := make([]int, len(serverURLs))

	// Cada servidor vai 20 de cada goroutine
	for serverIndex, urlBuy := range serverURLs {
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			startWg.Add(1)
			go func(urlBuy string) {
				defer wg.Done()
				startWg.Done()                // Indica que a goroutine está pronta para começar
				startWg.Wait()                // Aguarda até que todas as goroutines estejam prontas para começar
				for j := 0; j < numTry; j++ { // Número de compras por goroutine
					start := time.Now()
					count[serverIndex]++
					resp, err := http.Post(urlBuy, "application/json", bytes.NewReader(allJsonRoutes[serverIndex]))
					if err != nil {
						fmt.Println("Erro ao comprar passagem:", err)
						return
					}
					defer resp.Body.Close()
					fmt.Println(resp.StatusCode, "->", serverIndex+1, "->", i+1, "->", j+1)
					var responseBuy ResponseBuy
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
			}(urlBuy)
		}
	}

	startWg.Wait() // Aguarda até que todas as goroutines estejam prontas para começar
	wg.Wait()      // Espera todas as goroutines terminarem
	fmt.Printf("\nTodas as compras foram processadas. Maior tempo de resposta: %v\n", maxDuration)
	fmt.Println("Compras por servidor:", count)
}
