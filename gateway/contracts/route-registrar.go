package contracts

import "github.com/gin-gonic/gin"

type RouteRegistrar interface {
	RegisterRoutes(r *gin.Engine)
}
