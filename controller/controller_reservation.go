package controller

import (
	"net/http"

	"hotel-soa/model"
	"hotel-soa/service"

	"github.com/gin-gonic/gin"
)

// ReservationController gerencia endpoints de reservas
type ReservationController struct {
	service service.ReservationService
}

// NewReservationController cria um novo ReservationController
func NewReservationController(s service.ReservationService) *ReservationController {
	return &ReservationController{service: s}
}

// @Summary Cria uma nova reserva
// @Description Cria uma nova reserva com os dados fornecidos
// @Tags reservations
// @Accept json
// @Produce json
// @Param reservation body model.ReservationResponse true "Reserva"
// @Success 201 {object} model.Reservation
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /reservations [post]
func (rc *ReservationController) Create(c *gin.Context) {
	var req model.ReservationResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := req.Reservation()
	id, status, err := rc.service.Create(*res)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	res.ID = id
	c.JSON(http.StatusCreated, res)
}

// @Summary Atualiza uma reserva existente
// @Description Atualiza os dados de uma reserva pelo ID
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "ID da Reserva (UUID)"
// @Param reservation body model.ReservationResponse true "Reserva atualizada"
// @Success 200 {object} model.Reservation
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /reservations/{id} [put]
func (rc *ReservationController) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req model.ReservationResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = id
	res := req.Reservation()
	if status, err := rc.service.Update(*res); err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Deleta uma reserva
// @Description Deleta uma reserva pelo ID
// @Tags reservations
// @Param id path string true "ID da Reserva (UUID)"
// @Success 204
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /reservations/{id} [delete]
func (rc *ReservationController) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := rc.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Busca reserva pelo ID
// @Description Retorna uma reserva pelo seu ID
// @Tags reservations
// @Param id path string true "ID da Reserva (UUID)"
// @Success 200 {object} model.Reservation
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /reservations/{id} [get]
func (rc *ReservationController) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	res, err := rc.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Lista todas as reservas
// @Description Retorna todas as reservas cadastradas
// @Tags reservations
// @Success 200 {array} model.Reservation
// @Success 204 "No Content"
// @Failure 500 {object} model.ErrorResponse
// @Router /reservations [get]
func (rc *ReservationController) GetAll(c *gin.Context) {
	reservations, err := rc.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(reservations) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, reservations)
}
