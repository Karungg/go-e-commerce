package cart_test

import (
	"context"
	"testing"

	cartDTO "go-e-commerce/internal/dto/cart"
	"go-e-commerce/internal/entity"
	authMock "go-e-commerce/internal/mocks/auth"
	cartMock "go-e-commerce/internal/mocks/cart"
	productMock "go-e-commerce/internal/mocks/product"
	cartPort "go-e-commerce/internal/port/cart"
	"go-e-commerce/internal/usecase/cart"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	mockUserID    = uuid.New()
	mockCartID    = uuid.New()
	mockProductID = uuid.New()
	mockItemID    = uuid.New()
)

func setupTest() (*authMock.TransactionManagerMock, *cartMock.CartRepositoryMock, *productMock.ProductRepositoryMock, cartPort.CartUseCase) {
	txManager := &authMock.TransactionManagerMock{}
	cartRepo := &cartMock.CartRepositoryMock{}
	productRepo := &productMock.ProductRepositoryMock{}
	
	// Assuming cart.NewCartUseCase returns the cartPort.CartUseCase interface
	useCase := cart.NewCartUseCase(txManager, cartRepo, productRepo)
	
	return txManager, cartRepo, productRepo, useCase
}

func TestGetCart_Success(t *testing.T) {
	_, cartRepo, _, useCase := setupTest()

	expectedCart := &entity.Cart{
		ID:     mockCartID,
		UserID: mockUserID,
		Items:  []entity.CartItem{},
	}

	cartRepo.On("GetCartByUserID", mock.Anything, mockUserID).Return(expectedCart, nil)

	res, err := useCase.GetCart(context.Background(), mockUserID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, mockCartID, res.ID)
	cartRepo.AssertExpectations(t)
}

func TestGetCart_AutoCreate(t *testing.T) {
	_, cartRepo, _, useCase := setupTest()

	cartRepo.On("GetCartByUserID", mock.Anything, mockUserID).Return(nil, gorm.ErrRecordNotFound)
	cartRepo.On("CreateCart", mock.Anything, mock.AnythingOfType("*entity.Cart")).Return(nil)

	res, err := useCase.GetCart(context.Background(), mockUserID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, mockUserID, res.UserID)
	cartRepo.AssertExpectations(t)
}

func TestAddToCart_SuccessNewItem(t *testing.T) {
	_, cartRepo, productRepo, useCase := setupTest()

	req := &cartDTO.AddCartItemRequest{
		ProductID: mockProductID,
		Quantity:  2,
	}

	existingCart := &entity.Cart{ID: mockCartID, UserID: mockUserID}
	cartRepo.On("GetCartByUserID", mock.Anything, mockUserID).Return(existingCart, nil)
	
	productRepo.On("FindByID", mock.Anything, mockProductID).Return(&entity.Product{Stock: 10}, nil)
	cartRepo.On("GetCartItem", mock.Anything, mockCartID, mockProductID).Return(nil, gorm.ErrRecordNotFound)
	cartRepo.On("AddCartItem", mock.Anything, mock.MatchedBy(func(item *entity.CartItem) bool {
		return item.Quantity == 2 && item.ProductID == mockProductID
	})).Return(nil)

	err := useCase.AddToCart(context.Background(), mockUserID, req)

	assert.NoError(t, err)
	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestAddToCart_SuccessUpdateItem(t *testing.T) {
	_, cartRepo, productRepo, useCase := setupTest()

	req := &cartDTO.AddCartItemRequest{
		ProductID: mockProductID,
		Quantity:  2,
	}

	existingCart := &entity.Cart{ID: mockCartID, UserID: mockUserID}
	cartRepo.On("GetCartByUserID", mock.Anything, mockUserID).Return(existingCart, nil)
	
	productRepo.On("FindByID", mock.Anything, mockProductID).Return(&entity.Product{Stock: 10}, nil)
	
	existingItem := &entity.CartItem{ID: mockItemID, CartID: mockCartID, ProductID: mockProductID, Quantity: 3}
	cartRepo.On("GetCartItem", mock.Anything, mockCartID, mockProductID).Return(existingItem, nil)
	
	cartRepo.On("UpdateCartItemQuantity", mock.Anything, mockItemID, 5).Return(nil) // 3 + 2 = 5

	err := useCase.AddToCart(context.Background(), mockUserID, req)

	assert.NoError(t, err)
	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestAddToCart_InsufficientStock(t *testing.T) {
	_, cartRepo, productRepo, useCase := setupTest()

	req := &cartDTO.AddCartItemRequest{
		ProductID: mockProductID,
		Quantity:  12, // More than stock
	}

	existingCart := &entity.Cart{ID: mockCartID, UserID: mockUserID}
	cartRepo.On("GetCartByUserID", mock.Anything, mockUserID).Return(existingCart, nil)
	
	productRepo.On("FindByID", mock.Anything, mockProductID).Return(&entity.Product{Stock: 10}, nil)

	err := useCase.AddToCart(context.Background(), mockUserID, req)

	assert.Error(t, err)
	assert.Equal(t, "insufficient stock", err.Error())
	cartRepo.AssertNotCalled(t, "AddCartItem")
}

func TestUpdateCartItem_Success(t *testing.T) {
	_, cartRepo, productRepo, useCase := setupTest()

	req := &cartDTO.UpdateCartItemRequest{
		Quantity: 5,
	}

	existingCart := &entity.Cart{
		ID:     mockCartID,
		UserID: mockUserID,
		Items: []entity.CartItem{
			{ID: mockItemID, ProductID: mockProductID, Quantity: 2},
		},
	}
	
	cartRepo.On("GetCartByUserID", mock.Anything, mockUserID).Return(existingCart, nil)
	productRepo.On("FindByID", mock.Anything, mockProductID).Return(&entity.Product{Stock: 10}, nil)
	cartRepo.On("UpdateCartItemQuantity", mock.Anything, mockItemID, 5).Return(nil)

	err := useCase.UpdateCartItem(context.Background(), mockUserID, mockItemID, req)

	assert.NoError(t, err)
	cartRepo.AssertExpectations(t)
}

func TestUpdateCartItem_NotFoundInCart(t *testing.T) {
	_, cartRepo, _, useCase := setupTest()

	req := &cartDTO.UpdateCartItemRequest{Quantity: 5}

	existingCart := &entity.Cart{
		ID:     mockCartID,
		UserID: mockUserID,
		Items:  []entity.CartItem{}, // Empty cart
	}
	
	cartRepo.On("GetCartByUserID", mock.Anything, mockUserID).Return(existingCart, nil)

	err := useCase.UpdateCartItem(context.Background(), mockUserID, mockItemID, req)

	assert.Error(t, err)
	assert.Equal(t, "cart item not found in your cart", err.Error())
}

func TestBatchDeleteCartItems_Success(t *testing.T) {
	_, cartRepo, _, useCase := setupTest()

	req := &cartDTO.BatchDeleteCartItemsRequest{
		CartItemIDs: []uuid.UUID{mockItemID},
	}

	existingCart := &entity.Cart{ID: mockCartID, UserID: mockUserID}
	cartRepo.On("GetCartByUserID", mock.Anything, mockUserID).Return(existingCart, nil)
	cartRepo.On("DeleteCartItems", mock.Anything, mockCartID, req.CartItemIDs).Return(nil)

	err := useCase.BatchDeleteCartItems(context.Background(), mockUserID, req)

	assert.NoError(t, err)
	cartRepo.AssertExpectations(t)
}
