package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/copaerp/orders/shared/repositories"
	"github.com/go-chi/chi/v5"
)

var rdsClient *repositories.OrdersRDSClient
var router *chi.Mux

// Lambda handler
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Converter APIGatewayProxyRequest em http.Request
	req, err := http.NewRequest(request.HTTPMethod, request.Path, bytes.NewBufferString(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	// Headers
	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()

	// Passar pelo router chi
	router.ServeHTTP(rec, req)

	// Converter resposta
	resp := rec.Result()
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	return events.APIGatewayProxyResponse{
		StatusCode: resp.StatusCode,
		Headers:    flattenHeaders(resp.Header),
		Body:       string(bodyBytes),
	}, nil
}

// Helper para converter http.Header em map[string]string
func flattenHeaders(h http.Header) map[string]string {
	out := make(map[string]string)
	for k, v := range h {
		if len(v) > 0 {
			out[k] = v[0]
		}
	}
	return out
}

func main() {
	var err error
	rdsClient, err = repositories.NewOrdersRDSClient()
	if err != nil {
		log.Printf("Error creating RDS client: %v", err)
		panic(err)
	}

	// Criar router e registrar rotas
	router = chi.NewRouter()

	// Listar todos os pedidos
	router.Get("/orders", func(w http.ResponseWriter, r *http.Request) {
		orders, err := rdsClient.ListOrders()
		if err != nil {
			log.Printf("error listing orders: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"internal_error"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	})
	router.Get("/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Fetching order with ID: " + id))
	})

	router.Post("/orders", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Order created with body: " + string(body)))
	})

	// Buscar produto por ID
	router.Get("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		product, err := rdsClient.GetProductByID(id)
		if err != nil {
			log.Printf("error fetching product %s: %v", id, err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"not_found"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	})

	lambda.Start(handler)
}
