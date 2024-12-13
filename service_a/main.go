package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Definindo estrutura para requisição de CEP
type CepRequest struct {
	Cep string `json:"cep"`
}

// Estrutura para a resposta do Serviço B
type TemperatureResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// Inicializando Telemetria com Zipkin
func initTelemetry() {
	reporter, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		panic(err)
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "service_a"),
		),
	)
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(reporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
}

// Handler para processar requisições de CEP
func cepHandler(w http.ResponseWriter, r *http.Request) {
	var req CepRequest
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || len(req.Cep) != 8 || !isValidCep(req.Cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	// Chamar Serviço B e medir o tempo com uma span
	ctx, span := otel.Tracer("service_a").Start(r.Context(), "cepHandler")
	defer span.End()
	fmt.Println(ctx)

	resp, err := http.Get(fmt.Sprintf("http://localhost:8081/temperatura/%s", req.Cep)) // Chama o Serviço B
	if err != nil {
		http.Error(w, "error contacting Service B", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Lê a resposta do Serviço B
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "error reading response", http.StatusInternalServerError)
		return
	}

	// Incluindo o contexto na nova resposta
	span.SetAttributes(attribute.String("response.status", resp.Status))
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

// Função para validar o CEP
func isValidCep(cep string) bool {
	return len(cep) == 8 && strings.TrimSpace(cep) != ""
}

// Função main
func main() {
	initTelemetry()
	http.HandleFunc("/cep", cepHandler)
	fmt.Println("Serviço A rodando na porta 8082")
	http.ListenAndServe(":8082", nil)
}
