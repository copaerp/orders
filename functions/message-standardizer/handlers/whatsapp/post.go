package whatsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/copaerp/orders/functions/message-standardizer/entities"
	ms_services "github.com/copaerp/orders/functions/message-standardizer/services"
	"github.com/copaerp/orders/shared/constants"
	gorm_entities "github.com/copaerp/orders/shared/entities"
	"github.com/copaerp/orders/shared/repositories"
	"github.com/copaerp/orders/shared/services"
	"github.com/google/uuid"
)

func Post(ctx context.Context, request events.APIGatewayProxyRequest, rdsClient *repositories.OrdersRDSClient, eventBridgeClient *repositories.EventBridgeClient) (events.APIGatewayProxyResponse, error) {

	log.Println("request body: ", request.Body)

	var payload entities.WhatsAppMessage
	err := json.Unmarshal([]byte(request.Body), &payload)
	if err != nil {
		log.Printf("Erro ao decodificar JSON: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}

	var customerNumber, message, senderNumber string

	if len(payload.Entry) > 0 && len(payload.Entry[0].Changes) > 0 {
		value := payload.Entry[0].Changes[0].Value

		if value.Statuses != nil {
			log.Println("Received a status update, ignoring the message")
			return events.APIGatewayProxyResponse{StatusCode: 200}, nil
		}

		msgs := value.Messages
		if len(msgs) > 0 {
			customerNumber = msgs[0].From
			message = msgs[0].Text.Body
			senderNumber = value.Metadata.DisplayPhoneNumber
		}
	}

	customer, err := rdsClient.GetCustomerByNumber(customerNumber)
	if err != nil {
		log.Printf("Error fetching customer: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 404}, nil
	}

	unit, err := rdsClient.GetUnitByWhatsappNumber(senderNumber)
	if err != nil {
		log.Printf("Error fetching unit: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	order, err := rdsClient.GetActiveOrderByCustomerAndSender(customer.ID, unit.ID)
	if err != nil {
		log.Printf("Error fetching active orders: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	var menu map[string]map[string]string
	if order == nil || order.UsedMenu == nil {
		menu, err = ms_services.MountMenu(rdsClient, order.UnitID)
		if err != nil {
			log.Printf("Error mounting menu for existing order: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}
	} else {
		err := json.Unmarshal(order.UsedMenu, &menu)
		if err != nil {
			log.Printf("Error unmarshalling used menu: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}
	}

	n8nMessage := map[string]any{
		"number":         customerNumber,
		"message":        message,
		"channel":        "whatsapp",
		"sender":         senderNumber,
		"meta_number_id": unit.WhatsappNumber.MetaNumberID,
		"menu":           menu,
	}

	if order == nil {
		log.Println("No active order found for customer and sender, creating a new order")

		channelID, err := uuid.Parse(constants.WHATSAPP_CHANNEL_ID)
		if err != nil {
			log.Printf("Error parsing channel ID: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}

		order = &gorm_entities.Order{
			ID:         uuid.New(),
			CustomerID: customer.ID,
			UnitID:     unit.ID,
			ChannelID:  channelID,
			Status:     constants.ORDER_STATUS_JUST_STARTED,
		}
	} else {

		productsFromOrder, err := rdsClient.GetOrderProducts(order.ID)
		if err != nil {
			log.Printf("Error fetching order products: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}

		var products = make(map[string]map[string]string, len(productsFromOrder))
		for _, product := range productsFromOrder {
			products[menu[product.ID.String()]["index"]] = map[string]string{
				"name":        product.Name,
				"description": product.Description,
				"price":       product.BRLPrice.StringFixed(4),
				"category":    product.Category,
			}
		}

	}

	order.LastMessageAt = time.Now()
	if order.UsedMenu == nil {
		order.UsedMenu, err = json.Marshal(menu)
		if err != nil {
			log.Printf("Error marshalling used menu: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}
	}
	err = rdsClient.SaveOrder(order)
	if err != nil {
		log.Printf("Error saving order: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	n8nMessage["order_id"] = order.ID.String()
	n8nMessage["customer_id"] = order.CustomerID.String()
	n8nMessage["unit_id"] = order.UnitID.String()
	n8nMessage["channel_id"] = order.ChannelID.String()
	n8nMessage["order_status"] = order.Status
	n8nMessage["order_last_message_at"] = order.LastMessageAt.Format(time.RFC3339)

	services.NewN8NClient().Post("new_message", n8nMessage)

	eventBridgePayload := map[string]any{
		"order_id": order.ID.String(),
		"type":     "warn",
		"channel":  "whatsapp",
	}

	err = eventBridgeClient.PutEvent(
		ctx,
		fmt.Sprintf("order-%s", order.ID.String()),
		10*time.Minute,
		eventBridgePayload,
	)
	if err != nil {
		log.Printf("Error creating schedule: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	eventBridgePayload["type"] = "timeout"

	err = eventBridgeClient.PutEvent(
		ctx,
		fmt.Sprintf("order-%s", order.ID.String()),
		1*time.Hour,
		eventBridgePayload,
	)
	if err != nil {
		log.Printf("Error creating schedule: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	log.Println("Schedule created successfully for order:", order.ID.String())

	return events.APIGatewayProxyResponse{StatusCode: 201}, nil
}
