package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

type Response events.APIGatewayProxyResponse

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	u, ok := request.QueryStringParameters["url"]
	if !ok {
		return Response{StatusCode: 501}, fmt.Errorf("missing required name parameter")
	}

	s := "short"

	// Create a new AWS session and fail immediately on error
	sess := session.Must(session.NewSession())
	// Create the DynamoDB client
	dynamodbclient := dynamodb.New(sess)
	_, err := dynamodbclient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DYNAMO_DB_TABLE")),
		Item: map[string]*dynamodb.AttributeValue{
			"short_url": &dynamodb.AttributeValue{S: aws.String(s)},
			"url":       &dynamodb.AttributeValue{S: aws.String(u)},
		}})
	if err != nil {
		log.WithField("error", err).Error("Couldn't save URL")
	}

	resp := Response{
		StatusCode: 201,
		Body:       s,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
