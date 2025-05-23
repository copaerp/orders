package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

type RequestMessage struct {
	Message string `json:"message"`
	Number  string `json:"number"`
	Channel string `json:"channel"`
}

func handler(request RequestMessage) error {

	log.Println("Eureka")
	log.Printf("message to be sent: %s, number: %s, channel: %s", request.Message, request.Number, request.Channel)

	return nil
}

func main() {
	lambda.Start(handler)
}
