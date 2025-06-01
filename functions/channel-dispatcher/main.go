package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/copaerp/orders/functions/channel-dispatcher/services"
	"github.com/copaerp/orders/shared/repositories"
)

type RequestMessage struct {
	OrderID string `json:"order_id"`
	Message string `json:"message"`
	Channel string `json:"channel"`
}

var whatsappToken string
var rdsClient *repositories.OrdersRDSClient

func handler(ctx context.Context, request RequestMessage) error {

	log.Printf("message to be sent: %s, channel: %s, order_id: %s", request.Message, request.Channel, request.OrderID)

	order, err := rdsClient.GetOrderByID(request.OrderID)
	if err != nil {
		log.Printf("Error fetching order: %v", err)
		return fmt.Errorf("error fetching order: %v", err)
	}

	customerNumber := order.Customer.Phone
	senderMetaNumberID := order.Unit.WhatsappNumber.MetaNumberID

	log.Printf("message to be sent: %s, number: %s, channel: %s, sender: %s", request.Message, customerNumber, request.Channel, senderMetaNumberID)

	switch request.Channel {
	case "whatsapp":
		whatsappClient := services.NewWhatsAppService(whatsappToken)
		err := whatsappClient.SendMessage(senderMetaNumberID, customerNumber, request.Message)

		return err
	default:
		log.Printf("Channel %s not supported", request.Channel)
		return fmt.Errorf("channel %s not supported", request.Channel)
	}
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
