package main

import (
	"fmt"
	"log"
	"net/http"
	"project/config"
	"project/delivery/handlers"
	"project/delivery/routes"
	adminrepository "project/repository/admin"
	cartrepository "project/repository/cart"
	"project/repository/infrastructure"
	orderrepository "project/repository/order"
	productrepository "project/repository/product"
	repository "project/repository/user"
	adminUseCase "project/usecase/admin"
	cartusecase "project/usecase/cart"
	orderusecase "project/usecase/order"
	productusecase "project/usecase/product"
	usecase "project/usecase/user"

	"github.com/gin-gonic/gin"
)

func main() {
	config,err:=config.LoadConfig()
	if err != nil{
		log.Fatal("error loading files using viper")
	}
	db, err := infrastructure.ConnectDb(config.DB)
	if err != nil {
		log.Fatal(err)
	}
	userRepo := repository.NewUserRepository(db)
	adminRepo := adminrepository.NewAdminRepository(db)
	productRepo := productrepository.NewProductRepository(db)
	cartRepo := cartrepository.NewCartRepository(db)
	orderRepo := orderrepository.NewOrderRepository(db)

	userusecase := usecase.NewUser(userRepo,&config.Otp)
	adminUseCase := adminUseCase.NewAdmin(adminRepo)
	productUsecase := productusecase.NewProduct(productRepo,&config.S3aws)
	cartUsecase := cartusecase.NewCart(cartRepo, productRepo)
	orderUsecase := orderusecase.NewOrder(orderRepo, cartRepo, userRepo, productRepo,&config.Razopay)

	userHandler := handlers.NewUserhandler(userusecase, productUsecase, cartUsecase)
	adminHandler := handlers.NewAdminHandler(adminUseCase, productUsecase)
	orderHandler := handlers.NewOrderHandler(orderUsecase,config.Razopay)

	router := gin.Default()

	routes.UserRouter(router, userHandler)
	routes.AdminRouter(router, adminHandler)
	routes.OrderRouter(router, orderHandler)

	router.LoadHTMLGlob("template/*.html")
	fmt.Println("Templates loaded from:", "template/*.html")
	fmt.Println("Starting server on port 8080...")
	err1 := http.ListenAndServe(":8080", router)
	if err1 != nil {
		log.Fatal(err1)
	}
}
