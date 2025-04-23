package products_customer

import (
	"context"
	"errors"
	database "github.com/FireFly4ik/Lavka-products-customer/internal/db"
	"github.com/FireFly4ik/Lavka-products-customer/internal/models"
	pb "github.com/FireFly4ik/Lavka-products-customer/proto"
	"go.uber.org/zap"
)

type ProductsCustomerServiceServer struct {
	pb.UnsafeProductsCustomerServiceServer
	Logger *zap.Logger
}

func (s *ProductsCustomerServiceServer) Category(ctx context.Context, req *pb.CategoryRequest) (*pb.CategoryResponse, error) {
	var productsInCategory []models.Category_product

	if err := database.DB.Where("category_id = ?", req.Id).Find(&productsInCategory).Error; err != nil {
		return nil, err
	}

	var productIDs []string
	for _, product := range productsInCategory {
		productIDs = append(productIDs, product.ProductID.String())
	}

	return &pb.CategoryResponse{
		ProductIds: productIDs,
		Message:    "Successful",
	}, nil
}

func (s *ProductsCustomerServiceServer) Product(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {
	var product models.Product

	if err := database.DB.Where("product_id = ?", req.ProductId).First(&product).Error; err != nil {
		return nil, errors.New("invalid id")
	}

	var Image_url string
	if product.Image_url == nil {
		Image_url = ""
	} else {
		Image_url = *product.Image_url
	}

	return &pb.ProductResponse{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Discount:    product.Discount,
		Image:       Image_url,
	}, nil
}

func (s *ProductsCustomerServiceServer) SearchProduct(ctx context.Context, req *pb.SearchProductRequest) (*pb.SearchProductResponse, error) {
	var productsInCategory []models.Category_product

	// Построение запроса с учетом фильтров
	query := database.DB.Model(&models.Category_product{}).Joins("JOIN products ON products.id = category_products.product_id")

	if req.Prefix != "" {
		query = query.Where("products.name ILIKE ?", "%"+req.Prefix+"%")
	}

	if req.Category != "" {
		query = query.Where("category.Id = ?", req.Category)
	}
	if req.MinPrice > 0 {
		query = query.Where("products.price >= ?", req.MinPrice)
	}
	if req.MaxPrice > 0 {
		query = query.Where("products.price <= ?", req.MaxPrice)
	}
	if req.HasDiscount {
		query = query.Where("products.discount > 0")
	}

	if err := query.Find(&productsInCategory).Error; err != nil {
		return nil, err
	}

	var productIDs []string
	seen := make(map[string]struct{})
	for _, product := range productsInCategory {
		if _, ok := seen[product.ProductID.String()]; !ok {
			seen[product.ProductID.String()] = struct{}{}
			productIDs = append(productIDs, product.ProductID.String())
		}
	}

	return &pb.SearchProductResponse{
		ProductIds: productIDs,
		Message:    "Successful",
	}, nil
}

func (s *ProductsCustomerServiceServer) Discount(ctx context.Context, req *pb.DiscountRequest) (*pb.DiscountResponse, error) {
	var productsInCategory []models.Product

	if err := database.DB.Where("discount > 0").Find(&productsInCategory).Error; err != nil {
		return nil, err
	}

	var response []string

	for _, product := range productsInCategory {
		response = append(response, product.ID.String())
	}

	return &pb.DiscountResponse{
		ProductIds: response,
		Message:    "Successful",
	}, nil
}

func (s *ProductsCustomerServiceServer) GetStock(ctx context.Context, req *pb.GetStockRequest) (*pb.GetStockResponse, error) {
	var prodMarket models.Market_product

	if err := database.DB.Where("product_id = ? and market_id = ?", req.ProductId, req.MarketId).First(&prodMarket).Error; err != nil {
		return nil, errors.New("invalid id")
	}

	return &pb.GetStockResponse{
		Stock:   int32(prodMarket.Stock),
		Message: "Successfully",
	}, nil
}
