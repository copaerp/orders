package repositories

import (
	"fmt"
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

// ListOrders returns all orders with basic associations preloaded.
func (c *OrdersRDSClient) ListOrders() ([]entities.Order, error) {
	var orders []entities.Order
	result := c.DB.
		Preload("Customer").
		Preload("Unit").
		Preload("Channel").
		Where("order.finished_at IS NOT NULL").
		Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

// ListOrdersByBusinessID returns all orders for a specific business with basic associations preloaded.
func (c *OrdersRDSClient) ListOrdersByBusinessID(businessID uuid.UUID) ([]entities.Order, error) {
	var orders []entities.Order
	result := c.DB.
		Preload("Customer").
		Preload("Unit").
		Preload("Channel").
		Joins("JOIN unit ON unit.id = order.unit_id").
		Where("unit.business_id = ?", businessID).
		Where("order.finished_at IS NOT NULL").
		Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

// GetOrderByIDAndBusinessID returns an order by ID filtering by business ID.
func (c *OrdersRDSClient) GetOrderByIDAndBusinessID(orderID string, businessID uuid.UUID) (*entities.Order, error) {
	var order entities.Order
	result := c.DB.
		Joins("Customer").
		Joins("Unit").
		Joins("Channel").
		Joins("Unit.WhatsappNumber").
		Joins("JOIN unit ON unit.id = order.unit_id").
		Where("order.id = ?", orderID).
		Where("unit.business_id = ?", businessID).
		First(&order)

	if result.Error != nil {
		return nil, result.Error
	}

	return &order, nil
}

// ValidateUnitBelongsToBusiness validates if a unit belongs to a specific business.
func (c *OrdersRDSClient) ValidateUnitBelongsToBusiness(unitID, businessID uuid.UUID) error {
	var count int64
	result := c.DB.Model(&entities.Unit{}).
		Where("id = ? AND business_id = ?", unitID, businessID).
		Count(&count)

	if result.Error != nil {
		return result.Error
	}

	if count == 0 {
		return fmt.Errorf("unit does not belong to business")
	}

	return nil
}
