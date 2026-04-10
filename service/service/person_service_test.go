package service

import (
	"context"
	"errors"
	"testing"

	"person-service/domain"
)

type mockPersonRepository struct {
	save    func(ctx context.Context, person domain.Person) error
	findAll func(ctx context.Context) ([]domain.Person, error)
}

func (m *mockPersonRepository) Save(ctx context.Context, person domain.Person) error {
	return m.save(ctx, person)
}

func (m *mockPersonRepository) FindAll(ctx context.Context) ([]domain.Person, error) {
	return m.findAll(ctx)
}

func TestCreatePerson_Success(t *testing.T) {
	rr := &mockPersonRepository{
		save: func(ctx context.Context, person domain.Person) error {
			return nil
		},
	}

	svc := NewPersonService(rr)

	result, err := svc.CreatePerson(context.Background(), domain.Person{
		FirstName: "John",
		LastName:  "Doe",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID == "" {
		t.Error("expected ID to be set by service, got empty string")
	}
	if result.FirstName != "John" {
		t.Errorf("expected FirstName John, got %s", result.FirstName)
	}
}

func TestCreatePerson_Error(t *testing.T) {
	rr := &mockPersonRepository{
		save: func(ctx context.Context, person domain.Person) error {
			return errors.New("Repository error")
		},
	}

	svc := NewPersonService(rr)

	_, err := svc.CreatePerson(context.Background(), domain.Person{
		FirstName: "John",
		LastName:  "Doe",
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestListPersons_Success(t *testing.T) {

	persons := []domain.Person{
		{
			ID:          "1",
			FirstName:   "John",
			LastName:    "Doe",
			PhoneNumber: "1234567",
			Address:     "Amsterdam",
		},
		{
			ID:          "2",
			FirstName:   "Mary",
			LastName:    "Doe",
			PhoneNumber: "1234568",
			Address:     "Amsterdam",
		},
	}

	rr := &mockPersonRepository{
		findAll: func(ctx context.Context) ([]domain.Person, error) {
			return persons, nil
		},
	}

	svc := NewPersonService(rr)

	result, err := svc.ListPersons(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 persons, got %d", len(result))

	}
}

func TestListPersons_Error(t *testing.T) {
	rr := &mockPersonRepository{
		findAll: func(ctx context.Context) ([]domain.Person, error) {
			return nil, errors.New("Repository error")
		},
	}

	svc := NewPersonService(rr)

	_, err := svc.ListPersons(context.Background())

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
