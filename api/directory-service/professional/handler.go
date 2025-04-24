package professional

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/hacKRD0/trikona_go/internal/directory-service/usecase"
)

type Handler struct {
    uc usecase.ProfessionalUsecase
}

func NewHandler(uc usecase.ProfessionalUsecase) *Handler {
    return &Handler{uc: uc}
}

func (h *Handler) GetProfessionals(c *gin.Context) {
    var params domain.ProfessionalFilterParams
    if err := c.ShouldBindQuery(&params); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
        return
    }
    list, total, err := h.uc.FetchProfessionals(&params)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch professionals"})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "data":       list,
        "page":       params.Page,
        "pageSize":   params.PageSize,
        "totalItems": total,
    })
}
