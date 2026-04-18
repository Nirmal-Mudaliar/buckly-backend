package di

import (
	"buckly-ms/gateway/api/handlers/auth"
	"buckly-ms/gateway/config"
	"buckly-ms/gateway/contracts"

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
	container.Provide(auth.NewAuthHandler, dig.Group(G_HANDLERS), dig.As(new(contracts.RouteRegistrar)))
	return container
}
