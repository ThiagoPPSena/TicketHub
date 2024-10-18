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
