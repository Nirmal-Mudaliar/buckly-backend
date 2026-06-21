package config

import (
	core_constants "buckly-ms/core/constants"
	core_utils "buckly-ms/core/utils"
)

type BucketListServiceConfig struct {
	ServiceName            string
	GRPCPort               string
	DatabaseServiceAddress string
	Environment            string
}

func LoadBucketListServiceConfig() *BucketListServiceConfig {
	return &BucketListServiceConfig{
		ServiceName:            core_utils.GetEnv(core_constants.SERVICE_NAME),
		GRPCPort:               core_utils.GetEnv(core_constants.GRPC_PORT),
		DatabaseServiceAddress: core_utils.GetEnv(core_constants.DATABASE_SERVICE_ADDRESS),
		Environment:            core_utils.GetEnv(core_constants.ENVIRONMENT),
	}
}
