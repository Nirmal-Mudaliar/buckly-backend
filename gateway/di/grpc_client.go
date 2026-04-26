package di

import (
	core_constants "buckly-ms/core/constants"
	auth_gen "buckly-ms/proto/auth-gen"
	database_gen "buckly-ms/proto/database-gen"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func newDatabaseServiceClient(databaseServiceAddress string) database_gen.DatabaseServiceClient {
	creds, err := credentials.NewClientTLSFromFile(strings.Replace(core_constants.GRPC_CERT, "{service}", core_constants.DATABASE_SERVICE_NAME, -1), "")
	if err != nil {
		log.Fatalf("Failed to get transport credentials for database service grpc server: %v", err.Error())
	}

	conn, err := grpc.NewClient(
		databaseServiceAddress,
		grpc.WithTransportCredentials(creds),
	)

	if err != nil {
		log.Fatalf("Failed to connect to database service grpc server: %v", err.Error())
	}

	log.Println("Connected to database service grpc server at: ", databaseServiceAddress)
	client := database_gen.NewDatabaseServiceClient(conn)
	return client
}

func newAuthServiceClient(authServiceAddress string) auth_gen.AuthServiceClient {
	creds, err := credentials.NewClientTLSFromFile(strings.Replace(core_constants.GRPC_CERT, "{service}", core_constants.AUTH_SERVICE_NAME, -1), "")
	if err != nil {
		log.Fatalf("Failed to get transport credentials for auth service gRPC server: %v", err.Error())
	}

	conn, err := grpc.NewClient(
		authServiceAddress,
		grpc.WithTransportCredentials(creds),
	)

	if err != nil {
		log.Fatalf("Failed to connect to auth service server at: %v", err.Error())
	}

	log.Println("Connected to auth service gRPC server at: ", authServiceAddress)
	client := auth_gen.NewAuthServiceClient(conn)
	return client
}
