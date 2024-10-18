package routes

import (
	passagesController "sharedPass/passages/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	api := router.Group("/passages")
	{
		api.GET("/routes", passagesController.FindAll)
		api.POST("/buy", passagesController.Buy)
	}

}
