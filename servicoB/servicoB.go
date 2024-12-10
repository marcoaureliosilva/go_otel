package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LocationResponse struct {
	Localidade string `json:"localidade"`
}

type TemperatureResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func temperaturaHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Path[len("/temperatura/"):]

	city, err := getCityByCep(cep)
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	temperature, err := getTemperatureByCity(city)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	response := TemperatureResponse{
		City:  city,
		TempC: temperature,
		TempF: convertCtoF(temperature),
		TempK: convertCtoK(temperature),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getCityByCep(cep string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("CEP inválido")
	}

	var loc LocationResponse
	err = json.NewDecoder(resp.Body).Decode(&loc)
	if err != nil {
		return "", err
	}

	return loc.Localidade, nil
}

func getTemperatureByCity(city string) (float64, error) {
	// Aqui você deve integrar com a API de clima
	// Para simplicidade, retornando uma temperatura fixa.
	return 28.5, nil // Esta deveria ser a chamada para a API real
}

func convertCtoF(c float64) float64 {
	return c*1.8 + 32
}

func convertCtoK(c float64) float64 {
	return c + 273.15
}

func main() {
	http.HandleFunc("/temperatura/", temperaturaHandler)
	fmt.Println("Serviço B rodando na porta 8081")
	http.ListenAndServe(":8081", nil)
}
