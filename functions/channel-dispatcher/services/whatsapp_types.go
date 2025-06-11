package services

type WhatsAppMessage struct {
	MessagingProduct string                      `json:"messaging_product"`
	To               string                      `json:"to"`
	Type             string                      `json:"type"`
	Text             *WhatsAppMessageText        `json:"text,omitempty"`
	Interactive      *WhatsAppInteractiveContent `json:"interactive,omitempty"`
}

// Text message
type WhatsAppMessageText struct {
	Body string `json:"body"`
}

// Interactive message
type WhatsAppInteractiveContent struct {
	Type   string                    `json:"type"`
	Body   WhatsAppTextBody          `json:"body"`
	Action WhatsAppInteractiveAction `json:"action"`
}

type WhatsAppTextBody struct {
	Text string `json:"text"`
}

type WhatsAppInteractiveAction struct {
	Button   *string                       `json:"button,omitempty"`
	Sections *[]WhatsAppInteractiveSection `json:"sections,omitempty"`
	Buttons  *[]WhatsAppInteractiveButtons `json:"buttons,omitempty"`
}

type WhatsAppInteractiveSection struct {
	Title string                   `json:"title"`
	Rows  []WhatsAppInteractiveRow `json:"rows"`
}

type WhatsAppInteractiveButtons struct {
	Type  string                 `json:"type"`
	Reply WhatsAppInteractiveRow `json:"reply"`
}

type WhatsAppInteractiveRow struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
}
