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

// GetProductByID returns a single product with its associations.
func (c *OrdersRDSClient) GetProductByID(productID string) (*entities.Product, error) {
	var product entities.Product
	err := c.DB.
		Preload("Business").
		Preload("ProductsInUnits").
		Preload("ProductsOrders").
		First(&product, "id = ?", productID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProductByIDAndBusinessID returns a single product filtering by business ID.
func (c *OrdersRDSClient) GetProductByIDAndBusinessID(productID string, businessID uuid.UUID) (*entities.Product, error) {
	var product entities.Product
	err := c.DB.
		Preload("Business").
		Preload("ProductsInUnits").
		Preload("ProductsOrders").
		Where("id = ? AND business_id = ?", productID, businessID).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetMenuByBusinessID returns all products (menu) for a specific business.
func (c *OrdersRDSClient) GetMenuByBusinessID(businessID uuid.UUID) ([]entities.Product, error) {
	var products []entities.Product
	err := c.DB.
		Preload("Business").
		Where("business_id = ?", businessID).
		Order("category ASC, name ASC").
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// GetProductByIDAndUnitID returns a single product filtering by unit ID through product_in_unit.
func (c *OrdersRDSClient) GetProductByIDAndUnitID(productID string, unitID uuid.UUID) (*entities.Product, error) {
	var product entities.Product
	err := c.DB.
		Preload("Business").
		Preload("ProductsInUnits").
		Preload("ProductsOrders").
		Joins("JOIN product_in_unit piu ON piu.product_id = product.id").
		Where("product.id = ? AND piu.unit_id = ?", productID, unitID).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetMenuByUnitID returns all products (menu) for a specific unit through product_in_unit.
func (c *OrdersRDSClient) GetMenuByUnitID(unitID uuid.UUID) ([]entities.Product, error) {
	var products []entities.Product
	err := c.DB.
		Preload("Business").
		Joins("JOIN product_in_unit piu ON piu.product_id = product.id").
		Where("piu.unit_id = ?", unitID).
		Order("product.category ASC, product.name ASC").
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
