package controller

import (
	"errors"
	"net/http"

	"zeum-license-api/internal/apperror"
	"zeum-license-api/internal/service"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	service *service.ProfileService
}

func NewProfileController(service *service.ProfileService) *ProfileController {
	return &ProfileController{service: service}
}

func (c *ProfileController) FindAction(ctx *gin.Context) {

	email := ctx.Param("email")

	data, err := c.service.FindByEmail(email)

	if err != nil {

		if errors.Is(err, apperror.ErrProfileNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Perfil do usuário não encontrado."})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}
