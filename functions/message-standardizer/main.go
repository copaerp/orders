package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const verifyToken = "your_verify_token"

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	requestJSON, _ := json.MarshalIndent(request, "", "  ")
	log.Printf("Received request: %s", requestJSON)

	if request.HTTPMethod == "GET" {
		mode := request.QueryStringParameters["hub.mode"]
		token := request.QueryStringParameters["hub.verify_token"]
		challenge := request.QueryStringParameters["hub.challenge"]

		if mode == "subscribe" && token == verifyToken {
			log.Println("Webhook verificado com sucesso")
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       challenge,
			}, nil
		}

		log.Println("Token de verificação inválido")
		return events.APIGatewayProxyResponse{StatusCode: 403}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
