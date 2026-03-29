package product

import (
	"context"

	"go-e-commerce/internal/entity"

	"github.com/google/uuid"

	productDTO "go-e-commerce/internal/dto/product"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	FindAll(ctx context.Context, limit, offset int) ([]*entity.Product, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProductUseCase interface {
	CreateProduct(ctx context.Context, req *productDTO.CreateProductReq) (*productDTO.ProductRes, error)
	GetAllProducts(ctx context.Context, page, limit int) ([]*productDTO.ProductRes, int64, error)
	GetProductByID(ctx context.Context, id string) (*productDTO.ProductRes, error)
	UpdateProduct(ctx context.Context, id string, req *productDTO.UpdateProductReq) (*productDTO.ProductRes, error)
	DeleteProduct(ctx context.Context, id string) error
}
