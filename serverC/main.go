package main

import (
	"sharedPass/graphs"
	"sharedPass/passages/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Cria um novo roteador do gin
	router := gin.Default()
	routes.RegisterRoutes(router)

	// Roda o server
	graphs.ReadRoutes() // Carregando o gráfico na memória
	router.Run(":8082") // Roda na porta 8081
}
