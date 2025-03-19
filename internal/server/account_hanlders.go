package server

import (
	"net/http"

	"wallet/internal/server/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateAccountHandler creates a new account with the given user details
// @Summary Create a new account
// @Description Create a new account with the given user details
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body dto.CreateAccountRequest true "Account details"
// @Success 201 {object} dto.CreateAccountResponse "Account created successfully"
// @Failure 400 {object} dto.AccountBadRequestResponse  "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /accounts [post]
func (s *Server) CreateAccountHandler(c *gin.Context) {
	var request dto.CreateAccountRequest

	// Check if the request body is empty
	if c.Request.Body == nil || c.Request.ContentLength == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body cannot be empty"})
		return
	}
	// Bind the request body to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the user and account with 0 balance
	account, err := s.AccountService.CreateAccountWithUser(request.Email, request.FirstName, request.LastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// TopUpHandler tops up the account with the given amount
// @Summary Top up an account
// @Description Top up an account with the given amount
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param request body dto.TopUpRequest true "Top up details"
// @Success 200 {object} dto.TopUpResponse "Top up successful"
// @Failure 400 {object} dto.AccountBadRequestResponse  "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /accounts/{id}/top-up [post]
func (s *Server) TopUpHandler(c *gin.Context) {
	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	var request struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the account service to top up the account
	transaction, err := s.AccountService.TopUp(accountID, request.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// ChargeHandler charges the account with the given amount
// @Summary Charge an account
// @Description Charge an account with the given amount
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param request body dto.ChargeRequest true "Charge details"
// @Success 200 {object} dto.ChargeResponse "Charge successful"
// @Failure 400 {object} dto.AccountBadRequestResponse  "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /accounts/{id}/charge [post]
func (s *Server) ChargeHandler(c *gin.Context) {
	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	var request struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the account service to charge the account
	transaction, err := s.AccountService.Charge(accountID, request.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
