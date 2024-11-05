package passagesService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sharedPass/collections"
	"sharedPass/graphs"
	"sharedPass/queues"
	"sharedPass/vectorClock"
	"sort"
	"sync"
)

func GetAllRoutes(origin string, destination string) ([][]graphs.Route, error) {

	filghtsOne, flightsTwo := getOtherFlights() // Pegando vôos dos outros servers

	routesMerged := mergeRoutes(graphs.Graph, filghtsOne, flightsTwo) // Gerando um grafo completo

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
func SolicitationLocal(routes []graphs.Route, externalServerId int, externalClock vectorClock.VectorClock, isBuy bool) (bool, error) {
	var solicitation queues.Solicitation
	// Pega as rotas e o ID do servidor e passa para dentro da struct de solitação
	solicitation.Clock = externalClock
	solicitation.ServerID = externalServerId
	solicitation.IsBuy = isBuy
	solicitation.Routes = routes
	// Criar o canal de resposta para receber a resposta de efetuação de compra
	solicitation.ResponseCh = make(chan bool)
	// Encaminha a solitação através do canal de comunicação queue
	queues.SolicitationsQueue <- &solicitation
	// Recebe a resposta de efetuação de compra
	confirmation := <-solicitation.ResponseCh

	return confirmation, nil
}

func sendBuyRequest(
	routes []graphs.Route, serverId int,
	clock vectorClock.VectorClock,
	port string, wg *sync.WaitGroup,
	channelBuy chan *queues.RequestBuy,
	response chan bool) {
	defer wg.Done()

	data := collections.Body{
		Routes:   routes,
		Clock:    &clock,
		ServerId: &serverId,
	}
	jsonRoutes, err := json.Marshal(data)
	if err != nil {
		response <- false
		return
	}

	request := queues.RequestBuy{
		DataJson:   jsonRoutes,
		Port:       port,
		ResponseCh: make(chan bool),
	}
	channelBuy <- &request
	confirmation := <-request.ResponseCh

	response <- confirmation
}

func sendRequestRollBack(port string, jsonRoutes []byte) bool {
	resp, err := http.Post("http://localhost:"+port+"/passages/rollback", "application/json", bytes.NewBuffer(jsonRoutes))
	if err != nil {
		fmt.Println("Erro:", err)
		return false
	}
	defer resp.Body.Close()
	return true
}

func Buy(routes []graphs.Route, externalServerId int, externalClock vectorClock.VectorClock) (bool, error) {
	var wg sync.WaitGroup
	routesCompanylocal := filterByCompany(routes, os.Getenv("LOCAL_COMPANY"))
	routesCompanyOne := filterByCompany(routes, os.Getenv("ONE_COMPANY"))
	routesCompanyTwo := filterByCompany(routes, os.Getenv("TWO_COMPANY"))

	// Compra local
	responseLocal := false
	var err error = nil
	if routesCompanylocal != nil {
		responseLocal, err = SolicitationLocal(routesCompanylocal, externalServerId, externalClock, true)
		if err != nil {
			fmt.Println("Erro ao comprar passagem local:", err)
			return false, err
		}
	}

	// Atualizando o relógio vetorial
	vectorClock.LocalClock.Update(externalClock)

	if responseLocal {

		responseOne := make(chan bool, 1)
		responseTwo := make(chan bool, 1)
		if routesCompanyOne != nil {
			wg.Add(1)
			go sendBuyRequest(routesCompanyOne, externalServerId, externalClock, os.Getenv("ONE_PORT"), &wg, queues.RequestQueueOne, responseOne)
		} else {
			responseOne <- true
		}
		if routesCompanyTwo != nil {
			wg.Add(1)
			go sendBuyRequest(routesCompanyTwo, externalServerId, externalClock, os.Getenv("TWO_PORT"), &wg, queues.RequestQueueTwo, responseTwo)
		} else {
			responseTwo <- true
		}
		wg.Wait()

		// Verifica as respostas
		responseChOne := <-responseOne
		responseChTwo := <-responseTwo

		if (responseChOne && !responseChTwo) || (!responseChOne && responseChTwo) {
			clock := vectorClock.LocalClock.Copy()
			data := collections.Body{
				Routes:   nil,
				Clock:    &clock,
				ServerId: &vectorClock.ServerId,
			}
			SolicitationLocal(routesCompanylocal, externalServerId, externalClock, false)
			if !responseChOne {
				data.Routes = routesCompanyTwo
				jsonRoutesTwo, _ := json.Marshal(data)
				fmt.Println("Rollback na empresa 2")
				sendRequestRollBack(os.Getenv("TWO_PORT"), jsonRoutesTwo)
				} else if !responseChTwo {
					data.Routes = routesCompanyOne
					jsonRoutesOne, _ := json.Marshal(data)
					fmt.Println("Rollback na empresa 1")
				sendRequestRollBack(os.Getenv("ONE_PORT"), jsonRoutesOne)
			}
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

func RollBack(routes []graphs.Route, externalServerId int, externalClock vectorClock.VectorClock) (bool, error) {
	routesCompanylocal := filterByCompany(routes, os.Getenv("LOCAL_COMPANY"))

	// Atualizando o relógio vetorial
	vectorClock.LocalClock.Update(externalClock)

	if routesCompanylocal != nil {
		_, err := SolicitationLocal(routesCompanylocal, externalServerId, externalClock, false)
		if err != nil {
			fmt.Println("Erro ao comprar passagem local:", err)
			return false, err
		}
	}
	return true, nil
}

// Retorna todas os vôos
func GetAllFlights() (map[string][]graphs.Route, error) {
	allFlights := graphs.Graph

	return allFlights, nil
}

// Pega os vôos dos outros servers
func getOtherFlights() (map[string][]graphs.Route, map[string][]graphs.Route) {
	var flightsOne map[string][]graphs.Route
	var flightsTwo map[string][]graphs.Route

	respOne, _ := http.Get("http://localhost:" + os.Getenv("ONE_PORT") + "/passages/flights") // Fazendo uma requisição ao servidor B
	if respOne != nil {
		defer respOne.Body.Close()
		if err := json.NewDecoder(respOne.Body).Decode(&flightsOne); err != nil {
			fmt.Println("Erro ao decodificar resposta:", err)
		}
	}

	respTwo, _ := http.Get("http://localhost:" + os.Getenv("TWO_PORT") + "/passages/flights") // Fazendo uma requisição ao servidor C
	if respTwo != nil {
		defer respTwo.Body.Close()
		if err := json.NewDecoder(respTwo.Body).Decode(&flightsTwo); err != nil {
			fmt.Println("Erro ao decodificar resposta:", err)
		}
	}

	return flightsOne, flightsTwo
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
