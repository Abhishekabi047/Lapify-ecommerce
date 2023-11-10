package main

import (
	"fmt"
	"log"
	"net/http"
	"project/delivery/handlers"
	"project/delivery/routes"
	adminrepository "project/repository/admin"
	"project/repository/infrastructure"
	productrepository "project/repository/product"
	repository "project/repository/user"
	adminUseCase "project/usecase/admin"
	productusecase "project/usecase/product"
	usecase "project/usecase/user"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := infrastructure.ConnectDb()
	if err != nil {
		log.Fatal(err)
	}
	userRepo := repository.NewUserRepository(db)
	adminRepo := adminrepository.NewAdminRepository(db)
	productRepo := productrepository.NewProductRepository(db)

	userusecase := usecase.NewUser(userRepo)
	adminUseCase := adminUseCase.NewAdmin(adminRepo)
	productUsecase := productusecase.NewProduct(productRepo)

	userHandler := handlers.NewUserhandler(userusecase, productUsecase)
	adminHandler := handlers.NewAdminHandler(adminUseCase, productUsecase)

	router := gin.Default()
	// router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.UserRouter(router, userHandler)
	routes.AdminRouter(router, adminHandler)
	// routes.OrderRouter(router, orderHandler)

	fmt.Println("Starting server on port 8080...")
	err1 := http.ListenAndServe(":8080", router)
	if err1 != nil {
		log.Fatal(err1)
	}
}
