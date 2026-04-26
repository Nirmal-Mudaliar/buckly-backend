package handlers

import (
	auth_gen "buckly-ms/proto/auth-gen"
	database_gen "buckly-ms/proto/database-gen"
	"buckly-ms/services/auth-service/twilio"
)

type AuthServiceServer struct {
	auth_gen.UnimplementedAuthServiceServer
	DatabaseServiceClient database_gen.DatabaseServiceClient
	TwilioClient          *twilio.TwilioClient
}
