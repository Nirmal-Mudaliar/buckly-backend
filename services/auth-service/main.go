package main

import (
	core_constants "buckly-ms/core/constants"
	"buckly-ms/core/utils"
	"buckly-ms/services/auth-service/config"
	grpcclients "buckly-ms/services/auth-service/grpc-clients"
	"buckly-ms/services/auth-service/handlers"
	"buckly-ms/services/auth-service/twilio"
	"net"
	"strings"

	auth_gen "buckly-ms/proto/auth-gen"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	authConfig := config.LoadAuthServiceConfig()

	logger := utils.InitLog(authConfig.Environment == core_constants.DEVELOPMENT)

	listner, err := net.Listen(core_constants.NETWORK, ":"+authConfig.GRPCPort)
	if err != nil {
		logger.Fatal("Failed to listen on port: ", zap.String("port", authConfig.GRPCPort), zap.Error(err))
	}

	creds, err := credentials.NewServerTLSFromFile(strings.Replace(core_constants.GRPC_CERT, "{service}", authConfig.ServiceName, -1), strings.Replace(core_constants.GRPC_CERT_KEY, "{service}", authConfig.ServiceName, -1))
	if err != nil {
		logger.Fatal("Failed to load TLS credentials: ", zap.Error(err))
	}

	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(100),
		grpc.Creds(creds),
	)

	authServiceServer := &handlers.AuthServiceServer{
		DatabaseServiceClient: grpcclients.NewDatabaseServiceClient(authConfig.DatabaseServiceAddress),
		TwilioClient: twilio.NewTwilioClient(
			authConfig.TwilioAccountSID,
			authConfig.TwilioAuthToken,
			authConfig.TwilioVerifyServiceSID,
		),
	}

	auth_gen.RegisterAuthServiceServer(grpcServer, authServiceServer)
	reflection.Register(grpcServer)

	logger.Info("Starting AuthService", zap.Any("config: ", authConfig))

	if err := grpcServer.Serve(listner); err != nil {
		logger.Fatal("Failed to serve gRPC server: ", zap.Error(err))
	}
}
