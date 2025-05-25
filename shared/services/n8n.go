package services

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

var WORKFLOWS = map[string]string{
	"new_message": os.Getenv("n8n_webhook_url") + os.Getenv("new_message_workflow_id"),
}

type N8NClient struct{}

func NewN8NClient() *N8NClient {
	return &N8NClient{}
}

func (n *N8NClient) Post(workflow string, body map[string]any) ([]byte, error) {
	webhook, ok := WORKFLOWS[workflow]
	if !ok {
		log.Printf("Workflow não encontrado: %s", workflow)
		return nil, nil
	}

	log.Println("sending to ", webhook)

	jsonData, _ := json.Marshal(body)

	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Erro ao enviar para n8n: %v", err)
		return nil, err
	}

	resp_body, _ := io.ReadAll(resp.Body)
	log.Printf("Resposta do n8n: %s", resp_body)
	defer resp.Body.Close()

	return resp_body, nil
}
