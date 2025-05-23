package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/copaerp/orders/functions/message-standardizer/handlers/whatsapp"
	"github.com/copaerp/orders/functions/message-standardizer/services"
)

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

	return handler(ctx, request)
}

func main() {
	lambda.Start(handler)
}
