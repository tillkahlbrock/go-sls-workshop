package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

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

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	s, ok := request.PathParameters["short"]
	if !ok {
		return Response{StatusCode: 501}, fmt.Errorf("missing required short parameter")
	}

	// Create a new AWS session and fail immediately on error
	sess := session.Must(session.NewSession())
	// Create the DynamoDB client
	dynamodbclient := dynamodb.New(sess)
	u, err := dynamodbclient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("DYNAMO_DB_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"short_url": {
				S: aws.String(s),
			},
		},
	})
	if err != nil {
		logrus.WithField("error", err).Error("failed to fetch long url")
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
	lambda.Start(Handler)
}
