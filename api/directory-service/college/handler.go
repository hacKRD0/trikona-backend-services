package college

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/hacKRD0/trikona_go/internal/directory-service/usecase"
)

type Handler struct {
    uc usecase.CollegeUsecase
}

func NewHandler(uc usecase.CollegeUsecase) *Handler {
    return &Handler{uc: uc}
}

func (h *Handler) GetColleges(c *gin.Context) {
    var params domain.CollegeFilterParams
    if err := c.ShouldBindQuery(&params); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
        return
    }
    list, total, err := h.uc.FetchColleges(&params)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch colleges"})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "data":       list,
        "page":       params.Page,
        "pageSize":   params.PageSize,
        "totalItems": total,
    })
}
