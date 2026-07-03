package handlers

import (
	"net/http"

	"github.com/dat19/gin-ecommerce-api/internal/models"
	"github.com/dat19/gin-ecommerce-api/internal/service"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
	"github.com/dat19/gin-ecommerce-api/pkg/validator"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
	userService service.UserService
}

func NewAuthHandler(authSvc service.AuthService, userSvc service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authSvc,
		userService: userSvc,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Validate request
	if err := validator.ValidateRegisterRequest(req.Email, req.Username, req.Password, req.FirstName, req.LastName); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Normalize data
	req.Email = validator.NormalizeEmail(req.Email)
	req.Username = validator.NormalizeUsername(req.Username)

	resp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user with this email already exists" || err.Error() == "user with this username already exists" {
			status = http.StatusConflict
		}
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Validate request
	if err := validator.ValidateLoginRequest(req.Email, req.Password); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Normalize data
	req.Email = validator.NormalizeEmail(req.Email)

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		status := http.StatusUnauthorized
		if err.Error() == "account is disabled" {
			status = http.StatusForbidden
		}
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Logout successful", nil)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")

	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User retrieved successfully", user)
}
