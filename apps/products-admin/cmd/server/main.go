package main

import (
	"fmt"
	"github.com/FireFly4ik/Lavka-products-admin/internal/config"
	"github.com/FireFly4ik/Lavka-products-admin/internal/db"
	proAdm "github.com/FireFly4ik/Lavka-products-admin/internal/products-admin"
	pb "github.com/FireFly4ik/Lavka-products-admin/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.MustLoadConfig()

	if err := db.Connect(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB.Address, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.Port,
	)); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Initialize services
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	products_adminService := proAdm.ProductsAdminServiceServer{
		Logger: zapLogger,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProductsAdminServiceServer(grpcServer, &products_adminService)

	log.Printf("gRPC server starting on port %s", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
