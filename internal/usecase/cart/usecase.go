package cart

import (
	"context"
	"errors"

	cartDTO "go-e-commerce/internal/dto/cart"
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

func (u *cartUseCase) GetCart(ctx context.Context, userID uuid.UUID) (*cartDTO.CartResponse, error) {
	cart, err := u.getOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapCartToResponse(cart), nil
}

func (u *cartUseCase) AddToCart(ctx context.Context, userID uuid.UUID, req *cartDTO.AddCartItemRequest) error {
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

		if product.Stock < req.Quantity {
			return errors.New("insufficient stock")
		}

		existingItem, err := u.cartRepo.GetCartItem(ctx, cart.ID, req.ProductID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if existingItem != nil {
			newQuantity := existingItem.Quantity + req.Quantity
			if product.Stock < newQuantity {
				return errors.New("insufficient stock for combined quantity")
			}
			return u.cartRepo.UpdateCartItemQuantity(ctx, existingItem.ID, newQuantity)
		}

		newItem := &entity.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		return u.cartRepo.AddCartItem(ctx, newItem)
	})
}

func (u *cartUseCase) UpdateCartItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req *cartDTO.UpdateCartItemRequest) error {
	return u.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		cart, err := u.getOrCreateCart(ctx, userID)
		if err != nil {
			return err
		}

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

func (u *cartUseCase) BatchDeleteCartItems(ctx context.Context, userID uuid.UUID, req *cartDTO.BatchDeleteCartItemsRequest) error {
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

func mapCartToResponse(cart *entity.Cart) *cartDTO.CartResponse {
	res := &cartDTO.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		Items:     make([]cartDTO.CartItemResponse, 0, len(cart.Items)),
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}

	for _, item := range cart.Items {
		itemRes := cartDTO.CartItemResponse{
			ID:        item.ID,
			CartID:    item.CartID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}

		if item.Product != nil {
			itemRes.Product = cartDTO.ProductResponse{
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
