package repositories

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/copaerp/orders/shared/constants"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// var dbPassword string = os.Getenv("orders_db_password")
var dbUser string = os.Getenv("orders_db_username")
var dbEndpoint string = os.Getenv("orders_db_endpoint")
var dbName string = os.Getenv("orders_db_name")

type OrdersRDSClient struct {
	DB *gorm.DB
}

func NewOrdersRDSClient(ctx context.Context) (*OrdersRDSClient, error) {

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("configuration error: " + err.Error())
	}

	authenticationToken, err := auth.BuildAuthToken(
		ctx, dbEndpoint, constants.AWS_REGION, dbUser, cfg.Credentials)
	if err != nil {
		panic("failed to create authentication token: " + err.Error())
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true&allowCleartextPasswords=true",
		dbUser, authenticationToken, dbEndpoint, dbName,
	)

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

func (c *OrdersRDSClient) Execute(query string, args ...any) (int64, error) {
	result := c.DB.Exec(query, args...)
	if result.Error != nil {
		log.Printf("Error executing query: %v", result.Error)
		return 0, result.Error
	}

	rowsAffected := result.RowsAffected
	return rowsAffected, nil
}
