package config

import (
	core_constants "buckly-ms/core/constants"
	core_utils "buckly-ms/core/utils"
)

type AuthServiceConfig struct {
	ServiceName            string
	GRPCPort               string
	DatabaseServiceAddress string
	Environment            string
	TwilioAccountSID       string
	TwilioAuthToken        string
	TwilioVerifyServiceSID string
}

func LoadAuthServiceConfig() *AuthServiceConfig {
	return &AuthServiceConfig{
		ServiceName:            core_utils.GetEnv(core_constants.SERVICE_NAME),
		GRPCPort:               core_utils.GetEnv(core_constants.GRPC_PORT),
		DatabaseServiceAddress: core_utils.GetEnv(core_constants.DATABASE_SERVICE_ADDRESS),
		Environment:            core_utils.GetEnv(core_constants.ENVIRONMENT),
		TwilioAccountSID:       core_utils.GetEnv(core_constants.TWILIO_ACCOUNT_SID),
		TwilioAuthToken:        core_utils.GetEnv(core_constants.TWILIO_AUTH_TOKEN),
		TwilioVerifyServiceSID: core_utils.GetEnv(core_constants.TWILIO_VERIFY_SERVICE_SID),
	}
}
