package cart

import (
	"context"
	"errors"

	"go-e-commerce/internal/dto"
	"go-e-commerce/internal/entity"
	authPort "go-e-commerce/internal/port/auth"
	cartPort "go-e-commerce/internal/port/cart"
	productPort "go-e-commerce/internal/port/product"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type cartUseCase struct {
	txManager   authPort.TransactionManager
	cartRepo    cartPort.CartRepository
	productRepo productPort.ProductRepository
}

func NewCartUseCase(
	txManager authPort.TransactionManager,
	cartRepo cartPort.CartRepository,
	productRepo productPort.ProductRepository,
) cartPort.CartUseCase {
	return &cartUseCase{
		txManager:   txManager,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (u *cartUseCase) GetCart(ctx context.Context, userID uuid.UUID) (*dto.CartResponse, error) {
	cart, err := u.getOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapCartToResponse(cart), nil
}

func (u *cartUseCase) AddToCart(ctx context.Context, userID uuid.UUID, req *dto.AddCartItemRequest) error {
	return u.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		cart, err := u.getOrCreateCart(ctx, userID)
		if err != nil {
			return err
		}

		product, err := u.productRepo.FindByID(ctx, req.ProductID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("product not found")
			}
			return err
		}

		// Check stock
		if product.Stock < req.Quantity {
			return errors.New("insufficient stock")
		}

		// Check if item already in cart
		existingItem, err := u.cartRepo.GetCartItem(ctx, cart.ID, req.ProductID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if existingItem != nil {
			// Update quantity
			newQuantity := existingItem.Quantity + req.Quantity
			if product.Stock < newQuantity {
				return errors.New("insufficient stock for combined quantity")
			}
			return u.cartRepo.UpdateCartItemQuantity(ctx, existingItem.ID, newQuantity)
		}

		// Create new item
		newItem := &entity.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		return u.cartRepo.AddCartItem(ctx, newItem)
	})
}

func (u *cartUseCase) UpdateCartItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req *dto.UpdateCartItemRequest) error {
	return u.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		// Verify cart belongs to user
		cart, err := u.getOrCreateCart(ctx, userID)
		if err != nil {
			return err
		}

		// Since we don't have GetCartItemByID directly, we just verify it exists within the user's cart
		var targetItem *entity.CartItem
		for _, item := range cart.Items {
			if item.ID == itemID {
				targetItem = &item
				break
			}
		}

		if targetItem == nil {
			return errors.New("cart item not found in your cart")
		}

		product, err := u.productRepo.FindByID(ctx, targetItem.ProductID)
		if err != nil {
			return err
		}

		if product.Stock < req.Quantity {
			return errors.New("insufficient stock")
		}

		return u.cartRepo.UpdateCartItemQuantity(ctx, itemID, req.Quantity)
	})
}

func (u *cartUseCase) BatchDeleteCartItems(ctx context.Context, userID uuid.UUID, req *dto.BatchDeleteCartItemsRequest) error {
	return u.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		cart, err := u.getOrCreateCart(ctx, userID)
		if err != nil {
			return err
		}

		return u.cartRepo.DeleteCartItems(ctx, cart.ID, req.CartItemIDs)
	})
}

func (u *cartUseCase) getOrCreateCart(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	cart, err := u.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new cart
			newCart := &entity.Cart{
				UserID: userID,
			}
			if err := u.cartRepo.CreateCart(ctx, newCart); err != nil {
				return nil, err
			}
			return newCart, nil
		}
		return nil, err
	}
	return cart, nil
}

func mapCartToResponse(cart *entity.Cart) *dto.CartResponse {
	res := &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		Items:     make([]dto.CartItemResponse, 0, len(cart.Items)),
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}

	for _, item := range cart.Items {
		itemRes := dto.CartItemResponse{
			ID:        item.ID,
			CartID:    item.CartID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}

		if item.Product != nil {
			itemRes.Product = dto.ProductResponse{
				ID:          item.Product.ID,
				Title:       item.Product.Title,
				Description: item.Product.Description,
				Price:       item.Product.Price,
				Stock:       item.Product.Stock,
				Image:       item.Product.Image,
			}
			itemRes.Subtotal = item.Product.Price * float64(item.Quantity)
			res.TotalPrice += itemRes.Subtotal
		}

		res.Items = append(res.Items, itemRes)
	}

	return res
}
