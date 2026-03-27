package product

import "github.com/google/uuid"

type CreateProductReq struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	Image       string  `json:"image" binding:"required"`
	CategoryID  string  `json:"category_id" binding:"required,uuid"`
	SKU         string  `json:"sku" binding:"required"`
}

type UpdateProductReq struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	Image       string  `json:"image" binding:"required"`
	CategoryID  string  `json:"category_id" binding:"required,uuid"`
	SKU         string  `json:"sku" binding:"required"`
	Status      string  `json:"status" binding:"required,oneof=active inactive"`
}

type ProductRes struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Image       string    `json:"image"`
	CategoryID  uuid.UUID `json:"category_id"`
	SKU         string    `json:"sku"`
	Status      string    `json:"status"`
}
