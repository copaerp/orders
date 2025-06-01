package repositories

import (
	"github.com/copaerp/orders/shared/entities"
	"github.com/google/uuid"
)

func (c *OrdersRDSClient) SaveOrder(order *entities.Order) error {
	result := c.DB.Save(order)
	return result.Error
}

func (c *OrdersRDSClient) GetActiveOrderByCustomerAndSender(customerID, senderID uuid.UUID) (*entities.Order, error) {
	var orders []entities.Order
	result := c.DB.
		Joins("Customer").
		Joins("Unit").
		Joins("Unit.WhatsappNumber").
		Where("Customer.id = ?", customerID).
		Where("Unit__WhatsappNumber.id = ?", senderID).
		Where("order.finished_at IS NULL").
		Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(orders) == 0 {
		return nil, nil
	}

	return &orders[0], nil
}
