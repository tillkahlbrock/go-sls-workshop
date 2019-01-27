package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	name, ok := request.QueryStringParameters["name"]
	if !ok {
		return Response{StatusCode: 501}, fmt.Errorf("missing required name parameter")
	}

	body := fmt.Sprintf("Hello %s!", name)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            body,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
