package handler

import (
	"net/http"

	"register-system-be/business"

	"github.com/gin-gonic/gin"
)

type AffiliationHandler struct {
	biz business.AffiliationBusiness
}

func NewAffiliationHandler(biz business.AffiliationBusiness) *AffiliationHandler {
	return &AffiliationHandler{biz: biz}
}

// @Summary      Get all affiliations
// @Description  Returns all valid affiliations for register form
// @Tags         Affiliation
// @Produce      json
// @Success      200 {array} model.AffiliationResponse
// @Failure      500 {object} map[string]string
// @Router       /affiliations [get]
func (h *AffiliationHandler) ListAll(c *gin.Context) {
	list, err := h.biz.ListAll()
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	Success(c, list)
}