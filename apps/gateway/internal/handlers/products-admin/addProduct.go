package products_admin

import (
	"github.com/FireFly4ik/Lavka-gateway/internal/rpc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type AddProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Image_url   string   `json:"image_url"`
	Price       float32  `json:"price"`
	Discount    float32  `json:"discount"`
	Category    []string `json:"category"`
}

type AddProductResponse struct {
	ProductId string `json:"product_id"`
	Message   string `json:"message"`
}

func AddProduct(c *gin.Context) {
	var req AddProductRequest
	var res AddProductResponse
	proAdmClient := c.MustGet("proAdmClient").(*rpc.ProAdmClient)
	logger := c.MustGet("logger").(*zap.Logger)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		logger.Error("AddProduct: invalid request", zap.Error(err))
		return
	}
	logger.Info("AddProduct", zap.Any("request", req))
	var err error
	res.ProductId, res.Message, err = proAdmClient.AddProductRequest(c, req.Name, req.Description, req.Image_url, req.Price, req.Discount, req.Category)
	if err != nil {
		if res.Message == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "no connection with server"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": res.Message})
		}
		logger.Error("AddProduct: response error:", zap.String("message", res.Message), zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"Product_id": res.ProductId, "message": res.Message})
	logger.Info("AddProduct: response:", zap.String("message", res.Message))
}
