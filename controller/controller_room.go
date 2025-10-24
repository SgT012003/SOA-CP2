package controller

import (
	"net/http"

	"hotel-soa/model"
	"hotel-soa/service"

	"github.com/gin-gonic/gin"
)

// RoomController gerencia endpoints de quartos
type RoomController struct {
	service service.RoomService
}

// NewRoomController cria um novo RoomController
func NewRoomController(s service.RoomService) *RoomController {
	return &RoomController{service: s}
}

// @Summary Cria um novo quarto
// @Description Cria um novo quarto com os dados fornecidos
// @Tags rooms
// @Accept json
// @Produce json
// @Param room body model.RoomRequest true "Quarto"
// @Success 201 {object} model.Room
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /rooms [post]
func (rc *RoomController) Create(c *gin.Context) {
	var req model.RoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room := req.Room()
	if err := room.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := rc.service.Create(*room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	room.ID = id
	c.JSON(http.StatusCreated, room)
}

// @Summary Atualiza um quarto existente
// @Description Atualiza os dados de um quarto pelo ID
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path string true "ID do Quarto (UUID)"
// @Param room body model.RoomRequest true "Quarto atualizado"
// @Success 200 {object} model.Room
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /rooms/{id} [put]
func (rc *RoomController) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req model.RoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = id
	room := req.Room()
	if err := room.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := rc.service.Update(*room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, room)
}

// @Summary Deleta um quarto
// @Description Deleta um quarto pelo ID
// @Tags rooms
// @Param id path string true "ID do Quarto (UUID)"
// @Success 204
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /rooms/{id} [delete]
func (rc *RoomController) Delete(c *gin.Context) {
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

// @Summary Busca quarto pelo ID
// @Description Retorna um quarto pelo seu ID
// @Tags rooms
// @Param id path string true "ID do Quarto (UUID)"
// @Success 200 {object} model.Room
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /rooms/{id} [get]
func (rc *RoomController) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	room, err := rc.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, room)
}

// @Summary Lista todos os quartos
// @Description Retorna todos os quartos cadastrados
// @Tags rooms
// @Success 200 {array} model.Room
// @Success 204 "No Content"
// @Failure 500 {object} model.ErrorResponse
// @Router /rooms [get]
func (rc *RoomController) GetAll(c *gin.Context) {
	rooms, err := rc.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(rooms) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, rooms)
}
