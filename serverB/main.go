package main

import (
	"io"
	"sharedPass/graphs"
	"sharedPass/passages/routes"
	"sharedPass/queues"
	"sharedPass/vectorClock"
	"os"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Define a saída de logs para descartar
	gin.DefaultWriter = io.Discard
	// Seta o modo de execução do gin para release
	gin.SetMode(gin.ReleaseMode)

	// Cria um novo roteador do gin
	router := gin.Default()
	routes.RegisterRoutes(router)

	// Carregar a .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	// Iniciar o relogio vetorial local
	vectorClock.NewVectorClock(3)
	// Iniciar a queue de solicitações
	queues.StartProcessQueue()
	queues.StartProcessRequestQueue()
	// Seta o id do servidor
	vectorClock.SetServerId(1)

	// Roda o server
	graphs.ReadRoutes() // Carregando o gráfico na memória
	router.Run(":"+ os.Getenv("LOCAL_PORT")) // Roda na porta 8080
}
