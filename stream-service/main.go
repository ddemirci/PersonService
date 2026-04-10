package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

type DynamoDBEvent struct {
	Records []struct {
		EventName string `json:"eventName"`
		Dynamodb  struct {
			NewImage map[string]struct {
				S string `json:"S"`
			} `json:"NewImage"`
		} `json:"dynamodb"`
	} `json:"Records"`
}

func handler(ctx context.Context, event DynamoDBEvent) error {

	cfg, _ := config.LoadDefaultConfig(ctx)
	client := eventbridge.NewFromConfig(cfg)

	for _, record := range event.Records {

		if record.EventName != "INSERT" {
			continue
		}

		person := map[string]string{
			"id":          record.Dynamodb.NewImage["id"].S,
			"firstName":   record.Dynamodb.NewImage["firstName"].S,
			"lastName":    record.Dynamodb.NewImage["lastName"].S,
			"phoneNumber": record.Dynamodb.NewImage["phoneNumber"].S,
			"address":     record.Dynamodb.NewImage["address"].S,
		}

		detail, _ := json.Marshal(person)
		busName := os.Getenv("EVENT_BUS_NAME")

		_, err := client.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					Source:       awsString("person.service"),
					DetailType:   awsString("PersonCreated"),
					Detail:       awsString(string(detail)),
					EventBusName: awsString(busName),
				},
			},
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func awsString(s string) *string {
	return &s
}

func main() {
	lambda.Start(handler)
}
