package whatsapp

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

const verifyToken = "your_verify_token"

func Get(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	mode := request.QueryStringParameters["hub.mode"]
	token := request.QueryStringParameters["hub.verify_token"]
	challenge := request.QueryStringParameters["hub.challenge"]

	if mode == "subscribe" && token == verifyToken {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       challenge,
		}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 403}, nil
}
