package product

import (
	"context"
	"errors"

	productDTO "go-e-commerce/internal/dto/product"
	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/pkg/apperror"
	productPort "go-e-commerce/internal/port/product"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productUseCase struct {
	productRepo productPort.ProductRepository
}

func NewProductUseCase(productRepo productPort.ProductRepository) productPort.ProductUseCase {
	return &productUseCase{
		productRepo: productRepo,
	}
}

func (u *productUseCase) CreateProduct(ctx context.Context, req *productDTO.CreateProductReq) (*productDTO.ProductRes, error) {
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	product := &entity.Product{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Image:       req.Image,
		CategoryID:  categoryID,
		SKU:         req.SKU,
		Status:      "active", // default status
	}

	if err := u.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return &productDTO.ProductRes{
		ID:          product.ID,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Image:       product.Image,
		CategoryID:  product.CategoryID,
		SKU:         product.SKU,
		Status:      product.Status,
	}, nil
}

func (u *productUseCase) GetAllProducts(ctx context.Context) ([]*productDTO.ProductRes, error) {
	products, err := u.productRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var res []*productDTO.ProductRes
	for _, p := range products {
		res = append(res, &productDTO.ProductRes{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			Image:       p.Image,
			CategoryID:  p.CategoryID,
			SKU:         p.SKU,
			Status:      p.Status,
		})
	}

	if res == nil {
		res = make([]*productDTO.ProductRes, 0)
	}

	return res, nil
}

func (u *productUseCase) GetProductByID(ctx context.Context, id string) (*productDTO.ProductRes, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	product, err := u.productRepo.FindByID(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrProductNotFound
		}
		return nil, err
	}

	return &productDTO.ProductRes{
		ID:          product.ID,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Image:       product.Image,
		CategoryID:  product.CategoryID,
		SKU:         product.SKU,
		Status:      product.Status,
	}, nil
}

func (u *productUseCase) UpdateProduct(ctx context.Context, id string, req *productDTO.UpdateProductReq) (*productDTO.ProductRes, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	product, err := u.productRepo.FindByID(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrProductNotFound
		}
		return nil, err
	}

	// Update fields
	product.Title = req.Title
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.Image = req.Image
	product.CategoryID = categoryID
	product.SKU = req.SKU
	product.Status = req.Status

	if err := u.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return &productDTO.ProductRes{
		ID:          product.ID,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Image:       product.Image,
		CategoryID:  product.CategoryID,
		SKU:         product.SKU,
		Status:      product.Status,
	}, nil
}

func (u *productUseCase) DeleteProduct(ctx context.Context, id string) error {
	productID, err := uuid.Parse(id)
	if err != nil {
		return apperror.ErrBadRequest
	}

	_, err = u.productRepo.FindByID(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrProductNotFound
		}
		return err
	}

	return u.productRepo.Delete(ctx, productID)
}
