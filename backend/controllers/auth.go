package controllers

import (
	"app/models"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (ac *AuthController) Register(c *gin.Context) {
	var user models.User
	
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	
	if name == "" || email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome, email e senha são obrigatórios"})
		return
	}
	
	var existingUser models.User
	if err := ac.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email já cadastrado"})
		return
	}
	
	user.Name = name
	user.Email = email
	user.Password = password
	
	file, err := c.FormFile("profilePic")
	if err == nil {
		extension := filepath.Ext(file.Filename)
		newFilename := uuid.New().String() + extension
		
		if err := c.SaveUploadedFile(file, "uploads/"+newFilename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar imagem"})
			return
		}
		
		user.ProfilePic = newFilename
	}
	
	if err := ac.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar usuário"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "Usuário registrado com sucesso"})
}

func (ac *AuthController) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	
	var user models.User
	if err := ac.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}
	
	if !user.CheckPassword(input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	
	tokenString, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gerar token"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"profilePic": user.ProfilePic,
		},
	})
}

func (ac *AuthController) Verify(c *gin.Context) {
	userId, _ := c.Get("userId")
	
	var user models.User
	if err := ac.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"profilePic": user.ProfilePic,
		},
	})
}