package ports

import (
	"context"

	"person-service/domain"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

type PersonRepository interface {
	Save(ctx context.Context, person domain.Person) error
	FindAll(ctx context.Context) ([]domain.Person, error)
}

type PersonService interface {
	CreatePerson(ctx context.Context, person domain.Person) (domain.Person, error)
	ListPersons(ctx context.Context) ([]domain.Person, error)
}
