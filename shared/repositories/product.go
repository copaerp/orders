package repositories

import (
	"github.com/copaerp/orders/shared/entities"
	"github.com/google/uuid"
)

func (c *OrdersRDSClient) GetProductsByUnitID(unitID uuid.UUID) ([]entities.Product, error) {
	var products []entities.Product
	err := c.DB.
		Joins("JOIN product_in_unit piu ON piu.product_id = product.id").
		Where("piu.unit_id = ?", unitID).
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (c *OrdersRDSClient) GetOrderProducts(orderID uuid.UUID) ([]entities.Product, error) {
	var products []entities.Product
	err := c.DB.
		Joins("JOIN product_order ON product_order.product_id = product.id").
		Where("product_order.order_id = ?", orderID).
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}
