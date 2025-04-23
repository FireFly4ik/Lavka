package main

import (
	"fmt"
	"github.com/FireFly4ik/Lavka-auth/internal/auth"
	"github.com/FireFly4ik/Lavka-auth/internal/config"
	"github.com/FireFly4ik/Lavka-auth/internal/db"
	"github.com/FireFly4ik/Lavka-auth/pkg/utils"
	pb "github.com/FireFly4ik/Lavka-auth/proto"
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
	jwtService := utils.NewJWTService(cfg.JWTSecretKey)
	emailService := utils.NewEmailService(cfg.SMTP.SMTPHost, cfg.SMTP.SMTPPort, cfg.SMTP.SMTPUsername, cfg.SMTP.SMTPPassword, cfg.SMTP.EmailFrom)
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	authService := auth.AuthServiceServer{
		Logger:       zapLogger,
		JwtService:   jwtService,
		EmailService: emailService,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &authService)

	log.Printf("gRPC server starting on port %s", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
