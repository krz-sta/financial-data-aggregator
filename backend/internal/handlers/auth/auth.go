package auth

import (
	"errors"
	"financial-data-aggregator-backend/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	DB     *gorm.DB
	JWTKey string
}

func NewHandler(db *gorm.DB, jwtKey string) *Handler {
	return &Handler{
		DB:     db,
		JWTKey: jwtKey,
	}
}

func (h *Handler) Register(ctx *gin.Context) {
	var input registerInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var count int64
	err := h.DB.Model(&models.User{}).Where("email = ?", input.Email).Count(&count).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count > 0 {
		ctx.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser := models.User{
		ID:           uuid.New(),
		Email:        input.Email,
		DisplayName:  input.Name,
		PasswordHash: string(passwordHash),
	}

	err = h.DB.Create(&newUser).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": userResponse{
		ID:          newUser.ID,
		Email:       newUser.Email,
		DisplayName: newUser.DisplayName,
	}})
}

func (h *Handler) Login(ctx *gin.Context) {
	var input loginInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tempUsr models.User
	err := h.DB.Model(&models.User{}).Where("email = ?", input.Email).First(&tempUsr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(tempUsr.PasswordHash), []byte(input.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// sub (Subject) - ID użytkownika
	// exp (Expiration Time) - czas wygaśnięcia tokenu w formacie Unix
	// iat (Issued at) - czas wydania tokenu w formacie Unix
	claims := jwt.MapClaims{
		"sub": tempUsr.ID.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.JWTKey))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}
