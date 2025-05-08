package directory

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hacKRD0/trikona_go/internal/directory-service/usecase"
	"github.com/hacKRD0/trikona_go/pkg/logger"
)

type MasterHandler struct {
	uc usecase.MasterUsecase
}

func NewMasterHandler(uc usecase.MasterUsecase) *MasterHandler {
	return &MasterHandler{uc: uc}
}

func (h *MasterHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/api/v1/directory/masters")
	{
		group.GET("/industries", h.GetAllIndustries)
		group.GET("/companies", h.GetAllCompanies)
		group.GET("/sectors", h.GetAllSectors)
		group.GET("/services", h.GetAllServices)
		group.GET("/states", h.GetAllStates)
		group.GET("/skills", h.GetAllSkills)
		group.GET("/countries", h.GetAllCountries)
		group.GET("/colleges", h.GetAllColleges)
	}
}

func (h *MasterHandler) GetAllIndustries(c *gin.Context) {
	industries, err := h.uc.GetAllIndustries()
	if err != nil {
		logger.Error("failed to get industries", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get industries"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": industries})
}

func (h *MasterHandler) GetAllCompanies(c *gin.Context) {
	companies, err := h.uc.GetAllCompanies()
	if err != nil {
		logger.Error("failed to get companies", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get companies"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": companies})
}

func (h *MasterHandler) GetAllSectors(c *gin.Context) {
	sectors, err := h.uc.GetAllSectors()
	if err != nil {
		logger.Error("failed to get sectors", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sectors"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sectors})
}

func (h *MasterHandler) GetAllServices(c *gin.Context) {
	services, err := h.uc.GetAllServices()
	if err != nil {
		logger.Error("failed to get services", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get services"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": services})
}

func (h *MasterHandler) GetAllStates(c *gin.Context) {
	states, err := h.uc.GetAllStates()
	if err != nil {
		logger.Error("failed to get states", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get states"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": states})
}

func (h *MasterHandler) GetAllSkills(c *gin.Context) {
	skills, err := h.uc.GetAllSkills()
	if err != nil {
		logger.Error("failed to get skills", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get skills"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": skills})
}

func (h *MasterHandler) GetAllCountries(c *gin.Context) {
	countries, err := h.uc.GetAllCountries()
	if err != nil {
		logger.Error("failed to get countries", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get countries"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": countries})
}

func (h *MasterHandler) GetAllColleges(c *gin.Context) {
	colleges, err := h.uc.GetAllColleges()
	if err != nil {
		logger.Error("failed to get colleges", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get colleges"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": colleges})
}
