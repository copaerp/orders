package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func (w *WhatsAppService) sendDefaultMessage(from string, whatsAppMessage WhatsAppMessage) error {
	jsonData, err := json.Marshal(whatsAppMessage)
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodyString := string(bodyBytes)

	log.Printf("WhatsApp API response: %s", bodyString)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("WhatsApp API returned non-success status code: %d with message %s", resp.StatusCode, bodyString)
	}

	return nil
}

func (w *WhatsAppService) SendMessage(from, to, message string) error {

	whatsAppMessage := WhatsAppMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "text",
		Text: &WhatsAppMessageText{
			Body: message,
		},
	}

	return w.sendDefaultMessage(from, whatsAppMessage)
}

func (w *WhatsAppService) SendButtonMessage(from, to, message string, button map[string]any) error {

	interactiveRows := []WhatsAppInteractiveRow{}
	for _, row := range button["rows"].([]any) {
		rowMap := row.(map[string]any)
		interactiveRows = append(interactiveRows, WhatsAppInteractiveRow{
			ID:          rowMap["id"].(string),
			Title:       rowMap["title"].(string),
			Description: rowMap["description"].(*string),
		})
	}

	whatsAppMessage := WhatsAppMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "interactive",
		Interactive: &WhatsAppInteractiveContent{
			Type: "list",
			Body: WhatsAppTextBody{
				Text: message,
			},
			Action: WhatsAppInteractiveAction{
				Button: button["text"].(*string),
				Sections: &[]WhatsAppInteractiveSection{
					{
						Title: "",
						Rows:  interactiveRows,
					},
				},
			},
		},
	}

	return w.sendDefaultMessage(from, whatsAppMessage)
}

func (w *WhatsAppService) SendButtonsMessage(from, to, message string, buttons map[string]any) error {

	interactiveRows := []WhatsAppInteractiveButtons{}
	for _, row := range buttons["rows"].([]any) {
		rowMap := row.(map[string]any)
		interactiveRows = append(interactiveRows, WhatsAppInteractiveButtons{
			Type: "reply",
			Reply: WhatsAppInteractiveRow{
				ID:    rowMap["id"].(string),
				Title: rowMap["title"].(string),
			},
		})
	}
	whatsAppMessage := WhatsAppMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "interactive",
		Interactive: &WhatsAppInteractiveContent{
			Type: "button",
			Body: WhatsAppTextBody{
				Text: message,
			},
			Action: WhatsAppInteractiveAction{
				Buttons: &interactiveRows,
			},
		},
	}

	return w.sendDefaultMessage(from, whatsAppMessage)
}
