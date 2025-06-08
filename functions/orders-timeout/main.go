package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/copaerp/orders/shared/constants"
	"github.com/copaerp/orders/shared/repositories"
	"github.com/copaerp/orders/shared/services"
)

var rdsClient *repositories.OrdersRDSClient

type Request struct {
	OrderID string `json:"order_id"`
	Channel string `json:"channel"`
	Type    string `json:"type"` // "timeout" | "warn"
}

func handler(ctx context.Context, request Request) error {

	log.Println("received request:")
	log.Printf("%v", request)

	if request.Type == constants.ORDER_STATUS_TIMEOUT {
		err := rdsClient.CloseOrder(request.OrderID)
		if err != nil {
			log.Printf("Error closing order: %v", err)
			return err
		}
	}

	services.NewN8NClient().Post("order_timeout", map[string]any{
		"order_id": request.OrderID,
		"channel":  request.Channel,
		"type":     request.Type,
	})

	return nil
}

func main() {

	var err error
	rdsClient, err = repositories.NewOrdersRDSClient()
	if err != nil {
		log.Printf("Error creating RDS client: %v", err)
		panic(err)
	}

	lambda.Start(handler)
}
