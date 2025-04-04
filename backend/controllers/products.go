package controllers

import (
	"app/models"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductController struct {
	DB *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{DB: db}
}

func (pc *ProductController) GetProducts(c *gin.Context) {
	var products []models.Product
	
	if err := pc.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar produtos"})
		return
	}
	
	c.JSON(http.StatusOK, products)
}

func (pc *ProductController) GetProduct(c *gin.Context) {
	id := c.Param("id")
	
	var product models.Product
	if err := pc.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado"})
		return
	}
	
	c.JSON(http.StatusOK, product)
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	var product models.Product
	
	name := c.PostForm("name")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")
	quantityStr := c.PostForm("quantity")
	
	if name == "" || description == "" || priceStr == "" || quantityStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todos os campos são obrigatórios"})
		return
	}
	
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Preço inválido"})
		return
	}
	
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantidade inválida"})
		return
	}
	
	product.Name = name
	product.Description = description
	product.Price = price
	product.Quantity = quantity
	
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Imagem obrigatória"})
		return
	}
	
	extension := filepath.Ext(file.Filename)
	newFilename := uuid.New().String() + extension
	
	if err := c.SaveUploadedFile(file, "uploads/"+newFilename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar imagem"})
		return
	}
	
	product.Image = newFilename
	
	if err := pc.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar produto"})
		return
	}
	
	c.JSON(http.StatusCreated, product)
}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	
	var product models.Product
	if err := pc.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado"})
		return
	}
	
	name := c.PostForm("name")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")
	quantityStr := c.PostForm("quantity")
	
	if name != "" {
		product.Name = name
	}
	
	if description != "" {
		product.Description = description
	}
	
	if priceStr != "" {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Preço inválido"})
			return
		}
		product.Price = price
	}
	
	if quantityStr != "" {
		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Quantidade inválida"})
			return
		}
		product.Quantity = quantity
	}
	
	file, err := c.FormFile("image")
	if err == nil {
		extension := filepath.Ext(file.Filename)
		newFilename := uuid.New().String() + extension
		
		if err := c.SaveUploadedFile(file, "uploads/"+newFilename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar imagem"})
			return
		}
		
		product.Image = newFilename
	}
	
	if err := pc.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao atualizar produto"})
		return
	}
	
	c.JSON(http.StatusOK, product)
}

func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	
	var product models.Product
	if err := pc.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado"})
		return
	}
	
	if err := pc.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao excluir produto"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Produto excluído com sucesso"})
}