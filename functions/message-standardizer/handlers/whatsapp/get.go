package whatsapp

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	schedulersvc "github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/copaerp/orders/shared/repositories"
)

func Get(ctx context.Context, request events.APIGatewayProxyRequest, rdsClient *repositories.OrdersRDSClient, schedulerClient *schedulersvc.Client) (events.APIGatewayProxyResponse, error) {
	mode := request.QueryStringParameters["hub.mode"]
	token := request.QueryStringParameters["hub.verify_token"]
	challenge := request.QueryStringParameters["hub.challenge"]

	if mode == "subscribe" && token == os.Getenv("whatsapp_verify_token") {
		return events.APIGatewayProxyResponse{StatusCode: 200, Body: challenge}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 403}, nil
}
