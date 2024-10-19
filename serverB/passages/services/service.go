package passagesService

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sharedPass/graphs"
	"sort"
)

func GetAllRoutes(origin string, destination string) ([][]string, error) {

	filghtsB, flightsC, err := getOtherFlights() // Pegando vôos dos outros servers
	if err != nil {
		return nil, err
	}

	routesMerged := mergeRoutes(graphs.Graph, filghtsB, flightsC) // Gerando um grafo completo

	visited := make(map[string]bool) // Lista para mapear se um nó do grafo já foi visitado
	var path []string                // Salva uma rota
	var allPaths [][]string          // Salva todas as rotas possíveis

	// Método para saber todas as rotas possíveis
	graphs.FindRoutes(routesMerged, origin, destination, visited, path, &allPaths)

	// Ordenando da menor rota para a maior
	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	// Pegando as 10 menores rotas disponíveis
	if len(allPaths) >= 10 {
		allPaths = allPaths[:10]
	}

	return allPaths, nil
}

// Ainda implementar
func Buy(routes []string) (bool, error) {
	return false, nil
}

// Retorna todas os vôos
func GetAllFlights() (map[string][]graphs.Route, error) {
	allFlights := graphs.Graph

	return allFlights, nil
}

// Pega os vôos dos outros servers
func getOtherFlights() (map[string][]graphs.Route, map[string][]graphs.Route, error) {

	respA, err := http.Get("http://localhost:8080/passages/flights") // Fazendo uma requisição ao servidor B
	if err != nil {
		fmt.Println("Erro:", err)
		return nil, nil, err
	}
	defer respA.Body.Close()

	var flightsA map[string][]graphs.Route
	if err := json.NewDecoder(respA.Body).Decode(&flightsA); err != nil {
		fmt.Println("Erro ao decodificar resposta:", err)
		return nil, nil, err
	}

	respC, err := http.Get("http://localhost:8082/passages/flights") // Fazendo uma requisição ao servidor C
	if err != nil {
		fmt.Println("Erro:", err)
		return nil, nil, err
	}
	defer respC.Body.Close()

	var flightsC map[string][]graphs.Route
	if err := json.NewDecoder(respC.Body).Decode(&flightsC); err != nil {
		fmt.Println("Erro ao decodificar resposta:", err)
		return nil, nil, err
	}

	return flightsA, flightsC, nil
}

// Junta todos os vôos em um único grafo
func mergeRoutes(routeMaps ...map[string][]graphs.Route) map[string][]graphs.Route {
	finalGraph := make(map[string][]graphs.Route)

	for _, routeMap := range routeMaps {
		for city, routes := range routeMap {
			finalGraph[city] = append(finalGraph[city], routes...)
		}
	}

	return finalGraph
}
