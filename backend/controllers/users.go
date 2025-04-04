package controllers

import (
	"app/models"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) GetUsers(c *gin.Context) {
	var users []models.User
	
	if err := uc.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar usuários"})
		return
	}
	
	c.JSON(http.StatusOK, users)
}

func (uc *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	
	var user models.User
	if err := uc.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	
	if name == "" || email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome, email e senha são obrigatórios"})
		return
	}
	
	var existingUser models.User
	if err := uc.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
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
	
	if err := uc.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar usuário"})
		return
	}
	
	c.JSON(http.StatusCreated, user)
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	
	var user models.User
	if err := uc.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	
	if name != "" {
		user.Name = name
	}
	
	if email != "" && email != user.Email {
		var existingUser models.User
		if err := uc.DB.Where("email = ? AND id != ?", email, id).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email já cadastrado"})
			return
		}
		user.Email = email
	}
	
	if password != "" {
		user.Password = password
	}
	
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
	
	if err := uc.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao atualizar usuário"})
		return
	}
	
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	
	var user models.User
	if err := uc.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	
	if err := uc.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao excluir usuário"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Usuário excluído com sucesso"})
}