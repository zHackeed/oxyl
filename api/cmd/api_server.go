package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	fiberRecover "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/spf13/cobra"
	"zhacked.me/oxyl/api/internal/controllers/agent"
	agentAuth "zhacked.me/oxyl/api/internal/controllers/agent/auth"
	agentNotifications "zhacked.me/oxyl/api/internal/controllers/agent/notifications"
	"zhacked.me/oxyl/api/internal/controllers/company"
	"zhacked.me/oxyl/api/internal/controllers/company/endpoints"
	"zhacked.me/oxyl/api/internal/controllers/company/member"
	"zhacked.me/oxyl/api/internal/controllers/company/threshold"
	"zhacked.me/oxyl/api/internal/controllers/user"
	"zhacked.me/oxyl/api/internal/middlewares"
	apiModel "zhacked.me/oxyl/api/internal/models"
	"zhacked.me/oxyl/api/internal/validator"
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

	userStorage, companyStorage, agentStorage, notificationStorage, tokenStorage, monitoringStorage := createStorage(timescale, redis)

	userService, companyService, agentService, tokenService, metricsService, err := createServices(userStorage,
		companyStorage, agentStorage, notificationStorage, tokenStorage, monitoringStorage, redis)

	if err != nil {
		slog.Error("unable to create service", "error", err)
		return
	}

	ctx := cmd.Context()

	httpServer := fiber.New(fiber.Config{
		AppName:         fmt.Sprintf("oxyl-api-%s", version.CommitID),
		ServerHeader:    fmt.Sprintf("oxyl-api-%s", version.CommitID),
		CaseSensitive:   false,
		StructValidator: validator.NewStructValidator(),
	})

	httpServer.Use(fiberRecover.New(
		fiberRecover.Config{
			EnableStackTrace: true,
		},
	))

	httpServer.Use(requestid.New())
	httpServer.Use(middlewares.NewLoggingMiddleWare().Handle)
	httpServer.Use(healthcheck.LivenessEndpoint, healthcheck.New())
	/*
		httpServer.Use(limiter.New(
			limiter.Config{
				Next: func(c fiber.Ctx) bool {
					return c.IP() == "127.0.0.1" // do not apply rate limit to localhost
				},
				Max:        20,
				Expiration: 30 * time.Second,
				KeyGenerator: func(c fiber.Ctx) string {
					return c.Get("x-forwarded-for")
				},
				LimiterMiddleware: limiter.SlidingWindow{},
			}),
		)
	*/

	unprotectedRoutes := []apiModel.Registrable{
		// -------------- user routes
		user.NewRegisterController(userService),
		user.NewLoginController(userService, tokenService),
		user.NewRefreshController(tokenService),
		user.NewLogoutController(tokenService),

		// -------------- agent routes
		agentAuth.NewAgentRefreshController(tokenService),
		agentAuth.NewAgentLoginController(agentService, tokenService),
		agentAuth.NewAgentShutdownController(agentService, tokenService),
	}

	registerRoutes(httpServer, unprotectedRoutes...)

	apiGroupRoute := httpServer.Group("/api/v1")

	authMiddleware := middlewares.NewAuthMiddleware(tokenService)
	apiGroupRoute.Use(authMiddleware.Handle)

	protectedRoutes := []apiModel.Registrable{
		// -------------- user protected routes
		user.NewInfoController(userService),

		// -------------- company routes
		company.NewCreateCompanyController(companyService),
		company.NewListCompaniesController(companyService),
		company.NewInfoController(companyService),
		company.NewDeleteCompanyController(companyService),
		member.NewSelfPermissionController(companyService),
		member.NewAddMemberController(companyService),
		member.NewListMemberController(companyService),
		member.NewRemoveMemberController(companyService, userService),
		threshold.NewThresholdsController(companyService),
		threshold.NewModifyThresholdController(companyService),
		endpoints.NewCreateEntrypointController(companyService),
		endpoints.NewListEntrypointController(companyService),
		endpoints.NewDeleteEntrypointController(companyService),

		// -------------- agent routes
		agent.NewCreateAgentController(agentService),
		agent.NewListAgentsController(agentService),
		agent.NewAgentInfoController(agentService),
		agent.NewDeleteAgentController(agentService),
		agent.NewToggleMaintenanceController(agentService),
		agent.NewMetricsController(metricsService),
		agentNotifications.NewListController(agentService),
	}

	registerRoutes(apiGroupRoute, protectedRoutes...)

	go func() {
		slog.Info("starting http server", slog.Int("port", 19999))
		if err := httpServer.Listen(":19999", fiber.ListenConfig{
			DisableStartupMessage: true,
			EnablePrintRoutes:     true,
		}); err != nil {
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

func createStorage(timescale *datasource.TimescaleConnection, redis *datasource.RedisConnection) (
	*storage.UserStorage, *storage.CompanyStorage, *storage.AgentStorage, *storage.NotificationStorage, *storage.TokenStorage, *storage.MonitoringStorage,
) {
	return storage.NewUserStorage(timescale),
		storage.NewCompanyStorage(timescale),
		storage.NewAgentStorage(timescale),
		storage.NewNotificationStorage(timescale),
		storage.NewTokenStorage(redis),
		storage.NewMonitoringStorage(timescale)
}

func createServices(
	userStorage *storage.UserStorage,
	companyStorage *storage.CompanyStorage,
	agentStorage *storage.AgentStorage,
	notificationStorage *storage.NotificationStorage,
	tokenStorage *storage.TokenStorage,
	monitoringStorage *storage.MonitoringStorage,
	redis *datasource.RedisConnection,
) (*service.UserService, *service.CompanyService, *service.AgentService, *service.TokenService, *service.MetricsService, error) {
	tokenService, err := service.NewTokenService(tokenStorage, redis)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	userService := service.NewUserService(userStorage)
	companyService := service.NewCompanyService(redis, companyStorage, userStorage)
	agentService := service.NewAgentService(redis, companyStorage, agentStorage, notificationStorage)
	metricsService := service.NewAgentMetricsService(companyStorage, agentStorage, monitoringStorage)

	return userService, companyService, agentService, tokenService, metricsService, nil
}

func registerRoutes(router fiber.Router, registrable ...apiModel.Registrable) {
	for _, route := range registrable {
		if route.GetMethod() == apiModel.MethodNone {
			continue
		}

		if route.RequestRequirements() != nil {
			schemaValidator := middlewares.NewSchemaValidator(route.RequestRequirements().GetValidationType(), route.RequestRequirements().GetModel())
			router.Add([]string{string(route.GetMethod())}, route.GetPath(), schemaValidator.Handle, route.Handle)
		} else {
			router.Add([]string{string(route.GetMethod())}, route.GetPath(), route.Handle)
		}
	}
}
