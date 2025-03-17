package main

import (
	"auth0_test/internal/application"
	"auth0_test/internal/infrastructure/db"
	"auth0_test/internal/infrastructure/repository"
	"auth0_test/internal/interfaces/http"
)

func main() {
	db.InitDB()
	userRepo := repository.NewUserRepository(db.DB)
	userService := application.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userService)
	router := http.SetupRouter(userHandler)
	router.Run(":8080")
}
