package main

import (
	"app-server/internal/api"
	"app-server/internal/config"
	"app-server/internal/functions"
	"app-server/pkg/postgres"
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	conn, err := postgres.New(cfg.Postgres)
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			panic(err)
		}
	}(conn, context.Background())

	if err != nil {
		panic(err)
	}
	fmt.Println(conn.Ping(ctx))

	functions.InitDB(conn)
	router := mux.NewRouter()
	api.Create(router)

	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	logs := handlers.CORS(headers, methods, origins)(router)

	log.Printf("Server started on localhost:%d\n", cfg.Port)
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", cfg.Port), logs))
	}()
	time.Sleep(1 * time.Second)

	testRequest()

	select {
	case <-ctx.Done():
		//server.Stop()
		// Message to logger
	}
}

func testRequest() {
	jsonBody := []byte(`{"title": "World Cup", "address": "St. Petersburg", "season": "2022/2023"}`)
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", "http://localhost:8082/cup/register", bodyReader)
	if err != nil {
		log.Fatalf("Could not create request: %v\n", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Could not send request: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Expected status OK, got %v\n", resp.Status)
	}

	var responseBody bytes.Buffer
	responseBody.ReadFrom(resp.Body)
	fmt.Println("Response:", responseBody.String())
}
