package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const verifyToken = "your_verify_token"
const n8nWebhookURL = "https://n8n.copaerp.site/webhook/aba98742-debe-4f62-a283-55519635318b"

type WhatsAppMessage struct {
	Entry []struct {
		Changes []struct {
			Value struct {
				Messages []struct {
					From string `json:"from"`
					Text struct {
						Body string `json:"body"`
					} `json:"text"`
				} `json:"messages"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "GET" {
		mode := request.QueryStringParameters["hub.mode"]
		token := request.QueryStringParameters["hub.verify_token"]
		challenge := request.QueryStringParameters["hub.challenge"]

		if mode == "subscribe" && token == verifyToken {
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       challenge,
			}, nil
		}
		return events.APIGatewayProxyResponse{StatusCode: 403}, nil
	}

	if request.HTTPMethod == "POST" {
		var payload WhatsAppMessage
		err := json.Unmarshal([]byte(request.Body), &payload)
		if err != nil {
			log.Printf("Erro ao decodificar JSON: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 400}, nil
		}

		if len(payload.Entry) > 0 && len(payload.Entry[0].Changes) > 0 {
			msgs := payload.Entry[0].Changes[0].Value.Messages
			if len(msgs) > 0 {
				number := msgs[0].From
				message := msgs[0].Text.Body

				outgoing := map[string]string{
					"number":  number,
					"message": message,
				}
				jsonData, _ := json.Marshal(outgoing)

				resp, err := http.Post(n8nWebhookURL, "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					log.Printf("Erro ao enviar para n8n: %v", err)
				} else {
					body, _ := ioutil.ReadAll(resp.Body)
					log.Printf("Resposta do n8n: %s", body)
					resp.Body.Close()
				}
			}
		}

		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 405}, nil
}

func main() {
	lambda.Start(handler)
}
