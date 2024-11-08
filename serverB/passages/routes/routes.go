package routes

import (
	passagesController "sharedPass/passages/controllers"

	"github.com/gin-gonic/gin"
)
// Registra todas as rotas possiveis de passagens
func RegisterRoutes(router *gin.Engine) {

	api := router.Group("/passages")
	{
		//Busca todas as rotas possiveis
		api.GET("/routes", passagesController.FindAllRoutes)
		//Busca os todos os vôos de um servidor
		api.GET("/flights", passagesController.FindAllFlights)
		//Compra uma passagem
		api.POST("/buy", passagesController.Buy)
		//Cancela uma passagem
		api.POST("/rollback", passagesController.RollBack)
	}

}
