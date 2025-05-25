package repositories

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var connString string = os.Getenv("orders_db_connection_url")

type OrdersRDSClient struct {
	DB *gorm.DB
}

func NewOrdersRDSClient() *OrdersRDSClient {
	db, err := gorm.Open(mysql.Open(connString), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		panic("failed to connect database")
	}

	return &OrdersRDSClient{
		DB: db,
	}
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

func (c *OrdersRDSClient) Execute(query string, args ...any) (int64, error) {
	result := c.DB.Exec(query, args...)
	if result.Error != nil {
		log.Printf("Error executing query: %v", result.Error)
		return 0, result.Error
	}

	rowsAffected := result.RowsAffected
	return rowsAffected, nil
}
