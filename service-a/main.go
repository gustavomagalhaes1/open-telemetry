package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Cep struct {
	Cep string `json:"cep"`
}

func initTrace() {
	// Create the OTLP trace exporter
	ctx := context.Background()
	client := otlptracehttp.NewClient(otlptracehttp.WithEndpoint("otel-collector:4317"), otlptracehttp.WithInsecure())
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	// Create the Zipkin trace exporter
	zipkinExporter, err := zipkin.New("http://zipkin:9411/api/v2/spans")
	if err != nil {
		log.Fatalf("failed to create zipkin exporter: %v", err)
	}

	// Create the trace provider with the exporters
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithBatcher(zipkinExporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("service-a"),
		)),
	)
	otel.SetTracerProvider(tp)
}


func main() {
	initTrace()
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/", handleRequest)
	
	log.Println("Iniciando o servidor na porta 8080")
	http.ListenAndServe(":8080", router)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("service-a").Start(r.Context(), "GetCepWheather")

	// Parse the request
	var requestData Cep
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	// Check if the zipcode is valid
	if !isValidCEP(requestData.Cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	resp, err := sendToServiceB	(ctx, requestData.Cep)
	if err != nil {
		http.Error(w, "Failed to contact service b" + err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, "Failed to copy response body", http.StatusInternalServerError)
		return
	}
	span.End()
}

func isValidCEP(cep string) bool {
	if cep == "" || len(cep) == 0 {
		return false
	}

	if len(cep) != 8 {
		return false
	}

	for _, d := range cep {
		if d < '0' || d > '9' {
			return false
		}
	}

	return true
}

func sendToServiceB(ctx context.Context, cep string) (*http.Response, error) {
	// Create a new span to send the request to service b
	ctx, span := otel.Tracer("service-a").Start(ctx, "sendToServiceB")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://goapp-service-b:8081/"+cep, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}
