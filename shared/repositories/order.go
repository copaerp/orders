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
	var order entities.Order
	err := c.DB.
		Joins("Customer").
		Joins("Unit").
		Joins("Unit.WhatsappNumber").
		Where("customer.id = ?", customerID.String()).
		Where("whatsapp_number.id = ?", senderID.String()).
		Where("order.finished_at IS NULL").
		First(&order).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (c *OrdersRDSClient) GetActiveOrderByCustomerAndSenderNumbers(customerNumber, senderNumber string) ([]entities.Order, error) {

	var orders []entities.Order
	err := c.DB.
		Joins("Customer").
		Joins("Unit").
		Joins("Unit.WhatsappNumber").
		Where("customer.phone = ?", customerNumber).
		Where("whatsapp_number.number = ?", senderNumber).
		Where("order.finished_at IS NULL").
		First(&orders).Error

	return orders, err
}
