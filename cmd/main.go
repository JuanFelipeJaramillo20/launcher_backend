package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"venecraft-back/cmd/controller"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/repository"
	"venecraft-back/cmd/routes"
	"venecraft-back/cmd/seeds"
	"venecraft-back/cmd/service"
)

var DB *gorm.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Ensure you have set the environment variables.")
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

	fmt.Println("Database migrated successfully!")
}

func main() {
	connectDatabase()

	userRepo := repository.NewUserRepository(DB)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	roleRepo := repository.NewRoleRepository(DB)
	userRoleRepo := repository.NewUserRoleRepository(DB)

	registerRepo := repository.NewRegisterRepository(DB)
	registerService := service.NewRegisterService(registerRepo, userRepo, roleRepo, userRoleRepo)
	registerController := controller.NewRegisterController(registerService)

	server := gin.Default()

	routes.RegisterRoutes(server, registerController)
	routes.UserRoutes(server, userController)

	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server up and running",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := server.Run(":" + port)
	if err != nil {
		return
	}
}
