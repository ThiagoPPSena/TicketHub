package passagesController

import (
	"net/http"
	"sharedPass/collections"
	passagesService "sharedPass/passages/services"
	"sharedPass/vectorClock"

	"github.com/gin-gonic/gin"
)

func FindAllRoutes(context *gin.Context) {
	origin := context.Query("origem")       // Captura o parâmetro de query 'origem'
	destination := context.Query("destino") // Captura o parâmetro de query 'destino'

	routes, err := passagesService.GetAllRoutes(origin, destination)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, routes)
}

func FindAllFlights(context *gin.Context) {
	flights, err := passagesService.GetAllFlights()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, flights)
}

func Buy(context *gin.Context) {

	var data collections.Body

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
	_, err := passagesService.Buy(data.Routes, serverId, externalClock)

	if err != nil {
		context.JSON(http.StatusConflict, gin.H{"error": "Passagem já reservada ou esgotada"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "Compra realizada com sucesso!",
	})
}

func RollBack(context *gin.Context){
	var data collections.Body

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
	_, err := passagesService.RollBack(data.Routes, serverId, externalClock)

	if err != nil {
		context.JSON(http.StatusConflict, gin.H{"error": "Passagem já reservada ou esgotada"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "Compra realizada com sucesso!",
	})
}