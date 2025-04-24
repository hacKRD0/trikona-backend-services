// api/directory-service/student/handler.go
package student

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/hacKRD0/trikona_go/internal/directory-service/usecase"
	"github.com/hacKRD0/trikona_go/pkg/logger"
)

type StudentHandler struct {
    uc usecase.StudentUsecase
}

func NewStudentHandler(uc usecase.StudentUsecase) *StudentHandler {
    return &StudentHandler{uc: uc}
}

func (h *StudentHandler) GetStudents(c *gin.Context) {
    var params domain.StudentFilterParams
    // bind all the rest of the query params
    if err := c.ShouldBindQuery(&params); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

   // manually split comma‐separated skills into a proper slice
   if raw := c.Query("skills"); raw != "" {
       parts := strings.Split(raw, ",")
       for i := range parts {
           parts[i] = strings.TrimSpace(parts[i])
       }
       params.Skills = parts
   }

    students, total, err := h.uc.FetchStudents(&params)
    if err != nil {
        logger.Error("failed to fetch students", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch students"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data":       students,
        "page":       params.Page,
        "pageSize":   params.PageSize,
        "totalItems": total,
    })
}
