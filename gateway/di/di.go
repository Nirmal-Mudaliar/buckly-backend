package di

import (
	"buckly-ms/gateway/api/handlers/auth"
	bucket_list "buckly-ms/gateway/api/handlers/bucket-list"
	"buckly-ms/gateway/config"
	"buckly-ms/gateway/contracts"
	auth_gen "buckly-ms/proto/auth-gen"
	bucket_list_gen "buckly-ms/proto/bucket-list-gen"
	database_gen "buckly-ms/proto/database-gen"

	"go.uber.org/dig"
)

const (
	G_HANDLERS = "handlers"
)

type Params struct {
	dig.In
	Config  *config.GatewayConfig
	Handler []contracts.RouteRegistrar `group:"handlers"`
}

func BuildContainer() *dig.Container {
	container := dig.New()
	container.Provide(func() *config.GatewayConfig {
		return config.LoadGatewayConfig()
	})
	container.Provide(func(cfg *config.GatewayConfig) database_gen.DatabaseServiceClient {
		return newDatabaseServiceClient(cfg.DatabaseServiceAddress)
	})
	container.Provide(func(cfg *config.GatewayConfig) auth_gen.AuthServiceClient {
		return newAuthServiceClient(cfg.AuthServiceAddress)
	})
	container.Provide(func(cfg *config.GatewayConfig) bucket_list_gen.BucketListServiceClient {
		return newBucketListServiceClient(cfg.BucketListServiceAddress)
	})
	container.Provide(auth.NewAuthHandler, dig.Group(G_HANDLERS), dig.As(new(contracts.RouteRegistrar)))
	container.Provide(bucket_list.NewBucketListHandler, dig.Group(G_HANDLERS), dig.As(new(contracts.RouteRegistrar)))

	return container
}
