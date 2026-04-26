package main

import (
	core_constants "buckly-ms/core/constants"
	"buckly-ms/core/utils"
	config "buckly-ms/services/database-service/config"
	db "buckly-ms/services/database-service/db/generated"
	"buckly-ms/services/database-service/handlers"
	"context"
	"net"
	"strings"
	"time"

	pb "buckly-ms/proto/database-gen"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	config := config.LoadDatabaseServiceConfig()

	logger := utils.InitLog(config.Environment == core_constants.DEVELOPMENT)

	// Create a new database connection pool
	poolConfig, err := pgxpool.ParseConfig(config.DatabaseURL)

	if err != nil {
		logger.Fatal("Failed to parse database URL", zap.Error(err))
	}

	poolConfig.MaxConns = config.PoolMaxConnections
	poolConfig.MinConns = config.PoolMinConnections
	poolConfig.MaxConnIdleTime = 30 * time.Second
	poolConfig.HealthCheckPeriod = 30 * time.Second
	poolConfig.ConnConfig.Tracer = &BucklyQueryTracer{}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Fatal("Failed to create database connection pool", zap.Error(err))
	}
	defer pool.Close()

	dbServer := &handlers.DatabaseServiceServer{
		Queries: db.New(pool),
		Pool:    pool,
	}

	creds, err := credentials.NewServerTLSFromFile(strings.Replace(core_constants.GRPC_CERT, "{service}", config.ServiceName, -1), strings.Replace(core_constants.GRPC_CERT_KEY, "{service}", config.ServiceName, -1))

	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(100),
		grpc.Creds(creds),
	)

	listner, err := net.Listen(core_constants.NETWORK, ":"+config.GRPCPort)
	if err != nil {
		logger.Fatal("Failed to listen on port "+config.GRPCPort, zap.Error(err))
	}

	pb.RegisterDatabaseServiceServer(grpcServer, dbServer)
	reflection.Register(grpcServer)

	logger.Info("Starting Database Service", zap.Any("config: ", config))

	if err := grpcServer.Serve(listner); err != nil {
		logger.Fatal("Failed to serve gRPC server", zap.Error(err))
	}
}

type BucklyQueryTracer struct{ log *zap.Logger }

func (tracer *BucklyQueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	tracer.log = utils.GetLoggerFromContext(ctx)
	tracer.log = tracer.log.With(zap.String("SQL", data.SQL), zap.Any("Args", data.Args))
	// Store the logger in context
	ctx = utils.SetLoggerInContext(ctx, tracer.log)
	return ctx
}

func (tracer *BucklyQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	// Retrieve the logger from the context
	tracer.log = utils.GetLoggerFromContext(ctx)
	tracer.log = tracer.log.With(zap.String("Operation", data.CommandTag.String()))
	if data.Err != nil {
		tracer.log = tracer.log.With(zap.String("Error", data.Err.Error()))
	}
	tracer.log.Info("Command Executed")
}
