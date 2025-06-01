package service

import (
	"strconv"

	"github.com/copaerp/orders/shared/repositories"
	"github.com/google/uuid"
)

func MountMenu(rdsClient *repositories.OrdersRDSClient, unitID uuid.UUID) (menu map[string]string, err error) {

	products, err := rdsClient.GetProductsByUnitID(unitID)
	if err != nil {
		return nil, err
	}

	menu = make(map[string]string, len(products))

	for i, product := range products {
		menu[product.ID.String()] = strconv.Itoa(i)
	}

	return menu, nil
}
