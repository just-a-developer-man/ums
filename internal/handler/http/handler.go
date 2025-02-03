package handler_http

import (
	"log/slog"
	"net/http"
	"ums/internal/models"

	"github.com/gin-gonic/gin"
)

type errMessage string

var (
	invalidJSON    errMessage = "invalid JSON format"
	creationFailed errMessage = "failed to create user"
)

type UMS interface {
	CreateUser(user models.User) (int, error)
	ReadUsers() ([]models.User, error)
	UpdateUser(id int, user models.User) error
	DeleteUser(id int) error
}

// Handler is http handler responsible for handling RESTful requests to UMS service
type Handler struct {
	service UMS
}

// New creates new http handler with given service
func New(service UMS) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateUser is the hand for user creation
func (h *Handler) CreateUser(c *gin.Context) {
	user := models.User{}

	err := c.ShouldBindJSON(&user)
	if err != nil {
		slog.Error("Cannot bind body request to User model", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidJSON})
	}

	err = validateUser(user)
	if err != nil {
		errorMessage := getValidationErrorMessage(err)
		slog.Error("Failed to validate request body", "err", errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
	}

	id, err := h.service.CreateUser(user)
	if err != nil {
		slog.Error("Failed to create user", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": creationFailed})
	}

	slog.Info("User created", "id", id)
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
