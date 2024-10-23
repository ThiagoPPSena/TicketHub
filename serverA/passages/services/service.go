package passagesService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sharedPass/graphs"
	"sort"
)

func GetAllRoutes(origin string, destination string) ([][]graphs.Route, error) {

	filghtsB, flightsC, err := getOtherFlights() // Pegando vôos dos outros servers
	if err != nil {
		return nil, err
	}

	routesMerged := mergeRoutes(graphs.Graph, filghtsB, flightsC) // Gerando um grafo completo

	visited := make(map[string]bool) // Lista para mapear se um nó do grafo já foi visitado
	var path []graphs.Route          // Salva uma rota
	var allPaths [][]graphs.Route    // Salva todas as rotas possíveis

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
func Buy(routes []graphs.Route) (bool, error) {
	//routesCompanyA := filterByCompany(routes, "A")
	routesCompanyB := filterByCompany(routes, "B")
	routesCompanyC := filterByCompany(routes, "C")

	// Converte a estrutura para JSON
	jsonRoutesB, err := json.Marshal(routesCompanyB)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return false, err
	}

	// Converte a estrutura para JSON
	jsonRoutesC, err := json.Marshal(routesCompanyC)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return false, err
	}

	if routesCompanyB != nil {
		respB, err := http.Post("http://localhost:8081/passages/buy", "application/json", bytes.NewBuffer(jsonRoutesB)) // Fazendo uma requisição ao servidor B
		if err != nil {
			fmt.Println("Erro:", err)
			return false, err
		}
		defer respB.Body.Close()
	}

	if routesCompanyC != nil {
		respC, err := http.Post("http://localhost:8082/passages/buy", "application/json", bytes.NewBuffer(jsonRoutesC)) // Fazendo uma requisição ao servidor C
		if err != nil {
			fmt.Println("Erro:", err)
			return false, nil
		}
		defer respC.Body.Close()
	}

	// Efetuar atualização de compra da companhia A aqui
	return true, nil

}

// Retorna todas os vôos
func GetAllFlights() (map[string][]graphs.Route, error) {
	allFlights := graphs.Graph

	return allFlights, nil
}

// Pega os vôos dos outros servers
func getOtherFlights() (map[string][]graphs.Route, map[string][]graphs.Route, error) {

	respB, err := http.Get("http://localhost:8081/passages/flights") // Fazendo uma requisição ao servidor B
	if err != nil {
		fmt.Println("Erro:", err)
		return nil, nil, err
	}
	defer respB.Body.Close()

	var flightsB map[string][]graphs.Route
	if err := json.NewDecoder(respB.Body).Decode(&flightsB); err != nil {
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

	return flightsB, flightsC, nil
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

// Função que filtra as rotas pelo nome da empresa
func filterByCompany(routes []graphs.Route, company string) []graphs.Route {
	var filteredRoutes []graphs.Route
	for _, route := range routes {
		if route.Company == company {
			filteredRoutes = append(filteredRoutes, route)
		}
	}
	return filteredRoutes
}
