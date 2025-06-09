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

	if order == nil {
		log.Println("No active order found for customer and sender, creating a new order")

		channelID, err := uuid.Parse(constants.WHATSAPP_CHANNEL_ID)
		if err != nil {
			log.Printf("Error parsing channel ID: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}

		menu, err := ms_services.MountMenu(rdsClient, unit.ID)
		if err != nil {
			log.Printf("Error mounting menu for existing order: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}

		byteMenu, err := json.Marshal(menu)
		if err != nil {
			log.Printf("Error marshalling used menu: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}

		order = &gorm_entities.Order{
			ID:         uuid.New(),
			CustomerID: customer.ID,
			UnitID:     unit.ID,
			ChannelID:  channelID,
			Status:     constants.ORDER_STATUS_JUST_STARTED,
			UsedMenu:   byteMenu,
		}
	}

	order.LastMessageAt = time.Now()
	err = rdsClient.SaveOrder(order)
	if err != nil {
		log.Printf("Error saving order: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	menuAsMapArr, err := order.GetMenuAsMapArr()
	if err != nil {
		log.Printf("Error getting menu as map: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	productsFromOrder, err := rdsClient.GetOrderProducts(order.ID)
	if err != nil {
		log.Printf("Error fetching order products: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	var products = make([]map[string]string, len(productsFromOrder))
	for i, product := range productsFromOrder {
		products[i] = map[string]string{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.BRLPrice.StringFixed(2),
			"category":    product.Category,
		}
	}

	strOrderID := order.ID.String()
	n8nMessage := map[string]any{
		"message":               message,
		"menu":                  menuAsMapArr,
		"products":              products,
		"order_id":              strOrderID,
		"customer_id":           customer.ID.String(),
		"unit_id":               unit.ID.String(),
		"order_status":          order.Status,
		"order_last_message_at": order.LastMessageAt.Format(time.RFC3339),
	}

	services.NewN8NClient().Post("new_message", n8nMessage)

	eventBridgePayload := map[string]any{
		"order_id": strOrderID,
		"type":     "warn",
	}

	err = eventBridgeClient.PutEvent(
		ctx,
		fmt.Sprintf("order-warn-%s", strOrderID),
		15*time.Minute,
		eventBridgePayload,
	)
	if err != nil {
		log.Printf("Error creating schedule: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	eventBridgePayload["type"] = constants.ORDER_STATUS_TIMEOUT

	err = eventBridgeClient.PutEvent(
		ctx,
		fmt.Sprintf("order-timeout-%s", strOrderID),
		1*time.Hour,
		eventBridgePayload,
	)
	if err != nil {
		log.Printf("Error creating schedule: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	log.Println("Schedule created successfully for order:", n8nMessage["order_id"])

	return events.APIGatewayProxyResponse{StatusCode: 201}, nil
}
