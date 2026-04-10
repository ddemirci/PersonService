package handler

import (
	"context"
	"errors"
	"testing"

	"person-service/domain"
)

type mockPersonService struct {
	createPerson func(ctx context.Context, person domain.Person) (domain.Person, error)
	listPersons  func(ctx context.Context) ([]domain.Person, error)
}

func (m *mockPersonService) CreatePerson(ctx context.Context, person domain.Person) (domain.Person, error) {
	return m.createPerson(ctx, person)
}

func (m *mockPersonService) ListPersons(ctx context.Context) ([]domain.Person, error) {
	return m.listPersons(ctx)
}

func TestHandle_UnsupportedMethod(t *testing.T) {
	h := NewPersonHandler(&mockPersonService{})

	result, err := h.Handle(context.Background(), map[string]interface{}{
		"httpMethod": "DELETE",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["statusCode"] != 400 {
		t.Errorf("expected 400, got %v", result["statusCode"])
	}
}

func TestHandle_Post_Success(t *testing.T) {
	svc := &mockPersonService{
		createPerson: func(ctx context.Context, person domain.Person) (domain.Person, error) {
			person.ID = "123"
			return person, nil
		},
	}
	h := NewPersonHandler(svc)

	result, err := h.Handle(context.Background(), map[string]interface{}{
		"httpMethod": "POST",
		"body":       `{"firstName":"John","lastName":"Doe"}`,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["statusCode"] != 200 {
		t.Errorf("expected 200, got %v", result["statusCode"])
	}
}

func TestHandle_Post_InvalidBody(t *testing.T) {
	h := NewPersonHandler(&mockPersonService{})

	result, err := h.Handle(context.Background(), map[string]interface{}{
		"httpMethod": "POST",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["statusCode"] != 400 {
		t.Errorf("expected 400, got %v", result["statusCode"])
	}
}

func TestHandle_Post_InvalidJSON(t *testing.T) {
	h := NewPersonHandler(&mockPersonService{})

	result, err := h.Handle(context.Background(), map[string]interface{}{
		"httpMethod": "POST",
		"body":       `not-json`,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["statusCode"] != 400 {
		t.Errorf("expected 400, got %v", result["statusCode"])
	}
}

func TestHandle_Post_ServiceError(t *testing.T) {
	svc := &mockPersonService{
		createPerson: func(ctx context.Context, person domain.Person) (domain.Person, error) {
			return domain.Person{}, errors.New("service error")
		},
	}
	h := NewPersonHandler(svc)

	result, err := h.Handle(context.Background(), map[string]interface{}{
		"httpMethod": "POST",
		"body":       `{"firstName":"John"}`,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["statusCode"] != 500 {
		t.Errorf("expected 500, got %v", result["statusCode"])
	}
}

func TestHandle_Get_Success(t *testing.T) {
	svc := &mockPersonService{
		listPersons: func(ctx context.Context) ([]domain.Person, error) {
			return []domain.Person{
				{ID: "1", FirstName: "John", LastName: "Doe"},
			}, nil
		},
	}
	h := NewPersonHandler(svc)

	result, err := h.Handle(context.Background(), map[string]interface{}{
		"httpMethod": "GET",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["statusCode"] != 200 {
		t.Errorf("expected 200, got %v", result["statusCode"])
	}
}

func TestHandle_Post_HTTPv2Format(t *testing.T) {
	svc := &mockPersonService{
		createPerson: func(ctx context.Context, person domain.Person) (domain.Person, error) {
			person.ID = "123"
			return person, nil
		},
	}
	h := NewPersonHandler(svc)

	result, err := h.Handle(context.Background(), map[string]interface{}{
		"requestContext": map[string]interface{}{
			"http": map[string]interface{}{
				"method": "POST",
			},
		},
		"body": `{"firstName":"John","lastName":"Doe"}`,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["statusCode"] != 200 {
		t.Errorf("expected 200, got %v", result["statusCode"])
	}
}

func TestHandle_Get_ServiceError(t *testing.T) {
	svc := &mockPersonService{
		listPersons: func(ctx context.Context) ([]domain.Person, error) {
			return nil, errors.New("service error")
		},
	}
	h := NewPersonHandler(svc)

	result, err := h.Handle(context.Background(), map[string]interface{}{
		"httpMethod": "GET",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["statusCode"] != 500 {
		t.Errorf("expected 500, got %v", result["statusCode"])
	}
}
