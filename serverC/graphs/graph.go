package graphs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Route struct {
	From    string
	To      string
	Seats   int
	Company string
}

// Mapa de rotas
var Graph map[string][]Route

// Função para ler as rotas do arquivo JSON
func ReadRoutes() {
	// Abre o arquivo JSON
	file, err := os.Open("./files/routes.json")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	// Decodifica o arquivo JSON
	err = json.NewDecoder(file).Decode(&Graph)
	if err != nil {
		fmt.Println("Erro ao decodificar JSON:", err)
		return
	}
}

// Função auxiliar para verificar a disponibilidade de assentos
func checkAvailability(routesToCheck []Route) bool {
	for _, route := range routesToCheck {
		// Verifica se a rota existe e tem assentos disponíveis
		if routes, ok := Graph[route.From]; ok {
			found := false
			for _, r := range routes {
				if r.To == route.To && r.Seats > 0 {
					fmt.Println(r.Seats)
					found = true
					break
				}
			}
			if !found {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

// Função para decrementar o número de assentos
func buySeat(origin, destination string) {
	if routes, ok := Graph[origin]; ok {
		for i, route := range routes {
			if route.To == destination {
				if route.Seats > 0 {
					Graph[origin][i].Seats--
				}
			}
		}
	}
}

// Função para tentar comprar assentos para uma lista de rotas
func BuySeats(routesToBuy []Route) bool {
	// Primeiro verifica a disponibilidade de todas as rotas
	if !checkAvailability(routesToBuy) {
		return false
	}

	// Se todos os assentos estiverem disponíveis, realiza a compra
	for _, routeToBuy := range routesToBuy {
		buySeat(routeToBuy.From, routeToBuy.To)
	}
	return true
}

// Função para restaurar o número de assentos em caso de rollback
func RollBack(routesToRollback []Route) (bool) {
	for _, route := range routesToRollback {
		if routes, ok := Graph[route.From]; ok {
			for i, r := range routes {
				if r.To == route.To {
					Graph[route.From][i].Seats++ // Incrementa os assentos de volta
				}
			}
		}
	}

	return true
}

// Função para salvar o grafo atualizado no arquivo JSON
func SaveSeats() {
	file, err := os.OpenFile("./files/routes.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(Graph); err != nil {
		fmt.Println("Erro ao salvar o arquivo JSON:", err)
	}
}

// Método para encontrar todas as rotas disponíveis dada uma origem e destino
func FindRoutes(graph map[string][]Route, origin string, destination string, visited map[string]bool, path []Route, allpaths *[][]Route) {
	visited[origin] = true

	// Se a origem for igual ao destino, adiciona a rota encontrada
	if origin == destination {
		tempPath := make([]Route, len(path)) // Faz uma cópia do caminho
		copy(tempPath, path)
		*allpaths = append(*allpaths, tempPath)
	} else {
		// Verifica vizinhos (rotas possíveis) e faz a busca recursiva
		for _, neighbor := range graph[origin] {
			if neighbor.Seats > 0 && !visited[neighbor.To] {
				// Adiciona a rota ao caminho atual
				newPath := append(path, neighbor)
				FindRoutes(graph, neighbor.To, destination, visited, newPath, allpaths)
			}
		}
	}

	// Marca como não visitado (permite outras rotas usarem essa cidade novamente)
	visited[origin] = false
}
