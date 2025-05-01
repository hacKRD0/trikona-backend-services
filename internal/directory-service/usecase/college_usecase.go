package usecase

import (
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/hacKRD0/trikona_go/internal/directory-service/repository"
	"github.com/hacKRD0/trikona_go/pkg/pagination"
)

type CollegeUsecase interface {
	FetchColleges(params *domain.CollegeFilterParams) ([]domain.College, int64, error)
	GetCollegeByID(id uint) (*domain.College, error)
	CreateCollege(college *domain.College) error
	UpdateCollege(college *domain.College) error
	DeleteCollege(id uint) error
}

type collegeUsecase struct {
	repo repository.CollegeRepository
}

func NewCollegeUsecase(r repository.CollegeRepository) CollegeUsecase {
	return &collegeUsecase{repo: r}
}

func (u *collegeUsecase) GetCollegeByID(id uint) (*domain.College, error) {
	return u.repo.GetByID(id)
}

func (u *collegeUsecase) CreateCollege(college *domain.College) error {
	return u.repo.Create(college)
}

func (u *collegeUsecase) UpdateCollege(college *domain.College) error {
	return u.repo.Update(college)
}

func (u *collegeUsecase) DeleteCollege(id uint) error {
	return u.repo.Delete(id)
}

func (u *collegeUsecase) FetchColleges(params *domain.CollegeFilterParams) ([]domain.College, int64, error) {
	offset, limit := pagination.CalculateOffsetLimit(params.Page, params.PageSize)
	total, err := u.repo.Count(params)
	if err != nil {
		return nil, 0, err
	}
	list, err := u.repo.Find(params, offset, limit)
	return list, total, err
}
