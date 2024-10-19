package passagesService

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sharedPass/graphs"
)

func GetAllRoutes(origin string, destination string) ([]string, error) {

	getOtherFlights()

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

	respA, err := http.Get("http://localhost:8080/passages/flights") // Fazendo uma requisição ao servidor A
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}
	defer respA.Body.Close()

	var flightsA interface{}
	if err := json.NewDecoder(respA.Body).Decode(&flightsA); err != nil {
		fmt.Println("Erro ao decodificar resposta:", err)
		return
	}

	respB, err := http.Get("http://localhost:8082/passages/flights") // Fazendo uma requisição ao servidor C
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

	fmt.Println(flightsA)
	fmt.Print("\n\n")
	fmt.Println(flightsB)
}
