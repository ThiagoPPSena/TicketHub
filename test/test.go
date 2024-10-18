package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	//Ler parametro id da URL hello/id
	id := r.URL.Path[len("/hello/"):]
	//Escrever resposta
	fmt.Fprintf(w, "Hello, World! %s", id)
}

func main() {
	http.HandleFunc("/hello/", helloHandler)

	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
