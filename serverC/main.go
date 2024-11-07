package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sharedPass/graphs"
	"sharedPass/passages/routes"
	"sharedPass/queues"
	"sharedPass/vectorClock"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("AAAAAAAAAAA")
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
	serverID, err := strconv.Atoi(os.Getenv("SERVER_ID"))
	if err != nil {
		log.Fatal("Invalid SERVER_ID: ", err)
	}
	vectorClock.SetServerId(serverID)

	// Roda o server
	graphs.ReadRoutes()                       // Carregando o gráfico na memória
	fmt.Println("BBBBBBBBBBBB")
	router.Run(":" + os.Getenv("LOCAL_PORT")) // Roda na porta 8080
}
