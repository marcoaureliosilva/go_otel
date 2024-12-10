package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/exporters/trace/zipkin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Definindo estrutura para requisição de CEP
type CepRequest struct {
	Cep string `json:"cep"`
}

// Inicializando Telemetria com Zipkin
func initTelemetry() {
	zipkinURL := "http://localhost:9411/api/v2/spans"

	exporter, err := zipkin.New(zipkin.WithCollectorEndpoint(zipkin.WithEndpoint(zipkinURL)))
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(trace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)
}

// Handler para processar requisições de CEP
func cepHandler(w http.ResponseWriter, r *http.Request) {
	var req CepRequest
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Decodificando o JSON recebido
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || len(req.Cep) != 8 || !isValidCep(req.Cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	// Chama o Serviço B
	resp, err := http.Get(fmt.Sprintf("http://localhost:8081/temperatura/%s", req.Cep))
	if err != nil {
		http.Error(w, "error contacting Service B", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Retorna a resposta do Serviço B
	w.WriteHeader(resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
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
	fmt.Println("Serviço A rodando na porta 8080")
	http.ListenAndServe(":8080", nil)
}
