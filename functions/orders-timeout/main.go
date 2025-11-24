package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/copaerp/orders/shared/constants"
	"github.com/copaerp/orders/shared/repositories"
	"github.com/copaerp/orders/shared/services"
)

var rdsClient *repositories.OrdersRDSClient

type Request struct {
	OrderID string `json:"order_id"`
	Type    string `json:"type"` // "timeout" | "warn"
}

func handler(ctx context.Context, request Request) error {

	log.Println("received request:")
	log.Printf("%v", request)

	order, err := rdsClient.GetOrder(request.OrderID)
	if err != nil {
		log.Printf("Error fetching order: %v", err)
		return err
	}

	if order.FinishedAt != nil || order.CanceledAt != nil {
		log.Printf("Order %s is already finished or canceled, skipping timeout processing", request.OrderID)
		return nil
	}

	if request.Type == constants.ORDER_STATUS_TIMEOUT {

		canceledAt := time.Now()
		order.CanceledAt = &canceledAt
		err = rdsClient.SaveOrder(&order)
		if err != nil {
			log.Printf("Error finishing order: %v", err)
			return err
		}
	}

	services.NewN8NClient().Post("order_timeout", map[string]any{
		"order_id":           request.OrderID,
		"type":               request.Type,
		"escalated_to_human": order.EscalatedToHuman,
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
