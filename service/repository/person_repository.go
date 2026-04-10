package repository

import (
	"context"

	"person-service/domain"
	"person-service/ports"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type PersonRepository struct {
	client    ports.DynamoDBClient
	tableName string
}

func NewPersonRepository(client ports.DynamoDBClient, tableName string) *PersonRepository {
	return &PersonRepository{client: client, tableName: tableName}
}

func (r *PersonRepository) Save(ctx context.Context, person domain.Person) error {
	_, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item: map[string]types.AttributeValue{
			"id":          &types.AttributeValueMemberS{Value: person.ID},
			"firstName":   &types.AttributeValueMemberS{Value: person.FirstName},
			"lastName":    &types.AttributeValueMemberS{Value: person.LastName},
			"phoneNumber": &types.AttributeValueMemberS{Value: person.PhoneNumber},
			"address":     &types.AttributeValueMemberS{Value: person.Address},
		},
	})
	return err
}

func (r *PersonRepository) FindAll(ctx context.Context) ([]domain.Person, error) {
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: &r.tableName,
	})
	if err != nil {
		return nil, err
	}

	var persons []domain.Person
	for _, item := range result.Items {
		p := domain.Person{}
		if v, ok := item["id"].(*types.AttributeValueMemberS); ok {
			p.ID = v.Value
		}
		if v, ok := item["firstName"].(*types.AttributeValueMemberS); ok {
			p.FirstName = v.Value
		}
		if v, ok := item["lastName"].(*types.AttributeValueMemberS); ok {
			p.LastName = v.Value
		}
		if v, ok := item["phoneNumber"].(*types.AttributeValueMemberS); ok {
			p.PhoneNumber = v.Value
		}
		if v, ok := item["address"].(*types.AttributeValueMemberS); ok {
			p.Address = v.Value
		}
		persons = append(persons, p)
	}

	return persons, nil
}
