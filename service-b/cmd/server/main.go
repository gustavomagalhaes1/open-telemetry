package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/magalhaesgustavo/cloud-run/pkg/cep"
	"github.com/magalhaesgustavo/cloud-run/pkg/weather"
)

type TemperatureResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
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
			semconv.ServiceNameKey.String("service-b"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func main() {
	initTrace()

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	
	router.Route("/{cep}", func(r chi.Router) {
		
		r.Use(checkCepMiddleware)
		r.Get("/", handleGetTemperatureByCEP)
	})
	
	log.Println("Iniciando o servidor na porta 8081")
	http.ListenAndServe(":8081", router)
}

func checkCepMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cep := chi.URLParam(r, "cep")

		if cep == "" || len(cep) == 0 {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		if len(cep) != 8 {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		for _, d := range cep {
			if d < '0' || d > '9' {
				http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			}
		}

		next.ServeHTTP(w, r)
	})
}

func handleGetTemperatureByCEP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	propagator := propagation.TraceContext{}
	ctx1 := propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

	// Create a new span to validate the CEP
	ctx1, span := otel.Tracer("service-b").Start(ctx1, "GetCEP")

	cepReq := chi.URLParam(r, "cep")
	address, err := cep.GetAddressFromViaCEP(cepReq)
	if err != nil || address == nil{
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	span.End()

	_, span2 := otel.Tracer("service-b").Start(ctx1, "GetWeather")
	weatherResponse, err := weather.GetWeather(address.Localidade)
	log.Println(weatherResponse.Current.TempC)
	if err != nil {
		http.Error(w, "can not find weather", http.StatusNotFound)
		return
	}

	temperature := TemperatureResponse{
		City:  address.Localidade,
		TempC: weatherResponse.Current.TempC,
		TempF: weather.CelsiusToFahrenheit(weatherResponse.Current.TempC),
		TempK: weather.CelsiusToKelvin(weatherResponse.Current.TempC),
	}
	json.NewEncoder(w).Encode(temperature)
	span2.End()
}
