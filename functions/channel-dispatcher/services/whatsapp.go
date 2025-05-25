package services

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type WhatsAppMessageTemplateLanguage struct {
	Code string `json:"code"`
}

type WhatsAppMessageTemplate struct {
	Name     string `json:"name"`
	Language WhatsAppMessageTemplateLanguage
}

type WhatsAppMessage struct {
	MessagingProduct string `json:"messaging_product"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Template         WhatsAppMessageTemplate
}

func (w *WhatsAppService) SendMessage(to, message string) error {

	whatsappMessage := WhatsAppMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template: WhatsAppMessageTemplate{
			Name: "hello_world",
			Language: WhatsAppMessageTemplateLanguage{
				Code: "en_US",
			},
		},
	}

	jsonData, err := json.Marshal(whatsappMessage)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", os.Getenv("whatsapp_api_url"), bytes.NewBuffer(jsonData))
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("WhatsApp API returned non-success status code: %d", resp.StatusCode)
	}

	return nil
}
