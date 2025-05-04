package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler() error {

	log.Println("Eureka")

	return nil
}

func main() {
	lambda.Start(handler)
}
