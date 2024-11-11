package main

// Package classification of Venecraft API.
//
// Documentation for Venecraft API.
//
//  Schemes: http, https
//  Host: localhost:8080
//  BasePath: /
//  Version: 1.0.0
//  Consumes:
//  - application/json
//  Produces:
//  - application/json
//
// swagger:meta

import (
	"fmt"
	"log"
	"os"
	"time"
	"venecraft-back/cmd/controller"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/middlewares"
	"venecraft-back/cmd/repository"
	"venecraft-back/cmd/routes"
	"venecraft-back/cmd/seeds"
	"venecraft-back/cmd/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "production"
	}

	envFile := fmt.Sprintf(".env.%s", env)
	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("Error loading %s file: %v", envFile, err)
		log.Println("Falling back to default .env.production file")
		err = godotenv.Load()
		if err != nil {
			log.Println("No .env.production file found. Ensure you have set the environment variables.")
		}
	}
}

func connectDatabase() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	log.Println("Connecting to database...", dsn)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	fmt.Println("Database connection established successfully!")

	err = DB.AutoMigrate(&entity.Register{}, &entity.User{}, &entity.Role{}, &entity.Permission{},
		&entity.RolePermission{}, &entity.UserRole{}, &entity.Server{},
		&entity.Player{}, &entity.Ban{}, &entity.Log{}, &entity.Setting{},
		&entity.UserSetting{}, &entity.News{})
	if err != nil {
		log.Fatal("Failed to migrate the database: ", err)
	}

	seeds.SeedRoles(DB)
	seeds.SeedUsers(DB)

	fmt.Println("Database migrated successfully!")
}

func main() {
	connectDatabase()

	// Initialize repositories
	userRepo := repository.NewUserRepository(DB)
	roleRepo := repository.NewRoleRepository(DB)
	userRoleRepo := repository.NewUserRoleRepository(DB)
	registerRepo := repository.NewRegisterRepository(DB)
	newsRepo := repository.NewNewsRepository(DB)

	// Initialize services
	userService := service.NewUserService(userRepo, roleRepo)
	authService := service.NewAuthService(userRepo)
	registerService := service.NewRegisterService(registerRepo, userRepo, roleRepo, userRoleRepo)
	newsService := service.NewNewsService(newsRepo)

	// Initialize controllers
	userController := controller.NewUserController(userService)
	authController := controller.NewAuthController(authService)
	registerController := controller.NewRegisterController(registerService)
	newsController := controller.NewNewsController(newsService)

	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	server.POST("/api/password-reset-request", userController.PasswordResetRequest)
	server.POST("/api/reset-password", userController.ResetPassword)
	routes.AuthRoutes(server, authController)
	routes.RegisterRoutes(server, registerController)

	protected := server.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.PUT("/register/approve/:id", registerController.ApproveRegister)
		protected.PUT("/register/deny/:id", registerController.DenyRegister)
		protected.GET("/register", registerController.GetAllRegisters)
		routes.UserRoutes(protected, userController)
		routes.NewsRoutes(protected, newsController)
	}

	// Health check route
	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server up and running",
		})
	})

	server.GET("/docs", func(c *gin.Context) {
		opts := middleware.SwaggerUIOpts{
			SpecURL: "/swagger.yaml",
			Title:   "Venecraft API Documentation",
		}
		sh := middleware.SwaggerUI(opts, nil)
		sh.ServeHTTP(c.Writer, c.Request)
	})

	// Serve the Swagger spec
	server.StaticFile("/swagger.yaml", "./swagger.yaml")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := server.Run(":" + port)
	if err != nil {
		return
	}
}
