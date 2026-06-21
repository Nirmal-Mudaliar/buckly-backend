package handlers

import (
	bucket_list_gen "buckly-ms/proto/bucket-list-gen"
	database_gen "buckly-ms/proto/database-gen"
)

type BucketListServiceServer struct {
	bucket_list_gen.UnimplementedBucketListServiceServer
	DatabaseServiceClient database_gen.DatabaseServiceClient
}
