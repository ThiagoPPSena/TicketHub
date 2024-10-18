// server_a.go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func requestToServerB() {
	// Faz a requisição HTTP para o Servidor B
	resp, err := http.Get("http://localhost:8080/hello/22")
	if err != nil {
		log.Fatalf("Erro ao fazer requisição ao Servidor B: %v", err)
	}
	defer resp.Body.Close()

	// Lê a resposta do Servidor B
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erro ao ler a resposta do Servidor B: %v", err)
	}

	fmt.Printf("Resposta do Servidor B: %s\n", string(body))
}

func main() {
	requestToServerB()
}
