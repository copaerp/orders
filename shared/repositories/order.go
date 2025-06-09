package repositories

import (
	"log"

	"github.com/copaerp/orders/shared/entities"
	"github.com/google/uuid"
)

func (c *OrdersRDSClient) SaveOrder(order *entities.Order) error {
	result := c.DB.Save(order)
	return result.Error
}

func (c *OrdersRDSClient) GetOrderByID(orderID string) (*entities.Order, error) {
	var order entities.Order
	result := c.DB.
		Joins("Customer").
		Joins("Unit").
		Joins("Channel").
		Joins("Unit.WhatsappNumber").
		Where("order.id = ?", orderID).
		First(&order)

	if result.Error != nil {
		return nil, result.Error
	}

	return &order, nil
}

func (c *OrdersRDSClient) GetActiveOrderByCustomerAndSender(customerID, unitID uuid.UUID) (*entities.Order, error) {

	log.Printf("GetActiveOrderByCustomerAndSender: customerID: %s, unitID: %s", customerID, unitID)

	var orders []entities.Order
	result := c.DB.
		Joins("Customer").
		Joins("Unit").
		Where("Customer.id = ?", customerID).
		Where("Unit.id = ?", unitID).
		Where("order.finished_at IS NULL").
		Where("order.canceled_at IS NULL").
		Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(orders) == 0 {
		return nil, nil
	}

	return &orders[0], nil
}

func (c *OrdersRDSClient) GetOrder(orderID string) (entities.Order, error) {
	var order entities.Order
	result := c.DB.First(&order, "id = ?", orderID)
	if result.Error != nil {
		return entities.Order{}, result.Error
	}

	return order, result.Error
}
