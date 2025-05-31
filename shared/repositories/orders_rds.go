package repositories

import (
	"context"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbUser string = os.Getenv("orders_db_username")
var dbPassword string = os.Getenv("orders_db_password")
var dbEndpoint string = os.Getenv("orders_db_endpoint")
var dbName string = os.Getenv("orders_db_name")

type OrdersRDSClient struct {
	DB *gorm.DB
}

func NewOrdersRDSClient(ctx context.Context) (*OrdersRDSClient, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbEndpoint, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, err
	}

	return &OrdersRDSClient{
		DB: db,
	}, nil
}

func (c *OrdersRDSClient) GetDB() *gorm.DB {
	return c.DB
}

func (c *OrdersRDSClient) Query(query string, args ...any) (*gorm.DB, error) {
	result := c.DB.Raw(query, args...)
	if result.Error != nil {
		log.Printf("Error executing query: %v", result.Error)
		return nil, result.Error
	}

	return result, nil
}
