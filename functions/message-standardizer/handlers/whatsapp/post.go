package whatsapp

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/copaerp/orders/functions/message-standardizer/entities"
	"github.com/copaerp/orders/shared/services"
)

func Post(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Println("all whatsapp raw headers: ", request.Headers)
	log.Println("all whatsapp raw data: ", request.Body)
	log.Println("all whatsapp raw query: ", request.QueryStringParameters)

	var payload entities.WhatsAppMessage
	err := json.Unmarshal([]byte(request.Body), &payload)
	if err != nil {
		log.Printf("Erro ao decodificar JSON: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}

	if len(payload.Entry) > 0 && len(payload.Entry[0].Changes) > 0 {
		msgs := payload.Entry[0].Changes[0].Value.Messages
		if len(msgs) > 0 {
			services.NewN8NClient().Post("new_message", map[string]any{
				"number":  msgs[0].From,
				"message": msgs[0].Text.Body,
				"channel": "whatsapp",
			})
		}
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
