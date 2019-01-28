package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type Item struct {
	Url      string
	ShortUrl string `json:"short_url"`
}

type HandlerConfig struct {
	table string
	c     dynamodbiface.DynamoDBAPI
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func (hc *HandlerConfig) Handler(request events.APIGatewayProxyRequest) (Response, error) {
	s, ok := request.PathParameters["short"]
	if !ok {
		return Response{StatusCode: 501}, fmt.Errorf("missing required short parameter")
	}

	u, err := hc.c.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(hc.table),
		Key: map[string]*dynamodb.AttributeValue{
			"short_url": {
				S: aws.String(s),
			},
		},
	})
	if err != nil {
		logrus.WithField("error", err).Error("failed to fetch long url")
	}

	if len(u.Item) <= 1 {
		return Response{}, fmt.Errorf("borken")
	}

	item := Item{}

	err = dynamodbattribute.UnmarshalMap(u.Item, &item)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	resp := Response{
		StatusCode: 302,
		Headers: map[string]string{
			"Location": item.Url,
		},
	}

	return resp, nil
}

func main() {
	// Create a new AWS session and fail immediately on error
	sess := session.Must(session.NewSession())
	// Create the DynamoDB client
	dynamodbclient := dynamodb.New(sess)
	hc := HandlerConfig{c: dynamodbclient, table: os.Getenv("DYNAMO_DB_TABLE")}
	lambda.Start(hc.Handler)
}
