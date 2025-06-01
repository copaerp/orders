package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/copaerp/orders/functions/channel-dispatcher/services"
)

type RequestMessage struct {
	Message string `json:"message"`
	Number  string `json:"number"`
	Channel string `json:"channel"`
	Sender  string `json:"sender"`
}

var whatsappToken string

func handler(ctx context.Context, request RequestMessage) error {

	log.Printf("message to be sent: %s, number: %s, channel: %s, sender: %s", request.Message, request.Number, request.Channel, request.Sender)

	switch request.Channel {
	case "whatsapp":
		whatsappClient := services.NewWhatsAppService(whatsappToken)
		err := whatsappClient.SendMessage(request.Number, request.Message)

		return err
	default:
		log.Printf("Channel %s not supported", request.Channel)
		return fmt.Errorf("channel %s not supported", request.Channel)
	}
}

func main() {
	lambda.Start(handler)
}
