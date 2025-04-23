package products_admin

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterProAdmRoutes(r *gin.Engine, log *zap.Logger, proAdmClient *rpc.ProAdmClient) {
	authGroup := r.Group("products-admin")
	authGroup.Use(func(c *gin.Context) {
		c.Set("proAdmClient", proAdmClient)
		c.Set("logger", log)
		c.Next()
	})
	{
		authGroup.POST("/addproduct", AddProduct)
		authGroup.PUT("/updateproduct", UpdateProduct)
		authGroup.DELETE("/deleteproduct", DeleteProduct)
		authGroup.POST("/applydiscount", ApplyDiscount)
		authGroup.DELETE("/removediscount", RemoveDiscount)
		authGroup.PUT("/bulkupdateproducts", BulkUpdateProducts)
	}
	log.Info("products-admin routes were registered")
}
