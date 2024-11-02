package passagesService

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sharedPass/collections"
	"sharedPass/graphs"
	"sharedPass/queues"
	"sharedPass/vectorClock"
	"sort"
	"sync"
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

// Função para comprar passagem no proprio servidor
func BuyLocal(routes []graphs.Route, externalServerId int, externalClock vectorClock.VectorClock) (bool, error) {
	var solicitation queues.Solicitation
	// Pega as rotas e o ID do servidor e passa para dentro da struct de solitação
	solicitation.Clock = externalClock
	solicitation.ServerID = externalServerId
	solicitation.Routes = routes
	// Criar o canal de resposta para receber a resposta de efetuação de compra
	solicitation.ResponseCh = make(chan bool)
	// Encaminha a solitação através do canal de comunicação queue
	queues.SolicitationsQueue <- &solicitation
	// Recebe a resposta de efetuação de compra
	confirmation := <-solicitation.ResponseCh

	return confirmation, nil
}

func sendBuyRequest(routes []graphs.Route, serverId int, clock vectorClock.VectorClock, port string, wg *sync.WaitGroup, channelBuy chan *queues.RequestBuy) (bool) {
	defer wg.Done()

	if routes == nil {
		return false
	}

	data := collections.Body{
		Routes:   routes,
		Clock:    &clock,
		ServerId: &serverId,
	}
	jsonRoutes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return false
	}

	request := queues.RequestBuy{
		DataJson:   jsonRoutes,
		Port:       port,
		ResponseCh: make(chan bool),
	}
	channelBuy <- &request
	confirmation := <-request.ResponseCh

	return confirmation
}

// Ainda implementar
func Buy(routes []graphs.Route, externalServerId int, externalClock vectorClock.VectorClock) (bool, error) {
	var wg sync.WaitGroup
	routesCompanyA := filterByCompany(routes, "A")
	routesCompanyB := filterByCompany(routes, "B")
	routesCompanyC := filterByCompany(routes, "C")

	//Chama a função que compra passagem local
	if routesCompanyB != nil {
		_, err := BuyLocal(routesCompanyB, externalServerId, externalClock)
		if err != nil {
			fmt.Println("Erro ao comprar passagem local:", err)
			return false, err
		}
	}
	// Atualizando o relógio vetorial
	vectorClock.LocalClock.Update(externalClock)

	wg.Add(2)
	reponseOne := sendBuyRequest(routesCompanyA, externalServerId, externalClock, "8080", &wg, queues.RequestQueueOne)
	responseTwo := sendBuyRequest(routesCompanyC, externalServerId, externalClock, "8082", &wg, queues.RequestQueueTwo)
	wg.Wait()

	fmt.Println("Resposta 1:", reponseOne)
	fmt.Println("Resposta 2:", responseTwo)



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

	respB, err := http.Get("http://localhost:8080/passages/flights") // Fazendo uma requisição ao servidor A
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
