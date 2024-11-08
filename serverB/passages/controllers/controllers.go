package passagesController

import (
	"net/http"
	"sharedPass/collections"
	passagesService "sharedPass/passages/services"
	"sharedPass/vectorClock"

	"github.com/gin-gonic/gin"
)

// Função para retornar todas as rotas da a origem e destino
func FindAllRoutes(context *gin.Context) {
	origin := context.Query("origem")       // Captura o parâmetro de query 'origem'
	destination := context.Query("destino") // Captura o parâmetro de query 'destino'

	routes, err := passagesService.GetAllRoutes(origin, destination) // Chama o service para pegar todas as rotas

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, routes)
}

// Fução que solicita todas as rotas (todo o arquivo JSON) de um servidor
func FindAllFlights(context *gin.Context) {
	flights, err := passagesService.GetAllFlights() // Chama o service

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, flights)
}

// Função de compra chamada por um servidor ou cliente
func Buy(context *gin.Context) {

	var data collections.Body // Estabelece formato do Body a ser enviado na requisição de compra

	// Faz o binding dos dados do corpo da requisição para a estrutura `BuyRequest`
	if err := context.BindJSON(&data); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	serverId := vectorClock.ServerId // Valor Padrão pro server A
	if data.ServerId != nil {
		serverId = *data.ServerId
	}
	vectorClock.LocalClock.Increment() // Incrementa o relógio local
	externalClock := vectorClock.LocalClock.Copy()
	if data.Clock != nil {
		externalClock = *data.Clock
	}

	// Passa as rotas, o ID do servidor que requisita a compra (se for um cliente, o ID é do próprio server) e passa o relógio do servidor que requisitou (se foi o cliente, o clock é do próprio servidor)
	respBuy, err := passagesService.Buy(data.Routes, serverId, externalClock) // Chama o service de compra

	if err != nil || !respBuy { // Se houve erro na compra ou se a resposta for nula, a compra não pôde ser realizada
		context.JSON(http.StatusConflict, gin.H{"error": "Passagem já reservada ou esgotada"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "Compra realizada com sucesso!",
	})
}

// Função de rollback
func RollBack(context *gin.Context){
	var data collections.Body // Estabelece formato do Body a ser enviado na requisição de rollback 

	// Faz o binding dos dados do corpo da requisição para a estrutura `BuyRequest`
	if err := context.BindJSON(&data); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	serverId := vectorClock.ServerId // Valor Padrão pro server A
	if data.ServerId != nil {
		serverId = *data.ServerId
	}
	vectorClock.LocalClock.Increment() // Incrementa o relógio local
	externalClock := vectorClock.LocalClock.Copy()
	if data.Clock != nil {
		externalClock = *data.Clock
	}
	// Passa as rotas, o ID do servidor que requisita o rollback e passa o relógio do servidor que requisitou o rollback
	respRollBack, err := passagesService.RollBack(data.Routes, serverId, externalClock)

	if err != nil || !respRollBack { // Se houve erro ou o rollback não funcionou
		context.JSON(http.StatusConflict, gin.H{"error": "Rollback mal sucedido"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "Rollback realizado com sucesso!",
	})
}
