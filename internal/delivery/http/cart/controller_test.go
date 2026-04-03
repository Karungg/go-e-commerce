package cart_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cartCtrl "go-e-commerce/internal/delivery/http/cart"
	cartDTO "go-e-commerce/internal/dto/cart"
	cartMock "go-e-commerce/internal/mocks/cart"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	mockUserID    = uuid.New()
	mockCartID    = uuid.New()
	mockProductID = uuid.New()
	mockItemID    = uuid.New()
)

func setupControllerTest() (*cartMock.CartUseCaseMock, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	useCase := &cartMock.CartUseCaseMock{}
	controller := cartCtrl.NewCartController(useCase)

	router := gin.Default()
	
	authMiddleware := func(c *gin.Context) {
		c.Set("userID", mockUserID.String())
		c.Next()
	}

	cartGrp := router.Group("/api/cart")
	cartGrp.Use(authMiddleware)
	{
		cartGrp.GET("", controller.GetCart)
		cartGrp.POST("", controller.AddToCart)
		cartGrp.PUT("/:id", controller.UpdateCartItem)
		cartGrp.DELETE("/batch", controller.BatchDeleteCartItems)
	}

	return useCase, router
}

func TestGetCartController_Success(t *testing.T) {
	useCase, router := setupControllerTest()

	expectedResponse := &cartDTO.CartResponse{
		ID: mockCartID,
		UserID: mockUserID,
	}

	useCase.On("GetCart", mock.Anything, mockUserID).Return(expectedResponse, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/cart", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	useCase.AssertExpectations(t)
}

func TestAddToCartController_Success(t *testing.T) {
	useCase, router := setupControllerTest()

	payload := cartDTO.AddCartItemRequest{
		ProductID: mockProductID,
		Quantity:  2,
	}
	body, _ := json.Marshal(payload)

	useCase.On("AddToCart", mock.Anything, mockUserID, mock.AnythingOfType("*cart.AddCartItemRequest")).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/api/cart", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	useCase.AssertExpectations(t)
}

func TestAddToCartController_InvalidPayload(t *testing.T) {
	useCase, router := setupControllerTest()

	// Missing required fields
	payload := map[string]interface{}{
		"missing_fields": true,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/cart", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	useCase.AssertNotCalled(t, "AddToCart")
}

func TestUpdateCartItemController_Success(t *testing.T) {
	useCase, router := setupControllerTest()

	payload := cartDTO.UpdateCartItemRequest{
		Quantity: 5,
	}
	body, _ := json.Marshal(payload)

	useCase.On("UpdateCartItem", mock.Anything, mockUserID, mockItemID, mock.AnythingOfType("*cart.UpdateCartItemRequest")).Return(nil)

	req, _ := http.NewRequest(http.MethodPut, "/api/cart/"+mockItemID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	useCase.AssertExpectations(t)
}

func TestBatchDeleteCartItemsController_Success(t *testing.T) {
	useCase, router := setupControllerTest()

	payload := cartDTO.BatchDeleteCartItemsRequest{
		CartItemIDs: []uuid.UUID{mockItemID},
	}
	body, _ := json.Marshal(payload)

	useCase.On("BatchDeleteCartItems", mock.Anything, mockUserID, mock.AnythingOfType("*cart.BatchDeleteCartItemsRequest")).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/api/cart/batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	useCase.AssertExpectations(t)
}
