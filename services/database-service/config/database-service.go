package database_service

import (
	core_utils "buckly-ms/core/utils"
	database_service "buckly-ms/services/database-service/constants"
)

type DatabaseServiceConfig struct {
	ServiceName        string
	GRPCPort           string
	DatabaseURL        string
	Environment        string
	PoolMaxConnections int32
	PoolMinConnections int32
}

func LoadDatabaseServiceConfig() *DatabaseServiceConfig {
	return &DatabaseServiceConfig{
		ServiceName:        core_utils.GetEnv(database_service.SERVICE_NAME),
		GRPCPort:           core_utils.GetEnv(database_service.GRPC_PORT),
		DatabaseURL:        core_utils.GetEnv(database_service.DATABASE_URL),
		Environment:        core_utils.GetEnv(database_service.ENVIRONMENT),
		PoolMaxConnections: core_utils.GetInt32Env(database_service.MAX_POOL_SIZE, 30),
	}
}
