package repositories

import (
	"github.com/copaerp/orders/shared/entities"
)

func (c *OrdersRDSClient) GetCustomerByNumber(phoneNumber string) (*entities.Customer, error) {
	var customer entities.Customer
	err := c.DB.
		Where("phone = ?", phoneNumber).
		First(&customer).Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}
