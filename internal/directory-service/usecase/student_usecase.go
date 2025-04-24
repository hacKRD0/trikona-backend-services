// internal/directory-service/usecase/student_usecase.go
package usecase

import (
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/hacKRD0/trikona_go/internal/directory-service/repository"
	"github.com/hacKRD0/trikona_go/pkg/pagination"
)

type StudentUsecase interface {
    FetchStudents(params *domain.StudentFilterParams) ([]domain.Student, int64, error)
}

type studentUsecase struct {
    repo repository.StudentRepository
}

func NewStudentUsecase(repo repository.StudentRepository) StudentUsecase {
    return &studentUsecase{repo: repo}
}

func (u *studentUsecase) FetchStudents(params *domain.StudentFilterParams) ([]domain.Student, int64, error) {
    offset, limit := pagination.CalculateOffsetLimit(params.Page, params.PageSize)

    total, err := u.repo.Count(params)
    if err != nil {
        return nil, 0, err
    }

    students, err := u.repo.Find(params, offset, limit)
    if err != nil {
        return nil, 0, err
    }

    return students, total, nil
}
