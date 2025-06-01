package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request any) error {

	log.Println("received request:")
	log.Printf("%v", request)

	return nil
}

func main() {
	lambda.Start(handler)
}
