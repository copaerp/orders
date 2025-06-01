package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/copaerp/orders/functions/channel-dispatcher/services"
	"github.com/copaerp/orders/shared/entities"
	"github.com/copaerp/orders/shared/repositories"
)

type RequestMessage struct {
	Message string `json:"message"`
	Number  string `json:"number"`
	Channel string `json:"channel"`
}

var whatsappToken string

func handler(ctx context.Context, request RequestMessage) error {

	if request.Channel != "whatsapp" {
		log.Println("mock testing RDS connection")

		rdsClient, err := repositories.NewOrdersRDSClient(ctx)
		if err != nil {
			log.Printf("Error creating RDS client: %v", err)
			return err
		}

		log.Println(rdsClient.GetDB().Name())

		var channels []entities.Channel

		err = rdsClient.GetDB().Find(&channels).Error
		if err != nil {
			log.Printf("Erro ao buscar canais: %v", err)
			return err
		}

		log.Printf("Found %d channels", len(channels))
		for _, ch := range channels {
			log.Printf("Canal: ID=%s, Name=%s", ch.ID.String(), ch.Name)
		}

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
