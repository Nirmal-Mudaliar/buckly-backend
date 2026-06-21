package main

import (
	core_constants "buckly-ms/core/constants"
	"buckly-ms/core/utils"
	bucket_list_gen "buckly-ms/proto/bucket-list-gen"
	"buckly-ms/services/bucket-list-service/config"
	grpcclients "buckly-ms/services/bucket-list-service/grpc-clients"
	"buckly-ms/services/bucket-list-service/handlers"
	"net"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	bucketListServiceConfig := config.LoadBucketListServiceConfig()

	logger := utils.InitLog(bucketListServiceConfig.Environment == core_constants.DEVELOPMENT)

	listener, err := net.Listen(core_constants.NETWORK, ":"+bucketListServiceConfig.GRPCPort)
	if err != nil {
		logger.Fatal("Failed to listen on port: ", zap.String("port", bucketListServiceConfig.GRPCPort), zap.Error(err))
	}

	creds, err := credentials.NewServerTLSFromFile(
		strings.Replace(core_constants.GRPC_CERT, "{service}", bucketListServiceConfig.ServiceName, -1),
		strings.Replace(core_constants.GRPC_CERT_KEY, "{service}", bucketListServiceConfig.ServiceName, -1),
	)
	if err != nil {
		logger.Fatal("Failed to create TLS credentials: ", zap.Error(err))
	}

	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(100),
		grpc.Creds(creds),
	)

	bucketListServiceServer := &handlers.BucketListServiceServer{
		DatabaseServiceClient: grpcclients.NewDatabaseServiceClient(bucketListServiceConfig.DatabaseServiceAddress),
	}

	bucket_list_gen.RegisterBucketListServiceServer(grpcServer, bucketListServiceServer)
	reflection.Register(grpcServer)

	logger.Info("Starting BucketListService", zap.Any("config: ", bucketListServiceConfig))

	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatal("Failed to serve gRPC server: ", zap.Error(err))
	}
}
