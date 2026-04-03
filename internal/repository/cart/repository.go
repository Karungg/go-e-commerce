package cart

import (
	"context"

	"go-e-commerce/internal/entity"
	"go-e-commerce/internal/model"
	"go-e-commerce/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *cartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	var cartModel model.CartModel
	db := repository.ExtractTx(ctx, r.db)

	err := db.WithContext(ctx).
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Where("user_id = ?", userID).
		First(&cartModel).Error

	if err != nil {
		return nil, err
	}

	return mapCartModelToEntity(&cartModel), nil
}

func (r *cartRepository) CreateCart(ctx context.Context, cart *entity.Cart) error {
	cartModel := &model.CartModel{
		ID:     cart.ID,
		UserID: cart.UserID,
	}

	db := repository.ExtractTx(ctx, r.db)
	if err := db.WithContext(ctx).Create(cartModel).Error; err != nil {
		return err
	}

	return nil
}

func (r *cartRepository) GetCartItem(ctx context.Context, cartID, productID uuid.UUID) (*entity.CartItem, error) {
	var itemModel model.CartItemModel
	db := repository.ExtractTx(ctx, r.db)

	err := db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		First(&itemModel).Error

	if err != nil {
		return nil, err
	}

	return mapCartItemModelToEntity(&itemModel), nil
}

func (r *cartRepository) AddCartItem(ctx context.Context, item *entity.CartItem) error {
	itemModel := &model.CartItemModel{
		ID:        item.ID,
		CartID:    item.CartID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
	}

	db := repository.ExtractTx(ctx, r.db)
	return db.WithContext(ctx).Create(itemModel).Error
}

func (r *cartRepository) UpdateCartItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) error {
	db := repository.ExtractTx(ctx, r.db)
	return db.WithContext(ctx).Model(&model.CartItemModel{}).
		Where("id = ?", itemID).
		Update("quantity", quantity).Error
}

func (r *cartRepository) DeleteCartItems(ctx context.Context, cartID uuid.UUID, itemIDs []uuid.UUID) error {
	if len(itemIDs) == 0 {
		return nil
	}
	db := repository.ExtractTx(ctx, r.db)
	return db.WithContext(ctx).Where("cart_id = ? AND id IN ?", cartID, itemIDs).Delete(&model.CartItemModel{}).Error
}

func mapCartModelToEntity(m *model.CartModel) *entity.Cart {
	cart := &entity.Cart{
		ID:        m.ID,
		UserID:    m.UserID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Items:     make([]entity.CartItem, 0, len(m.Items)),
	}

	for _, itemModel := range m.Items {
		item := mapCartItemModelToEntity(&itemModel)
		cart.Items = append(cart.Items, *item)
	}

	return cart
}

func mapCartItemModelToEntity(m *model.CartItemModel) *entity.CartItem {
	item := &entity.CartItem{
		ID:        m.ID,
		CartID:    m.CartID,
		ProductID: m.ProductID,
		Quantity:  m.Quantity,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	if m.Product.ID != uuid.Nil {
		item.Product = &entity.Product{
			ID:          m.Product.ID,
			Title:       m.Product.Title,
			Description: m.Product.Description,
			Price:       m.Product.Price,
			Stock:       m.Product.Stock,
			Image:       m.Product.Image,
			CategoryID:  m.Product.CategoryID,
			SKU:         m.Product.SKU,
			Status:      m.Product.Status,
		}
	}
	return item
}
