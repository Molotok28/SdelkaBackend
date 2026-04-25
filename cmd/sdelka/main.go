package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_auth "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/auth"
	"github.com/nikolaimoiseev01/SdelkaBackend/internal/core/database"
	core_logger "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/logger"
	core_http_middleware "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/middleware"
	core_http_server "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/server"
	ads_postgres_repository "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/ads/repository/postgres"
	ads_service "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/ads/service"
	ads_transport_http "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/ads/transport/http"
	auth_postgres_repository "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/auth/repository/postgres"
	auth_service "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/auth/service"
	auth_transport_http "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/auth/transport/http"
	users_transport_http "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger: ", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("SDELKA application starting")

	pool, err := database.NewPool(ctx, database.NewConfigMust())
	if err != nil {
		logger.Error("failed to connect to postgres", zap.Error(err))
		os.Exit(1)
	}
	defer pool.Close()

	logger.Debug("connected to postgres")

	jwtConfig := core_auth.NewConfigMust()

	authRepo := auth_postgres_repository.NewAuthPostgresRepository(pool)
	authSvc := auth_service.NewAuthService(authRepo, jwtConfig)
	authHandler := auth_transport_http.NewAuthHTTPHandler(authSvc, jwtConfig)

	adsRepo := ads_postgres_repository.NewAdsPostgresRepository(pool)
	adsSvc := ads_service.NewAdsService(adsRepo)
	adsHandler := ads_transport_http.NewAdsHTTPHandler(adsSvc, jwtConfig)

	usersHandler := users_transport_http.NewUsersHTTPHandler(nil)

	apiRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiRouter.RegisterRoutes(authHandler.Routes()...)
	apiRouter.RegisterRoutes(adsHandler.Routes()...)
	apiRouter.RegisterRoutes(usersHandler.Routes()...)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestId(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)
	httpServer.RegisterAPIRouters(apiRouter)
	httpServer.RegisterStaticDir("web")

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("failed to run http server", zap.Error(err))
	}
}
