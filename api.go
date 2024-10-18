package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Item representa um item básico com ID e Nome
type Item struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

// Lista de itens (como um banco de dados em memória)
var items []Item
var idCounter int = 1

// Função que retorna todos os itens (GET)
func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// Função que adiciona um novo item (POST)
func addItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newItem.ID = idCounter
	idCounter++
	items = append(items, newItem)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newItem)
}

// Função que atualiza um item existente (PUT)
func updateItem(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/items/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var updatedItem Item
	err = json.NewDecoder(r.Body).Decode(&updatedItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, item := range items {
		if item.ID == id {
			items[i].Nome = updatedItem.Nome
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(items[i])
			return
		}
	}

	http.Error(w, "Item não encontrado", http.StatusNotFound)
}

// Função que deleta um item (DELETE)
func deleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/items/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	for i, item := range items {
		if item.ID == id {
			items = append(items[:i], items[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Item não encontrado", http.StatusNotFound)
}

// Função principal
func main() {
	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getItems(w, r)
		} else if r.Method == http.MethodPost {
			addItem(w, r)
		} else {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			updateItem(w, r)
		} else if r.Method == http.MethodDelete {
			deleteItem(w, r)
		} else {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
