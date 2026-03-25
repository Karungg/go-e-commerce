package category

import (
	"context"

	"go-e-commerce/internal/entity"

	"github.com/google/uuid"

	categoryDTO "go-e-commerce/internal/dto/category"
)

// CategoryRepository defines the contract for category data access
type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	FindAll(ctx context.Context) ([]*entity.Category, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// CategoryUseCase defines the contract for category business logic
type CategoryUseCase interface {
	CreateCategory(ctx context.Context, req *categoryDTO.CreateCategoryReq) (*categoryDTO.CategoryRes, error)
	GetAllCategories(ctx context.Context) ([]*categoryDTO.CategoryRes, error)
	GetCategoryByID(ctx context.Context, id string) (*categoryDTO.CategoryRes, error)
	UpdateCategory(ctx context.Context, id string, req *categoryDTO.UpdateCategoryReq) (*categoryDTO.CategoryRes, error)
	DeleteCategory(ctx context.Context, id string) error
}
