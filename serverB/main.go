package main

import (
	"sharedPass/graphs"
	"sharedPass/passages/routes"
	"sharedPass/queues"
	"sharedPass/vectorClock"

	"github.com/gin-gonic/gin"
)

func main() {
	// Cria um novo roteador do gin
	router := gin.Default()
	routes.RegisterRoutes(router)

	// Iniciar o relogio vetorial local
	vectorClock.NewVectorClock(3)
	// Iniciar a queue de solicitações
	queues.StartProcessQueue()
	// Seta o id do servidor
	vectorClock.SetServerId(1)

	// Roda o server
	graphs.ReadRoutes() // Carregando o gráfico na memória
	router.Run(":8081") // Roda na porta 8080
}
