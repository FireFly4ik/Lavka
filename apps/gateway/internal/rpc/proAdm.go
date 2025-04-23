package rpc

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"

	proto "github.com/FireFly4ik/Lavka-products-admin/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProAdmClient struct {
	api proto.ProductsAdminServiceClient
}

func NewProAdm(addr string) (*ProAdmClient, error) {
	const op = "proAdm.New"

	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, fmt.Errorf("%s: dial: %w", op, err)
	}

	return &ProAdmClient{api: proto.NewProductsAdminServiceClient(cc)}, nil
}

func (c *ProAdmClient) AddProductRequest(ctx context.Context, name string, description string, image string, price float32, discount float32, category []string) (string, string, error) {
	const op = "proAdm.AddProductRequest"

	resp, err := c.api.AddProduct(ctx, &proto.AddProductRequest{
		Name:        name,
		Description: description,
		Image:       image,
		Price:       price,
		Discount:    discount,
		Category:    category,
	})
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.ProductId, resp.Message, nil
}

func (c *ProAdmClient) UpdateProductRequest(ctx *gin.Context, id string, name string, description string, image string, price float32, discount float32, category []string) (string, error) {
	const op = "proAdm.AddProductRequest"

	resp, err := c.api.UpdateProduct(ctx, &proto.UpdateProductRequest{
		ProductId:   id,
		Name:        name,
		Description: description,
		Image:       image,
		Price:       price,
		Discount:    discount,
		Category:    category,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Message, nil
}

func (c *ProAdmClient) DeleteProductRequest(ctx *gin.Context, id string) (string, error) {
	const op = "proAdm.AddProductRequest"

	resp, err := c.api.DeleteProduct(ctx, &proto.DeleteProductRequest{
		ProductId: id,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Message, nil
}

func (c *ProAdmClient) ApplyDiscountRequest(ctx *gin.Context, id string, discount float32) (string, error) {
	const op = "proAdm.ApplyDiscountRequest"

	resp, err := c.api.ApplyDiscount(ctx, &proto.ApplyDiscountRequest{
		ProductId: id,
		Discount:  discount,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Message, nil
}

func (c *ProAdmClient) RemoveDiscountRequest(ctx *gin.Context, id string) (string, error) {
	const op = "proAdm.RemoveDiscountRequest"

	resp, err := c.api.RemoveDiscount(ctx, &proto.RemoveDiscountRequest{
		ProductId: id,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Message, nil
}
