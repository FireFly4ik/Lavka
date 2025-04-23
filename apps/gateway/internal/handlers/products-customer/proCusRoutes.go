package products_customer

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterProCusRoutes(r *gin.Engine, log *zap.Logger, proCusClient *rpc.ProCusClient) {
	authGroup := r.Group("products-customer")
	authGroup.Use(func(c *gin.Context) {
		c.Set("proCusClient", proCusClient)
		c.Set("logger", log)
		c.Next()
	})
	{
		authGroup.GET("/catalog", Catalog)
		authGroup.GET("/product", Product)
		authGroup.GET("/searchproduct", SearchProduct)
		authGroup.GET("/discounts", Discount)
		authGroup.GET("/getstock", GetStock)
	}
	log.Info("products-customer routes were registered")
}
