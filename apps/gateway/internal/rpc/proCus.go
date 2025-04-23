package rpc

import (
	"context"
	"fmt"

	proto "github.com/FireFly4ik/Lavka-products-customer/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProCusClient struct {
	api proto.ProductsCustomerServiceClient
}

func NewProCus(addr string) (*ProCusClient, error) {
	const op = "proCus.New"

	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, fmt.Errorf("%s: dial: %w", op, err)
	}

	return &ProCusClient{api: proto.NewProductsCustomerServiceClient(cc)}, nil
}

func (c *ProCusClient) ProductRequest(ctx context.Context, productID string) (string, string, *string, float32, float32, string, string, error) {
	const op = "proCus.ProductRequest"

	resp, err := c.api.Product(ctx, &proto.ProductRequest{
		ProductId: productID,
	})
	if err != nil {
		return "", "", nil, 0, 0, "", "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Name, resp.Description, &resp.Image, resp.Price, resp.Discount, resp.Category, resp.Message, nil
}

func (c *ProCusClient) CategoryRequest(ctx context.Context, id string) ([]string, string, error) {
	const op = "proCus.CategoryRequest"

	resp, err := c.api.Category(ctx, &proto.CategoryRequest{
		Id: id,
	})
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.ProductIds, resp.Message, nil
}

func (c *ProCusClient) SearchProductRequest(ctx context.Context, prefix string, category string, minPrice float32, maxPrice float32, hasDiscount bool) ([]string, string, error) {
	const op = "proCus.SearchProductRequest"

	resp, err := c.api.SearchProduct(ctx, &proto.SearchProductRequest{
		Prefix:      prefix,
		Category:    category,
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
		HasDiscount: hasDiscount,
	})
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.ProductIds, resp.Message, nil
}

func (c *ProCusClient) DiscountRequest(ctx context.Context) ([]string, string, error) {
	const op = "proCus.DiscountRequest"

	resp, err := c.api.Discount(ctx, &proto.DiscountRequest{})
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.ProductIds, resp.Message, nil
}

func (c *ProCusClient) GetStockRequest(ctx context.Context, productId string, marketId string) (int32, string, error) {
	const op = "proCus.GetStockRequest"

	resp, err := c.api.GetStock(ctx, &proto.GetStockRequest{
		ProductId: productId,
		MarketId:  marketId,
	})
	if err != nil {
		return 0, "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Stock, resp.Message, nil
}
