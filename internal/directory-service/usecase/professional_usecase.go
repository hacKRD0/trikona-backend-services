package usecase

import (
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/hacKRD0/trikona_go/internal/directory-service/repository"
	"github.com/hacKRD0/trikona_go/pkg/pagination"
)

type ProfessionalUsecase interface {
    FetchProfessionals(params *domain.ProfessionalFilterParams) ([]domain.Professional, int64, error)
}

type professionalUsecase struct {
    repo repository.ProfessionalRepository
}

func NewProfessionalUsecase(r repository.ProfessionalRepository) ProfessionalUsecase {
    return &professionalUsecase{repo: r}
}

func (u *professionalUsecase) FetchProfessionals(params *domain.ProfessionalFilterParams) ([]domain.Professional, int64, error) {
    offset, limit := pagination.CalculateOffsetLimit(params.Page, params.PageSize)
    total, err := u.repo.Count(params)
    if err != nil {
        return nil, 0, err
    }
    list, err := u.repo.Find(params, offset, limit)
    return list, total, err
}
