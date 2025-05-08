package repository

import (
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"gorm.io/gorm"
)

type MasterRepository interface {
	GetAllIndustries() ([]domain.IndustryMaster, error)
	GetAllCompanies() ([]domain.CompanyMaster, error)
	GetAllSectors() ([]domain.SectorMaster, error)
	GetAllServices() ([]domain.ServiceMaster, error)
	GetAllStates() ([]domain.StateMaster, error)
	GetAllSkills() ([]domain.SkillMaster, error)
	GetAllCountries() ([]domain.CountryMaster, error)
	GetAllColleges() ([]domain.CollegeMaster, error)
}

type masterRepository struct {
	db *gorm.DB
}

func NewMasterRepository(db *gorm.DB) MasterRepository {
	return &masterRepository{db: db}
}

func (r *masterRepository) GetAllIndustries() ([]domain.IndustryMaster, error) {
	var industries []domain.IndustryMaster
	if err := r.db.Find(&industries).Error; err != nil {
		return nil, err
	}
	return industries, nil
}

func (r *masterRepository) GetAllCompanies() ([]domain.CompanyMaster, error) {
	var companies []domain.CompanyMaster
	if err := r.db.Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

func (r *masterRepository) GetAllSectors() ([]domain.SectorMaster, error) {
	var sectors []domain.SectorMaster
	if err := r.db.Find(&sectors).Error; err != nil {
		return nil, err
	}
	return sectors, nil
}

func (r *masterRepository) GetAllServices() ([]domain.ServiceMaster, error) {
	var services []domain.ServiceMaster
	if err := r.db.Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func (r *masterRepository) GetAllStates() ([]domain.StateMaster, error) {
	var states []domain.StateMaster
	if err := r.db.Preload("Country").Find(&states).Error; err != nil {
		return nil, err
	}
	return states, nil
}

func (r *masterRepository) GetAllSkills() ([]domain.SkillMaster, error) {
	var skills []domain.SkillMaster
	if err := r.db.Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}

func (r *masterRepository) GetAllCountries() ([]domain.CountryMaster, error) {
	var countries []domain.CountryMaster
	if err := r.db.Find(&countries).Error; err != nil {
		return nil, err
	}
	return countries, nil
}

func (r *masterRepository) GetAllColleges() ([]domain.CollegeMaster, error) {
	var colleges []domain.CollegeMaster
	if err := r.db.Find(&colleges).Error; err != nil {
		return nil, err
	}
	return colleges, nil
}
