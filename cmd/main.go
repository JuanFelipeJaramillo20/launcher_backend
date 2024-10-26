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
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"venecraft-back/cmd/controller"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/middlewares"
	"venecraft-back/cmd/repository"
	"venecraft-back/cmd/routes"
	"venecraft-back/cmd/seeds"
	"venecraft-back/cmd/service"
)

var DB *gorm.DB

func init() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
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

	// Initialize services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)
	registerService := service.NewRegisterService(registerRepo, userRepo, roleRepo, userRoleRepo)

	// Initialize controllers
	userController := controller.NewUserController(userService)
	authController := controller.NewAuthController(authService)
	registerController := controller.NewRegisterController(registerService)

	server := gin.Default()

	routes.AuthRoutes(server, authController)
	routes.RegisterRoutes(server, registerController)

	protected := server.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.PUT("/register/approve/:id", registerController.ApproveRegister)
		routes.UserRoutes(protected, userController)
	}

	// Health check route
	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server up and running",
		})
	})

	// Swagger UI endpoint
	server.GET("/docs", func(c *gin.Context) {
		opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
		sh := middleware.SwaggerUI(opts, nil)
		sh.ServeHTTP(c.Writer, c.Request)
	})

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
