package service

import (
	"strconv"

	"github.com/copaerp/orders/shared/repositories"
	"github.com/google/uuid"
)

func MountMenu(rdsClient *repositories.OrdersRDSClient, unitID uuid.UUID) (menu []map[string]string, err error) {

	products, err := rdsClient.GetProductsByUnitID(unitID)
	if err != nil {
		return nil, err
	}

	menu = make([]map[string]string, len(products))
	for i, product := range products {
		menu[i] = map[string]string{
			"id":          product.ID.String(),
			"index":       strconv.Itoa(i + 1), // Não estamos mais utilizando o índice devido aos testes com IA
			"name":        product.Name,
			"description": product.Description,
			"price":       product.BRLPrice.StringFixed(2),
			"category":    product.Category,
		}
	}

	return menu, nil
}
