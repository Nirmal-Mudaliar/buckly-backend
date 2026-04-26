package config

import (
	core_constants "buckly-ms/core/constants"
	core_utils "buckly-ms/core/utils"
	"buckly-ms/gateway/constants"
)

type GatewayConfig struct {
	ServiceName            string
	HTTPPort               string
	Environment            string
	DatabaseServiceAddress string
	AuthServiceAddress     string
}

func LoadGatewayConfig() *GatewayConfig {
	return &GatewayConfig{
		ServiceName:            core_utils.GetEnv(constants.SERVICE_NAME),
		HTTPPort:               core_utils.GetEnv(constants.HTTP_PORT),
		Environment:            core_utils.GetEnv(constants.ENVIRONMENT),
		DatabaseServiceAddress: core_utils.GetEnv(core_constants.DATABASE_SERVICE_ADDRESS),
		AuthServiceAddress:     core_utils.GetEnv(core_constants.AUTH_SERVICE_ADDRESS),
	}
}
