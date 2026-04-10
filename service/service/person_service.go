package service

import (
	"context"

	"person-service/domain"
	"person-service/ports"

	"github.com/google/uuid"
)

type PersonService struct {
	repository ports.PersonRepository
}

func NewPersonService(repository ports.PersonRepository) *PersonService {
	return &PersonService{repository: repository}
}

func (s *PersonService) CreatePerson(ctx context.Context, person domain.Person) (domain.Person, error) {
	person.ID = uuid.New().String()

	if err := s.repository.Save(ctx, person); err != nil {
		return domain.Person{}, err
	}

	return person, nil
}

func (s *PersonService) ListPersons(ctx context.Context) ([]domain.Person, error) {
	return s.repository.FindAll(ctx)
}
