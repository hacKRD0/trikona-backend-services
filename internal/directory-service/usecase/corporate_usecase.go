package usecase

import (
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/hacKRD0/trikona_go/internal/directory-service/repository"
	"github.com/hacKRD0/trikona_go/pkg/pagination"
)

type CorporateUsecase interface {
    FetchCorporates(params *domain.CorporateFilterParams) ([]domain.Corporate, int64, error)
}

type corporateUsecase struct {
    repo repository.CorporateRepository
}

func NewCorporateUsecase(r repository.CorporateRepository) CorporateUsecase {
    return &corporateUsecase{repo: r}
}

func (u *corporateUsecase) FetchCorporates(params *domain.CorporateFilterParams) ([]domain.Corporate, int64, error) {
    offset, limit := pagination.CalculateOffsetLimit(params.Page, params.PageSize)
    total, err := u.repo.Count(params)
    if err != nil {
        return nil, 0, err
    }
    list, err := u.repo.Find(params, offset, limit)
    return list, total, err
}
