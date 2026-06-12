package auth

import (
	"financial-data-aggregator-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type registerInput struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"displayName" binding:"required,min=4"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type userResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"displayName"`
}

type Handler struct {
	authService service.AuthService
}

func NewHandler(authService service.AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) Register(ctx *gin.Context) {
	var input registerInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(input.Email, input.Name, input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h *Handler) Login(ctx *gin.Context) {
	var input loginInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenString, err := h.authService.Login(input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}
