package passagesController

import (
	"net/http"
	passagesService "sharedPass/passages/services"

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
	confirmation, err := passagesService.Buy([]string{"Salvador", "Sao Paulo"})

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, confirmation)
}
