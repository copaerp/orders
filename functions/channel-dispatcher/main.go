package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/copaerp/orders/functions/channel-dispatcher/services"
)

type RequestMessage struct {
	Message string `json:"message"`
	Number  string `json:"number"`
	Channel string `json:"channel"`
}

var whatsappToken string

func handler(request RequestMessage) error {

	if request.Channel == "dummy" {
		log.Printf("Dummy channel selected, message: %s, number: %s", request.Message, request.Number)

		// lets test the db integration

		log.Println(os.Getenv("orders_db_connection_url"))

		return nil
	}

	log.Printf("message to be sent: %s, number: %s, channel: %s", request.Message, request.Number, request.Channel)

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
