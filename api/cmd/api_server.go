package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	bind "github.com/idan-fishman/fiber-bind"
	"github.com/spf13/cobra"
	"zhacked.me/oxyl/api/internal/controllers/agent"
	"zhacked.me/oxyl/api/internal/controllers/company"
	"zhacked.me/oxyl/api/internal/controllers/user"
	"zhacked.me/oxyl/api/internal/middlewares"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/shared/pkg/datasource"
	"zhacked.me/oxyl/shared/pkg/logger"
	"zhacked.me/oxyl/shared/pkg/service"
	"zhacked.me/oxyl/shared/pkg/storage"
	"zhacked.me/oxyl/shared/pkg/version"
)

var apiServer = cobra.Command{
	Use:   "",
	Short: "Start the api server",
	Run:   startAPIServer,
}

func init() {
	cobra.OnInitialize(func() {
		logger.Register(slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})
	})
}

func Execute(ctx context.Context) error {
	return apiServer.ExecuteContext(ctx)
}

func startAPIServer(cmd *cobra.Command, _ []string) {
	slog.Info("starting api server", slog.String("version", version.CommitID), slog.String("branch", version.Branch))

	initCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	timescale, redis, err := createDatabases(initCtx)
	if err != nil {
		slog.Error("unable to create databases", "error", err)
		return
	}
	defer func() {
		if err := redis.Close(); err != nil {
			slog.Error("unable to close redis connection", "error", err)
		}

		timescale.Close()
	}()

	userStorage, companyStorage, agentStorage, tokenStorage := createStorage(timescale, redis)
	userService, companyService, agentService, tokenService, err := createServices(userStorage, companyStorage, agentStorage, tokenStorage, redis)
	if err != nil {
		slog.Error("unable to create services", "error", err)
		return
	}

	ctx := cmd.Context()

	httpServer := fiber.New(fiber.Config{
		AppName:       "oxyl-api",
		ServerHeader:  "oxyl",
		CaseSensitive: false,
	})

	unprotectedRoutes := []apiModel.Registrable{
		user.NewRegisterController(userService),
		user.NewLoginController(userService, tokenService),
		user.NewLogoutController(tokenService),
	}

	registerRoutes(httpServer, unprotectedRoutes...)

	apiGroupRoute := httpServer.Group("/api/v1")

	authMiddleware := middlewares.NewAuthMiddleware(tokenService)
	apiGroupRoute.Use(authMiddleware.Handle)

	protectedRoutes := []apiModel.Registrable{
		// -------------- company routes
		company.NewCreateCompanyController(companyService),
		company.NewListCompaniesController(companyService),
		company.NewDeleteCompanyController(companyService),
		company.NewAddMemberController(companyService),
		company.NewRemoveMemberController(companyService),
		company.NewThresholdsController(companyService),
		company.NewModifyThresholdController(companyService),

		// -------------- agent routes
		agent.NewCreateAgentController(agentService),
		agent.NewListAgentsController(agentService),
		agent.NewDeleteAgentController(agentService),
		agent.NewToggleMaintenanceController(agentService),
	}

	registerRoutes(apiGroupRoute, protectedRoutes...)

	go func() {
		if err := httpServer.Listen(":19999"); err != nil {
			slog.Error("unable to start http server", "error", err)
		}
	}()

	<-ctx.Done()

	if err := httpServer.ShutdownWithTimeout(10 * time.Second); err != nil {
		slog.Error("unable to shutdown http server", "error", err)
	}

}

func createDatabases(ctx context.Context) (*datasource.TimescaleConnection, *datasource.RedisConnection, error) {
	timescale, err := datasource.NewTimescaleConnection(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect: %w", err)
	}

	redis, err := datasource.NewRedisConnection()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect: %w", err)
	}

	return timescale, redis, nil
}

func createStorage(timescale *datasource.TimescaleConnection, redis *datasource.RedisConnection) (*storage.UserStorage, *storage.CompanyStorage, *storage.AgentStorage, *storage.TokenStorage) {
	return storage.NewUserStorage(timescale), storage.NewCompanyStorage(timescale), storage.NewAgentStorage(timescale), storage.NewTokenStorage(redis)
}

func createServices(
	userStorage *storage.UserStorage,
	companyStorage *storage.CompanyStorage,
	agentStorage *storage.AgentStorage,
	tokenStorage *storage.TokenStorage,
	redis *datasource.RedisConnection,
) (*service.UserService, *service.CompanyService, *service.AgentService, *service.TokenService, error) {
	tokenService, err := service.NewTokenService(tokenStorage)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	userService := service.NewUserService(userStorage)
	companyService := service.NewCompanyService(redis, companyStorage)
	agentService := service.NewAgentService(redis, companyStorage, agentStorage)

	return userService, companyService, agentService, tokenService, nil
}

func registerRoutes(router fiber.Router, registrable ...apiModel.Registrable) {
	for _, route := range registrable {
		if route.GetMethod() == apiModel.MethodNone {
			continue
		}

		if route.GetRequestModel() != nil {
			router.Use(route.GetPath(), bind.New(bind.Config{
				Source: bind.JSON,
				// todo: create global struct validator
			}, route.GetRequestModel()))
		}

		router.Add([]string{string(route.GetMethod())}, route.GetPath(), route.Handle)
	}
}
