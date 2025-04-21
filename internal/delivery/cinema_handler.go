package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/nurbol/cinema/internal/models"
	"github.com/nurbol/cinema/internal/services"
	"net/http"
	"strconv"
)

type CinemaHandler struct {
	service *services.CinemaService
}

func NewCinemaHandler(service *services.CinemaService) *CinemaHandler {
	return &CinemaHandler{service: service}
}

func (h *CinemaHandler) GetAllCinemas(c *gin.Context) {
	cinemas, err := h.service.GetAllCinemas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve cinemas"})
		return
	}
	c.JSON(http.StatusOK, cinemas)
}

func (h *CinemaHandler) GetCinema(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cinema ID"})
		return
	}
	cinema, err := h.service.GetCinemaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cinema not found"})
		return
	}
	c.JSON(http.StatusOK, cinema)
}

func (h *CinemaHandler) CreateCinema(c *gin.Context) {
	var cinema models.Cinema
	if err := c.ShouldBindJSON(&cinema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := h.service.CreateCinema(&cinema); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create cinema"})
		return
	}
	c.JSON(http.StatusCreated, cinema)
}

func (h *CinemaHandler) UpdateCinema(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cinema ID"})
		return
	}
	var cinema models.Cinema
	if err := c.ShouldBindJSON(&cinema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := h.service.UpdateCinema(id, &cinema); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cinema not found"})
		return
	}
	c.JSON(http.StatusOK, cinema)
}

func (h *CinemaHandler) DeleteCinema(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cinema ID"})
		return
	}
	if err := h.service.DeleteCinema(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cinema not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cinema deleted successfully"})
}
