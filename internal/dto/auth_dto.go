package dto

type RegisterCustomerReq struct {
	Email     string `json:"email" binding:"required,email,max=100"`
	Password  string `json:"password" binding:"required,min=8,max=64"`
	FirstName string `json:"first_name" binding:"required,alpha,min=2,max=50"`
	LastName  string `json:"last_name" binding:"required,alpha,min=2,max=50"`
	Phone     string `json:"phone" binding:"omitempty,numeric,min=10,max=15"` 
	Address   string `json:"address" binding:"omitempty,max=255"`
}

type RegisterSellerReq struct {
	Email            string `json:"email" binding:"required,email,max=100"`
	Password         string `json:"password" binding:"required,min=8,max=64"`
	StoreName        string `json:"store_name" binding:"required,min=3,max=100"`
	StoreDescription string `json:"store_description" binding:"omitempty,max=500"`
	LogoUrl          string `json:"logo_url" binding:"omitempty,url,max=255"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}
