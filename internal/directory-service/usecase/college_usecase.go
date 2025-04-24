package usecase

import (
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/hacKRD0/trikona_go/internal/directory-service/repository"
	"github.com/hacKRD0/trikona_go/pkg/pagination"
)

type CollegeUsecase interface {
    FetchColleges(params *domain.CollegeFilterParams) ([]domain.College, int64, error)
}

type collegeUsecase struct {
    repo repository.CollegeRepository
}

func NewCollegeUsecase(r repository.CollegeRepository) CollegeUsecase {
    return &collegeUsecase{repo: r}
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
