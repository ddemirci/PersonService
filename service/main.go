package main

import (
	"context"
	"os"

	"person-service/handler"
	"person-service/repository"
	"person-service/service"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var hd *handler.PersonHandler

func main() {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	client := dynamodb.NewFromConfig(cfg)
	repo := repository.NewPersonRepository(client, os.Getenv("TABLE_NAME"))
	svc := service.NewPersonService(repo)
	hd = handler.NewPersonHandler(svc)
	lambda.Start(hd.Handle)
}
