package college

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/hacKRD0/trikona_go/internal/directory-service/usecase"
	"github.com/hacKRD0/trikona_go/pkg/logger"
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

func (h *Handler) GetCollege(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	college, err := h.uc.GetCollegeByID(uint(id))
	if err != nil {
		logger.Error("failed to get college", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch college"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": college})
}

func (h *Handler) CreateCollege(c *gin.Context) {
	var college domain.College
	if err := c.ShouldBindJSON(&college); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.CreateCollege(&college); err != nil {
		logger.Error("failed to create college", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create college"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": college})
}

func (h *Handler) UpdateCollege(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid college ID"})
		return
	}

	var college domain.College
	if err := c.ShouldBindJSON(&college); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the ID in the URL matches the college object
	college.ID = uint(id)

	if err := h.uc.UpdateCollege(&college); err != nil {
		logger.Error("failed to update college", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update college"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": college})
}

func (h *Handler) DeleteCollege(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid college ID"})
		return
	}

	if err := h.uc.DeleteCollege(uint(id)); err != nil {
		logger.Error("failed to delete college", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete college"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "College deleted successfully"})
}
