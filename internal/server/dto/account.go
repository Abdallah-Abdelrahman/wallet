package dto

type CreateAccountRequest struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type CreateAccountResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Balance   float64 `json:"balance"`
}

type TopUpRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type TopUpResponse struct {
	TransactionID string  `json:"transaction_id"`
	AccountID     string  `json:"account_id"`
	Amount        float64 `json:"amount"`
	NewBalance    float64 `json:"new_balance"`
}

type ChargeRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type ChargeResponse struct {
	TransactionID string  `json:"transaction_id"`
	AccountID     string  `json:"account_id"`
	Amount        float64 `json:"amount"`
	NewBalance    float64 `json:"new_balance"`
}

type AccountBadRequestResponse struct {
	Error string `json:"error" default:"Invalid request"`
}
