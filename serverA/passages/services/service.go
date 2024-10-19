package passagesService

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sharedPass/graphs"
)

func GetAllRoutes(origin string, destination string) ([]string, error) {

	fmt.Println(graphs.Graph)

	return []string{"Salvador", "Sao Paulo", "Recife"}, nil
}

func Buy(routes []string) (bool, error) {
	return false, nil
}

func GetAllFlights() (map[string][]graphs.Route, error) {
	allFlights := graphs.Graph

	return allFlights, nil
}

func getOtherFlights() {

	respB, err := http.Get("http://localhost:8081/passages/flights") // Fazendo uma requisição ao servidor B
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}
	defer respB.Body.Close()

	var flightsB interface{}
	if err := json.NewDecoder(respB.Body).Decode(&flightsB); err != nil {
		fmt.Println("Erro ao decodificar resposta:", err)
		return
	}

	respC, err := http.Get("http://localhost:8081/passages/flights") // Fazendo uma requisição ao servidor C
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}
	defer respC.Body.Close()

	var flightsC interface{}
	if err := json.NewDecoder(respC.Body).Decode(&flightsC); err != nil {
		fmt.Println("Erro ao decodificar resposta:", err)
		return
	}

	fmt.Println(flightsB)
	fmt.Print("\n\n")
	fmt.Println(flightsC)
}
