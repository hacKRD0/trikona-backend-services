package usecase

import (
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/hacKRD0/trikona_go/internal/directory-service/repository"
)

type MasterUsecase interface {
	GetAllIndustries() ([]domain.IndustryMaster, error)
	GetAllCompanies() ([]domain.CompanyMaster, error)
	GetAllSectors() ([]domain.SectorMaster, error)
	GetAllServices() ([]domain.ServiceMaster, error)
	GetAllStates() ([]domain.StateMaster, error)
	GetAllSkills() ([]domain.SkillMaster, error)
	GetAllCountries() ([]domain.CountryMaster, error)
	GetAllColleges() ([]domain.CollegeMaster, error)
}

type masterUsecase struct {
	repo repository.MasterRepository
}

func NewMasterUsecase(repo repository.MasterRepository) MasterUsecase {
	return &masterUsecase{repo: repo}
}

func (u *masterUsecase) GetAllIndustries() ([]domain.IndustryMaster, error) {
	return u.repo.GetAllIndustries()
}

func (u *masterUsecase) GetAllCompanies() ([]domain.CompanyMaster, error) {
	return u.repo.GetAllCompanies()
}

func (u *masterUsecase) GetAllSectors() ([]domain.SectorMaster, error) {
	return u.repo.GetAllSectors()
}

func (u *masterUsecase) GetAllServices() ([]domain.ServiceMaster, error) {
	return u.repo.GetAllServices()
}

func (u *masterUsecase) GetAllStates() ([]domain.StateMaster, error) {
	return u.repo.GetAllStates()
}

func (u *masterUsecase) GetAllSkills() ([]domain.SkillMaster, error) {
	return u.repo.GetAllSkills()
}

func (u *masterUsecase) GetAllCountries() ([]domain.CountryMaster, error) {
	return u.repo.GetAllCountries()
}

func (u *masterUsecase) GetAllColleges() ([]domain.CollegeMaster, error) {
	return u.repo.GetAllColleges()
}
