package entities

type WhatsAppMessage struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Field string `json:"field"`
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					From      string `json:"from"`
					ID        string `json:"id"`
					Timestamp string `json:"timestamp"`
					Text      struct {
						Body string `json:"body"`
					} `json:"text"`
					Interactive struct {
						Type      string `json:"type"`
						ListReply struct {
							ID          string `json:"id"`
							Title       string `json:"title"`
							Description string `json:"description,omitempty"`
						} `json:"list_reply"`
						ButtonReply struct {
							ID    string `json:"id"`
							Title string `json:"title"`
						}
					}
					Type string `json:"type"`
				} `json:"messages"`
				Statuses []struct {
					ID           string `json:"id"`
					Status       string `json:"status"`
					Timestamp    string `json:"timestamp"`
					RecipientID  string `json:"recipient_id"`
					Conversation struct {
						ID                  string `json:"id"`
						ExpirationTimestamp string `json:"expiration_timestamp"`
						Origin              struct {
							Type string `json:"type"`
						} `json:"origin"`
					}
					Pricing struct {
						Billable     bool   `json:"billable"`
						PricingModel string `json:"pricing_model"`
						Category     string `json:"category"`
					} `json:"pricing"`
				}
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}
