package main

import (
	"sharedPass/passages/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Cria um novo roteador do gin
	router := gin.Default()
	routes.RegisterRoutes(router)

	// Rota o
	router.Run(":8080") // Roda na porta 8080
}
