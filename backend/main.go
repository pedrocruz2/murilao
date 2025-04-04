package main

import (
	"app/controllers"
	"app/middleware"
	"app/models"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	dbHost := getEnv("DB_HOST", "postgres")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "appdb")
	dbPort := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&models.User{}, &models.Product{})

	r := gin.Default()

	r.Use(cors.Default())

	r.Static("/uploads", "./uploads")

	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", 0755)
	}

	authController := controllers.NewAuthController(db)
	userController := controllers.NewUserController(db)
	productController := controllers.NewProductController(db)

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.GET("/verify", middleware.AuthMiddleware(), authController.Verify)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		users := api.Group("/users")
		{
			users.GET("", userController.GetUsers)
			users.GET("/:id", userController.GetUser)
			users.POST("", userController.CreateUser)
			users.PUT("/:id", userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}

		products := api.Group("/products")
		{
			products.GET("", productController.GetProducts)
			products.GET("/:id", productController.GetProduct)
			products.POST("", productController.CreateProduct)
			products.PUT("/:id", productController.UpdateProduct)
			products.DELETE("/:id", productController.DeleteProduct)
		}
	}

	port := getEnv("PORT", "8080")
	r.Run(":" + port)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func saveFile(file *gin.Context, fieldName string) (string, error) {
	formFile, err := file.FormFile(fieldName)
	if err != nil {
		return "", err
	}

	filename := filepath.Base(formFile.Filename)
	uploadPath := filepath.Join("uploads", filename)

	if err := file.SaveUploadedFile(formFile, uploadPath); err != nil {
		return "", err
	}

	return filename, nil
}