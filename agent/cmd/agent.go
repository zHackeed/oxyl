package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"zhacked.me/oxyl/agent/internal/logger"
	"zhacked.me/oxyl/agent/internal/service"
	protocolV1 "zhacked.me/oxyl/protocol/v1"
)

var (
	agentId          string
	grpcEndpoint     string
	loginEndpoint    string
	refreshEndpoint  string
	shutdownEndpoint string
)

var agentCmd = &cobra.Command{
	Use:   "",
	Short: "Monitor agent",
	Run:   startAgent,
}

func Execute(ctx context.Context) error {
	return agentCmd.ExecuteContext(ctx)
}

func init() {
	cobra.OnInitialize(func() {
		logger.Register(slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	})

	agentCmd.PersistentFlags().StringVarP(&grpcEndpoint, "grpc-endpoint", "g", "https://ingress.oxyl.zhacked.me", "gRPC endpoint")
	agentCmd.PersistentFlags().StringVarP(&agentId, "id", "i", "", "Agent ID")
	agentCmd.PersistentFlags().StringVarP(&loginEndpoint, "login-endpoint", "l", "https://api.oxyl.zhacked.me/agent/auth/login", "Login endpoint")
	agentCmd.PersistentFlags().StringVarP(&refreshEndpoint, "refresh-endpoint", "r", "https://api.oxyl.zhacked.me/agent/auth/refresh", "Refresh endpoint")
	agentCmd.PersistentFlags().StringVarP(&shutdownEndpoint, "shutdown-endpoint", "s", "https://api.oxyl.zhacked.me/agent/auth/shutdown", "Shutdown endpoint")

	if err := agentCmd.MarkPersistentFlagRequired("id"); err != nil {
		slog.Error("failed to mark flag required", "error", err)
		os.Exit(1)
	}
}

func startAgent(cmd *cobra.Command, _ []string) {
	slog.Info("starting agent", slog.String("id", agentId))
	authService, err := service.NewAuthService(agentId, loginEndpoint, refreshEndpoint, shutdownEndpoint)
	if err != nil {
		slog.Error("failed to create authentication service", "error", err)
		os.Exit(1)
	}

	authService.StartTicking(cmd.Context())

	systemInfoService, err := service.NewSystemInfoService()
	if err != nil {
		slog.Error("failed to create system info service", "error", err)
		os.Exit(1)
	}

	if err := systemInfoService.CaptureData(); err != nil {
		slog.Error("failed to capture system info", "error", err)
		os.Exit(1)
	}

	/*
		//todo: enforce TLS
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			slog.Error("failed to load system cert pool", "error", err)
			os.Exit(1)
		}
	*/

	grpcClient, err := grpc.NewClient(grpcEndpoint,
		grpc.WithPerRPCCredentials(authService),
		/*
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				RootCAs: systemRoots,
			}))
		*/
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		slog.Error("failed to create gRPC client", "error", err)
		os.Exit(1)
	}
	defer grpcClient.Close()

	enrollmentService := service.NewEnrollmentService(systemInfoService)
	enrollmentService.EnrollmentServiceClient = protocolV1.NewEnrollmentServiceClient(grpcClient)

	if err := enrollmentService.Start(cmd.Context()); err != nil {
		slog.Error("failed to start enrollment service", "error", err)
		os.Exit(1)
	}

	token, err := enrollmentService.ProvideEnrollmentIdentifier()
	if err == nil {
		authService.SetEnrollmentToken(&token)
	}

	monitoringService, err := service.NewMonitoringService(systemInfoService)

	if err != nil {
		slog.Error("failed to create monitoring service", "error", err)
		os.Exit(1)
	}

	monitoringService.MonitoringServiceClient = protocolV1.NewMonitoringServiceClient(grpcClient)

	if err := monitoringService.Start(cmd.Context()); err != nil {
		slog.Error("failed to start monitoring service", "error", err)
		os.Exit(1)
	}

	// wait for signal
	<-cmd.Context().Done()

	err = authService.RequestShutdown()
	if err != nil {
		slog.Error("failed to request shutdown authentication invalidation", "error", err)
	}

}
