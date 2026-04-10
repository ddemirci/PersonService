package repository

import (
	"context"
	"errors"
	"testing"

	"person-service/domain"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type mockDynamoDBClient struct {
	putItem func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	scan    func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

func (m *mockDynamoDBClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m.putItem(ctx, params, optFns...)
}

func (m *mockDynamoDBClient) Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return m.scan(ctx, params, optFns...)
}

func TestSave_Success(t *testing.T) {
	client := &mockDynamoDBClient{
		putItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	}
	repo := NewPersonRepository(client, "test-table")

	err := repo.Save(context.Background(), domain.Person{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSave_Error(t *testing.T) {
	client := &mockDynamoDBClient{
		putItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return nil, errors.New("dynamo error")
		},
	}
	repo := NewPersonRepository(client, "test-table")

	err := repo.Save(context.Background(), domain.Person{ID: "1"})

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestFindAll_Success(t *testing.T) {
	client := &mockDynamoDBClient{
		scan: func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
			return &dynamodb.ScanOutput{
				Items: []map[string]types.AttributeValue{
					{
						"id":          &types.AttributeValueMemberS{Value: "1"},
						"firstName":   &types.AttributeValueMemberS{Value: "John"},
						"lastName":    &types.AttributeValueMemberS{Value: "Doe"},
						"phoneNumber": &types.AttributeValueMemberS{Value: "123"},
						"address":     &types.AttributeValueMemberS{Value: "Street 1"},
					},
				},
			}, nil
		},
	}
	repo := NewPersonRepository(client, "test-table")

	persons, err := repo.FindAll(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(persons) != 1 {
		t.Fatalf("expected 1 person, got %d", len(persons))
	}
	if persons[0].FirstName != "John" {
		t.Errorf("expected John, got %s", persons[0].FirstName)
	}
}

func TestFindAll_Empty(t *testing.T) {
	client := &mockDynamoDBClient{
		scan: func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
			return &dynamodb.ScanOutput{Items: []map[string]types.AttributeValue{}}, nil
		},
	}
	repo := NewPersonRepository(client, "test-table")

	persons, err := repo.FindAll(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(persons) != 0 {
		t.Errorf("expected 0 persons, got %d", len(persons))
	}
}

func TestFindAll_Error(t *testing.T) {
	client := &mockDynamoDBClient{
		scan: func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
			return nil, errors.New("dynamo error")
		},
	}
	repo := NewPersonRepository(client, "test-table")

	_, err := repo.FindAll(context.Background())

	if err == nil {
		t.Error("expected error, got nil")
	}
}
