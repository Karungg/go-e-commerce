package dto

import "github.com/google/uuid"

type CreateCategoryReq struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateCategoryReq struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type CategoryRes struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}
