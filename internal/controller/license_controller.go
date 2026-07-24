package controller

import (
	"errors"
	"net/http"
	"strconv"

	"zeum-license-api/internal/apperror"
	"zeum-license-api/internal/dto"
	"zeum-license-api/internal/service"

	"github.com/gin-gonic/gin"
)

type LicenseController struct {
	service *service.LicenseService
}

func NewLicenseController(service *service.LicenseService) *LicenseController {
	return &LicenseController{service: service}
}

func (c *LicenseController) FindAction(ctx *gin.Context) {

	var request dto.FindLicenseRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"license_number": "This value should not be blank."})
		return
	}

	data, err := c.service.FindByLicenseNumber(request.LicenseNumber)

	if err != nil {

		if errors.Is(err, apperror.ErrLicenseNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Licença não encontrada."})
			return
		}

		if errors.Is(err, apperror.ErrLicenseExpired) {
			ctx.JSON(http.StatusPaymentRequired, gin.H{"message": "Licença expirada."})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (c *LicenseController) FindLicenseNumberByTenantAction(ctx *gin.Context) {

	tenantID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"id": "This value should be a valid tenant ID."})
		return
	}

	data, err := c.service.FindLicenseNumberByTenantID(uint(tenantID))

	if err != nil {

		if errors.Is(err, apperror.ErrTenantNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Tenant não encontrado."})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}
