package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/copaerp/orders/shared/entities"
	"github.com/copaerp/orders/shared/repositories"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var rdsClient *repositories.OrdersRDSClient
var router *chi.Mux

// Middleware customizado para CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Lambda handler
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := request.Path
	path = strings.TrimPrefix(path, "/prod")

	req, err := http.NewRequest(request.HTTPMethod, path, bytes.NewBufferString(request.Body))
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

	router = chi.NewRouter()

	// Configurar CORS customizado
	router.Use(corsMiddleware)

	// Middleware básico
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

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

	// Salvar um pedido completo
	router.Post("/orders", func(w http.ResponseWriter, r *http.Request) {
		var order entities.Order

		// Decodificar o JSON do body
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			log.Printf("error decoding order JSON: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid_json"}`))
			return
		}

		// Salvar o pedido usando o RDS client
		if err := rdsClient.SaveOrder(&order); err != nil {
			log.Printf("error saving order: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"save_failed"}`))
			return
		}

		// Retornar sucesso com o pedido salvo
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	})
	router.Get("/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Fetching order with ID: " + id))
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
