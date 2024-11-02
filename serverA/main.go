package main

import (
	"io"
	"sharedPass/graphs"
	"sharedPass/passages/routes"
	"sharedPass/queues"
	"sharedPass/vectorClock"

	"github.com/gin-gonic/gin"
)

func main() {
	// Define a saída de logs para descartar
	gin.DefaultWriter = io.Discard
	// Seta o modo de execução do gin para release
	gin.SetMode(gin.ReleaseMode)

	// Cria um novo roteador do gin
	router := gin.Default()
	routes.RegisterRoutes(router)

	// Iniciar o relogio vetorial local
	vectorClock.NewVectorClock(3)
	// Iniciar a queue de solicitações
	queues.StartProcessQueue()
	queues.StartProcessRequestQueue()
	// Seta o id do servidor
	vectorClock.SetServerId(0)

	// Roda o server
	graphs.ReadRoutes() // Carregando o gráfico na memória
	router.Run(":8080") // Roda na porta 8080
}
