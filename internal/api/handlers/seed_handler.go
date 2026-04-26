package handlers

import (
	"net/http"

	"github.com/dat19/gin-ecommerce-api/internal/service"
	"github.com/dat19/gin-ecommerce-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type SeedHandler struct {
	seedService service.SeedService
}

func NewSeedHandler(seedSvc service.SeedService) *SeedHandler {
	return &SeedHandler{seedService: seedSvc}
}

func (h *SeedHandler) SeedData(c *gin.Context) {
	err := h.seedService.SeedAll(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed data: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Database seeded successfully", nil)
}
