// internal/directory-service/usecase/student_usecase.go
package usecase

import (
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/hacKRD0/trikona_go/internal/directory-service/repository"
	"github.com/hacKRD0/trikona_go/pkg/pagination"
)

type StudentUsecase interface {
    FetchStudents(params *domain.StudentFilterParams) ([]domain.Student, int64, error)
    GetStudentByID(id uint) (*domain.Student, error)
    CreateStudent(student *domain.Student) error
    UpdateStudent(student *domain.Student) error
    DeleteStudent(id uint) error
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

// GetStudentByID retrieves a student by ID
func (u *studentUsecase) GetStudentByID(id uint) (*domain.Student, error) {
    return u.repo.GetByID(id)
}

// CreateStudent adds a new student
func (u *studentUsecase) CreateStudent(student *domain.Student) error {
    return u.repo.Create(student)
}

// UpdateStudent modifies an existing student
func (u *studentUsecase) UpdateStudent(student *domain.Student) error {
    return u.repo.Update(student)
}

// DeleteStudent removes a student by ID
func (u *studentUsecase) DeleteStudent(id uint) error {
    return u.repo.Delete(id)
}
