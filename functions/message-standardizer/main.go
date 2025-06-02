package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/copaerp/orders/functions/message-standardizer/handlers/whatsapp"
	"github.com/copaerp/orders/shared/repositories"
	"github.com/copaerp/orders/shared/services"
)

var eventBridgeClient *repositories.EventBridgeClient
var rdsClient *repositories.OrdersRDSClient

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	router := services.NewRouter()

	router.Add("GET", "/whatsapp", whatsapp.Get)
	router.Add("POST", "/whatsapp", whatsapp.Post)

	handler, found := router.Find(request.HTTPMethod, request.Path)
	if !found {
		log.Println("Route not found:", request.HTTPMethod, request.Path)
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       fmt.Sprintf("Route not found: %s %s", request.HTTPMethod, request.Path),
		}, nil
	}

	return handler(ctx, request, rdsClient, eventBridgeClient)
}

func main() {

	var err error
	eventBridgeClient, err = repositories.NewEventBridgeClient()
	if err != nil {
		log.Printf("Error creating EventBridge client: %v", err)
		panic(err)
	}

	rdsClient, err = repositories.NewOrdersRDSClient()
	if err != nil {
		log.Printf("Error creating RDS client: %v", err)
		panic(err)
	}

	lambda.Start(handler)
}
