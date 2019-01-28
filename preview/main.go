package main

import (
	"fmt"
	"os"

	"github.com/badoux/goscraper"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

type HandlerConfig struct {
	c     dynamodbiface.DynamoDBAPI
	table string
}

func (hc *HandlerConfig) Handler(event events.DynamoDBEvent) (Response, error) {

	for _, r := range event.Records {
		url, ok := r.Change.NewImage["url"]
		if !ok {
			return Response{StatusCode: 501}, fmt.Errorf("cant handle event: %v", event)
		}

		s, err := goscraper.Scrape(url.String(), 5)
		if err != nil {
			logrus.WithField("error", err).Errorf("failed to scrape '%s'", url)
		}

		_, err = hc.c.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(hc.table),
			Item: map[string]*dynamodb.AttributeValue{
				"url":   &dynamodb.AttributeValue{S: aws.String(url.String())},
				"image": &dynamodb.AttributeValue{S: aws.String(s.Preview.Images[0])},
				"name":  &dynamodb.AttributeValue{S: aws.String(s.Preview.Name)},
				"title": &dynamodb.AttributeValue{S: aws.String(s.Preview.Title)},
			}})

		if err != nil {
			logrus.WithField("error", err).Error("Couldn't save Preview")
		}
	}

	resp := Response{
		StatusCode: 201,
		Body:       "s",
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}

	return resp, nil
}

func main() {
	// Create a new AWS session and fail immediately on error
	sess := session.Must(session.NewSession())
	// Create the DynamoDB client
	dynamodbclient := dynamodb.New(sess)

	hc := HandlerConfig{c: dynamodbclient, table: os.Getenv("PREVIEW_TABLE")}
	lambda.Start(hc.Handler)
}
