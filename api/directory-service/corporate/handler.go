package directory

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/hacKRD0/trikona_go/internal/directory-service/usecase"
	"github.com/hacKRD0/trikona_go/pkg/utils"
)

type Handler struct {
	uc usecase.CorporateUsecase
}

func NewCorporateHandler(uc usecase.CorporateUsecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) GetCorporates(c *gin.Context) {
	var params domain.CorporateFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	params.Industries = utils.SplitStringArray(params.Industries)
	params.Sectors = utils.SplitStringArray(params.Sectors)
	params.Services = utils.SplitStringArray(params.Services)
	params.Country = utils.SplitStringArray(params.Country)
	params.States = utils.SplitStringArray(params.States)

	list, total, err := h.uc.FetchCorporates(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch corporates"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":       list,
		"page":       params.Page,
		"pageSize":   params.PageSize,
		"totalItems": total,
	})
}
