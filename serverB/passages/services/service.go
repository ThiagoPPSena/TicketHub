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

	respC, err := http.Get("http://localhost:8082/passages/flights") // Fazendo uma requisição ao servidor C
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

	fmt.Println(flightsA)
	fmt.Print("\n\n")
	fmt.Println(flightsC)
}
