// @title ERP Hotelaria SOA API
// @version 1.0
// @description API de exemplo com Gin + Swagger
// @host localhost:8080
// @BasePath /
package main

import (
	"hotel-soa/controller"
	"hotel-soa/service"
	"net/http"

	_ "hotel-soa/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @contact.name
// @contact.url
// @contact.email
func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	roomController := controller.NewRoomController(service.NewRoomService())
	reservationController := controller.NewReservationController(service.NewReservationService())

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Endpoint raiz
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/swagger/index.html")
	})

	// URI e handlers para rooms
	rooms := r.Group("/rooms")
	{
		rooms.POST("/", roomController.Create)
		rooms.PUT("/:id", roomController.Update)
		rooms.DELETE("/:id", roomController.Delete)
		rooms.GET("/:id", roomController.GetByID)
		rooms.GET("/", roomController.GetAll)
	}

	reservation := r.Group("/reservation")
	{
		reservation.POST("/", reservationController.Create)
		reservation.PUT("/:id", reservationController.Update)
		reservation.DELETE("/:id", reservationController.Delete)
		reservation.GET("/:id", reservationController.GetByID)
		reservation.GET("/", reservationController.GetAll)
	}

	// Inicia o servidor
	r.Run("0.0.0.0:8080")
}
