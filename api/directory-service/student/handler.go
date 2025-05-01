// api/directory-service/student/handler.go
package student

import (
	"net/http"
	"strconv"
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

// GetStudent retrieves a student by ID
func (h *StudentHandler) GetStudent(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student ID"})
        return
    }

    student, err := h.uc.GetStudentByID(uint(id))
    if err != nil {
        logger.Error("failed to get student", err)
        c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": student})
}

// CreateStudent adds a new student
func (h *StudentHandler) CreateStudent(c *gin.Context) {
    var student domain.Student
    if err := c.ShouldBindJSON(&student); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.uc.CreateStudent(&student); err != nil {
        logger.Error("failed to create student", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create student"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"data": student})
}

// UpdateStudent modifies an existing student
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student ID"})
        return
    }

    var student domain.Student
    if err := c.ShouldBindJSON(&student); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Ensure the ID in the URL matches the student object
    student.ID = uint(id)

    if err := h.uc.UpdateStudent(&student); err != nil {
        logger.Error("failed to update student", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update student"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": student})
}

// DeleteStudent removes a student
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student ID"})
        return
    }

    if err := h.uc.DeleteStudent(uint(id)); err != nil {
        logger.Error("failed to delete student", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete student"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "student deleted successfully"})
}
