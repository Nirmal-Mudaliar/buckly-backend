package handlers

import (
	pb "buckly-ms/proto/database-gen"
	db "buckly-ms/services/database-service/db/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseServiceServer struct {
	pb.UnimplementedDatabaseServiceServer
	Queries *db.Queries
	Pool    *pgxpool.Pool
}
