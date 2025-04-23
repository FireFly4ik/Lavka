package products_admin

import (
	"context"
	"github.com/FireFly4ik/Lavka-products-admin/internal/db"
	"github.com/FireFly4ik/Lavka-products-admin/internal/models"
	pb "github.com/FireFly4ik/Lavka-products-admin/proto"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"log"
)

type ProductsAdminServiceServer struct {
	pb.UnsafeProductsAdminServiceServer
	Logger *zap.Logger
}

func (s *ProductsAdminServiceServer) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.AddProductResponse, error) {
	var product models.Product

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Discount = req.Discount
	product.Image_url = &req.Image

	tx := db.DB.Begin()
	if tx.Error != nil {
		return &pb.AddProductResponse{Message: "Couldn't start transaction"}, tx.Error
	}

	if err := tx.Create(&product).Error; err != nil {
		return &pb.AddProductResponse{Message: "Couldn't create product"}, err
	}

	//TODO добавить продукт во все магазины

	for _, el := range req.Category {
		var newCategory models.Category_product

		newCatId, _ := uuid.Parse(el)
		newCategory.CategoryID = newCatId
		newCategory.ProductID = product.ID

		if err := tx.Create(&newCategory).Error; err != nil {
			return &pb.AddProductResponse{Message: "Couldn't create categories for product"}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return &pb.AddProductResponse{Message: "Could not commit transaction"}, err
	}

	return &pb.AddProductResponse{
		Message:   "Product added successfully",
		ProductId: product.ID.String(),
	}, nil
}

func (s *ProductsAdminServiceServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {

	tx := db.DB.Begin()

	if req.Name != "" {
		if err := tx.Model(&models.Product{}).Where("id = ?", req.ProductId).Update("name", req.Name).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if req.Description != "" {
		if err := tx.Model(&models.Product{}).Where("id = ?", req.ProductId).Update("description", req.Description).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if req.Price != -1 {
		if err := tx.Model(&models.Product{}).Where("id = ?", req.ProductId).Update("price", req.Price).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if req.Discount != -1 {
		if err := tx.Model(&models.Product{}).Where("id = ?", req.ProductId).Update("discount", req.Discount).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if req.Image != "" {
		if err := tx.Model(&models.Product{}).Where("id = ?", req.ProductId).Update("image_url", &req.Image).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if len(req.Category) != 0 {
		var productCategories []models.Category_product

		query := tx.Model(&models.Category_product{}).Where("product_id = ?", req.ProductId)

		if err := query.Find(&productCategories).Error; err != nil {
			return nil, err
		}

		seenNew := make(map[string]struct{})
		seenOld := make(map[string]struct{})
		common := make(map[string]struct{})

		for _, v := range req.Category {
			seenNew[v] = struct{}{}
		}
		for _, v := range productCategories {
			seenOld[v.CategoryID.String()] = struct{}{}
			if _, ok := seenNew[v.CategoryID.String()]; ok {
				common[v.CategoryID.String()] = struct{}{}
			}
		}

		toAdd := []string{}
		for _, v := range req.Category {
			if _, isCommon := common[v]; !isCommon {
				toAdd = append(toAdd, v)
			}
		}

		toDelete := []string{}
		for _, v := range productCategories {
			if _, isCommon := common[v.CategoryID.String()]; !isCommon {
				toDelete = append(toDelete, v.CategoryID.String())
			}
		}

		prodID, _ := uuid.Parse(req.ProductId)

		for _, el := range toAdd {
			log.Println(len(toAdd))
			var newCategory models.Category_product
			newCatID, _ := uuid.Parse(el)
			newCategory.CategoryID = newCatID
			newCategory.ProductID = prodID
			if err := tx.Create(&newCategory).Error; err != nil {
				return &pb.UpdateProductResponse{Message: "Couldn't create new category"}, err
			}
		}

		for _, el := range toDelete {
			delCatID, _ := uuid.Parse(el)
			if err := tx.Where("product_id = ? and category_id = ?", prodID, delCatID).Delete(&models.Category_product{}).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return &pb.UpdateProductResponse{Message: "Couldn't commit transaction"}, err
	}

	return &pb.UpdateProductResponse{Message: "Product updated successfully"}, nil
}

func (s *ProductsAdminServiceServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	prodId, _ := uuid.Parse(req.ProductId)

	tx := db.DB.Begin()

	if err := tx.Where("product_id = ?", prodId).Delete(&models.Category_product{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Where("product_id = ?", prodId).Delete(&models.Market_product{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Where("id = ?", prodId).Delete(&models.Product{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return &pb.DeleteProductResponse{Message: "Couldn't commit transaction"}, err
	}

	return nil, nil
}

func (s *ProductsAdminServiceServer) ApplyDiscount(ctx context.Context, req *pb.ApplyDiscountRequest) (*pb.ApplyDiscountResponse, error) {
	productId, _ := uuid.Parse(req.ProductId)

	tx := db.DB.Begin()

	if err := tx.Model(&models.Product{}).Where("id = ?", productId).Update("discount", req.Discount).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return &pb.ApplyDiscountResponse{Message: "Couldn't commit transaction"}, err
	}

	return &pb.ApplyDiscountResponse{Message: "Discount has been applied successfully"}, nil
}

func (s *ProductsAdminServiceServer) RemoveDiscount(ctx context.Context, req *pb.RemoveDiscountRequest) (*pb.RemoveDiscountResponse, error) {
	productId, _ := uuid.Parse(req.ProductId)

	tx := db.DB.Begin()

	if err := tx.Model(&models.Product{}).Where("id = ?", productId).Update("discount", 0).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return &pb.RemoveDiscountResponse{Message: "Couldn't commit transaction"}, err
	}

	return &pb.RemoveDiscountResponse{Message: "Discount has been deleted successfully"}, nil

}

func (s *ProductsAdminServiceServer) AddStock(ctx context.Context, req *pb.AddStockRequest) (*pb.AddStockResponse, error) {

	//TODO позже)

	return nil, nil
}
