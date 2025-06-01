package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type WhatsAppService struct {
	Token string
}

func NewWhatsAppService(token string) *WhatsAppService {
	return &WhatsAppService{
		Token: token,
	}
}

type WhatsAppMessageText struct {
	Body string `json:"body"`
}

type WhatsAppMessage struct {
	MessagingProduct string              `json:"messaging_product"`
	To               string              `json:"to"`
	Type             string              `json:"type"`
	Text             WhatsAppMessageText `json:"text"`
}

func (w *WhatsAppService) SendMessage(from, to, message string) error {

	whatsappMessage := WhatsAppMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "text",
		Text: WhatsAppMessageText{
			Body: message,
		},
	}

	jsonData, err := json.Marshal(whatsappMessage)
	if err != nil {
		return err
	}

	fullUrl := os.Getenv("whatsapp_api_url") + from + "/messages"

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+w.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Printf("WhatsApp API response: %v", resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("WhatsApp API returned non-success status code: %d with message %v", resp.StatusCode, err)
	}

	return nil
}
