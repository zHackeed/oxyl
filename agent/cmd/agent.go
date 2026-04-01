package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"zhacked.me/oxyl/agent/internal/logger"
	"zhacked.me/oxyl/agent/internal/service"
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

	agentCmd.Flags().StringVarP(&grpcEndpoint, "grpc-endpoint", "g", "https://ingress.oxyl.zhacked.me", "gRPC endpoint")
	agentCmd.Flags().StringVarP(&agentId, "id", "i", "", "Agent ID")
	agentCmd.Flags().StringVarP(&loginEndpoint, "login-endpoint", "l", "https://api.oxyl.zhacked.me/agent/auth/login", "Login endpoint")
	agentCmd.Flags().StringVarP(&refreshEndpoint, "refresh-endpoint", "r", "https://api.oxyl.zhacked.me/agent/auth/refresh", "Refresh endpoint")
	agentCmd.Flags().StringVarP(&shutdownEndpoint, "shutdown-endpoint", "s", "https://api.oxyl.zhacked.me/agent/auth/shutdown", "Shutdown endpoint")

	if err := agentCmd.MarkFlagRequired("id"); err != nil {
		slog.Error("failed to mark flag required", "error", err)
		os.Exit(1)
	}
}

func startAgent(cmd *cobra.Command, _ []string) {
	authService, err := service.NewAuthService(agentId, loginEndpoint, refreshEndpoint, shutdownEndpoint)
	if err != nil {
		slog.Error("failed to create authentication service", "error", err)
		os.Exit(1)
	}

	authService.StartTicking(cmd.Context())

	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		slog.Error("failed to load system cert pool", "error", err)
		os.Exit(1)
	}

	grpcClient, err := grpc.NewClient(grpcEndpoint,
		grpc.WithPerRPCCredentials(authService),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})))

	if err != nil {
		slog.Error("failed to create gRPC client", "error", err)
		os.Exit(1)
	}

	defer grpcClient.Close()
	
	// wait for signal
	<-cmd.Context().Done()

	err = authService.RequestShutdown()
	if err != nil {
		slog.Error("failed to request shutdown authentication invalidation", "error", err)
	}

}
