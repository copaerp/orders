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
	OrderID string         `json:"order_id"`
	Message map[string]any `json:"message"`
}

var whatsappToken string
var rdsClient *repositories.OrdersRDSClient

func handler(ctx context.Context, request RequestMessage) error {

	log.Printf("message to be sent: %s, order_id: %s", request.Message, request.OrderID)

	order, err := rdsClient.GetOrderByID(request.OrderID)
	if err != nil {
		log.Printf("Error fetching order: %v", err)
		return fmt.Errorf("error fetching order: %v", err)
	}

	channel := order.Channel.Name
	customerNumber := order.Customer.Phone
	senderMetaNumberID := order.Unit.WhatsappNumber.MetaNumberID

	log.Printf("message to be sent: %s, number: %s, channel: %s, sender: %s", request.Message, customerNumber, channel, senderMetaNumberID)

	switch channel {
	case "WhatsApp":
		whatsappClient := services.NewWhatsAppService(whatsappToken)

		messageMainText := request.Message["main_text"].(string)

		if request.Message["button"] != nil {
			item := request.Message["button"].(map[string]any)
			return whatsappClient.SendButtonMessage(senderMetaNumberID, customerNumber, messageMainText, item)
		}

		if request.Message["buttons"] != nil {
			item := request.Message["buttons"].(map[string]any)
			return whatsappClient.SendButtonsMessage(senderMetaNumberID, customerNumber, messageMainText, item)
		}

		return whatsappClient.SendMessage(senderMetaNumberID, customerNumber, messageMainText)
	default:
		log.Printf("Channel %s not supported", channel)
		return fmt.Errorf("channel %s not supported", channel)
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
