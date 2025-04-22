package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"

	prettyconsole "github.com/thessem/zap-prettyconsole" //logger
	"go.uber.org/zap"                                    // logger
	"go.uber.org/zap/zapcore"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title Лавка Gateway API
// @version 1.0
// @description API Gateway для проекта Лавка
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/ultard/fusion

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg := config.MustLoadConfig()

	logger := setupLogger(cfg.Env)
	defer logger.Sync()

	logger.Info("config loaded", zap.Any("config", cfg))

	authClient, err := rpc.NewAuth(cfg.Clients.Auth.Address + ":" + cfg.Clients.Auth.Port)

	if err != nil {
		logger.Fatal("failed to create auth client", zap.Error(err))
		os.Exit(1)
	}

	proCusClient, err := rpc.NewProCus(cfg.Clients.Products_customer.Address + ":" + cfg.Clients.Products_customer.Port)

	if err != nil {
		logger.Fatal("failed to create products-customer client", zap.Error(err))
		os.Exit(1)
	}

	proAdmClient, err := rpc.NewProAdm(cfg.Clients.Products_admin.Address + ":" + cfg.Clients.Products_admin.Port)

	if err != nil {
		logger.Fatal("failed to create products-admin client", zap.Error(err))
		os.Exit(1)
	}

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth.RegisterAuthRoutes(r, logger, authClient, []byte(cfg.Clients.Auth.AppSecret))
	proCus.RegisterProCusRoutes(r, logger, proCusClient)
	proAdm.RegisterProAdmRoutes(r, logger, proAdmClient)

	r.Run(cfg.Port)
}

func setupLogger(env string) *zap.Logger {
	var logger *zap.Logger
	switch env {
	case envLocal:
		logger = prettyconsole.NewLogger(zapcore.DebugLevel)
	case envDev:
		atomic := zap.NewAtomicLevel()
		atomic.SetLevel(zapcore.DebugLevel)
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.Lock(os.Stdout),
			atomic,
		)
		logger = zap.New(core)
	case envProd:
		atomic := zap.NewAtomicLevel()
		atomic.SetLevel(zapcore.InfoLevel)
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.Lock(os.Stdout),
			atomic,
		)
		logger = zap.New(core)
	}
	return logger
}
