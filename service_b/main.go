package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Estrutura para a resposta de temperatura
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
			attribute.String("service.name", "service_b"),
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

// Handler para processar requisições de CEP e buscar temperatura
func temperaturaHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Path[len("/temperatura/"):]

	// Inicia uma nova span
	ctx, span := otel.Tracer("service_b").Start(r.Context(), "temperaturaHandler")
	defer span.End()
	fmt.Println(ctx)

	// Consultar a API ViaCEP para obter a cidade
	location, err := getLocationByCEP(cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Consultar a API WeatherAPI para obter a temperatura
	temperature, err := getTemperature(location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := TemperatureResponse{
		City:  location,
		TempC: temperature.TempC,
		TempF: temperature.TempF,
		TempK: temperature.TempK,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Função que consulta a API ViaCEP
func getLocationByCEP(cep string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return "", fmt.Errorf("invalid zipcode")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid zipcode")
	}

	var data struct {
		Localidade string `json:"localidade"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode location data")
	}
	return data.Localidade, nil
}

// Função que faz a consulta de temperaturas
func getTemperature(city string) (*TemperatureResponse, error) {
	apiKey := "SUA_API_KEY_AQUI" // Substitua por sua chave da WeatherAPI
	resp, err := http.Get(fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, city))
	if err != nil {
		return nil, fmt.Errorf("error fetching temperature")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cannot find zipcode")
	}

	var data struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}

	// Ler a resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	// Calcula as temperaturas em diferentes escalas
	tempC := data.Current.TempC
	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15

	return &TemperatureResponse{
		City:  city,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}

// Função main
func main() {
	initTelemetry()
	http.HandleFunc("/temperatura/", temperaturaHandler)
	fmt.Println("Serviço B rodando na porta 8080")
	http.ListenAndServe(":8080", nil)
}
