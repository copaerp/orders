package repositories

import (
	"github.com/copaerp/orders/shared/entities"
)

func (c *OrdersRDSClient) GetUnitByWhatsappNumber(whatsappNumber string) (*entities.Unit, error) {
	var unit entities.Unit
	err := c.DB.
		Joins("WhatsappNumber").
		Where("WhatsappNumber.number = ?", whatsappNumber).
		First(&unit).Error

	if err != nil {
		return nil, err
	}

	return &unit, nil
}
