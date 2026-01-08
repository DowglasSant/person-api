package contract

type RegisterDTO struct {
	Username string `json:"username" example:"john.doe" binding:"required,min=3,max=50"`
	Email    string `json:"email" example:"john.doe@company.com" binding:"required,email,max=100"`
	Password string `json:"password" example:"SecurePass123!" binding:"required,min=8,max=72"`
}
