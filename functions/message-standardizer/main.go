package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	schedulersvc "github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/copaerp/orders/functions/message-standardizer/handlers/whatsapp"
	"github.com/copaerp/orders/shared/repositories"
	"github.com/copaerp/orders/shared/services"
)

var schedulerClient *schedulersvc.Client
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

	return handler(ctx, request, rdsClient, schedulerClient)
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("erro ao carregar config da AWS: " + err.Error())
	}
	schedulerClient = schedulersvc.NewFromConfig(cfg)

	rdsClient, err = repositories.NewOrdersRDSClient()
	if err != nil {
		log.Printf("Error creating RDS client: %v", err)
		panic(err)
	}

	lambda.Start(handler)
}
