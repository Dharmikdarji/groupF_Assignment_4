package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Weather struct for JSON response
type Weather struct {
	City        string `json:"city"`
	Temperature string `json:"temperature"`
	Weather     string `json:"weather"`
}

// Handler for GET /city
func getCityHandler(w http.ResponseWriter, r *http.Request) {
	cityName := r.URL.Query().Get("name")
	if cityName == "" {
		http.Error(w, "City name is required", http.StatusBadRequest)
		return
	}

	// Call external weather API (OpenWeatherMap)
	apiKey := "2b4e984d77b5a3d7d3e78833f3c398fa"
	apiUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", cityName, apiKey)
	resp, err := http.Get(apiUrl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch weather data, error: %s", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to fetch weather data, status code: %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	var weatherData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode weather data, error: %s", err), http.StatusInternalServerError)
		return
	}

	// Extract relevant information
	weather := Weather{
		City:        cityName,
		Temperature: fmt.Sprintf("%.1f°C", weatherData["main"].(map[string]interface{})["temp"].(float64)),
		Weather:     weatherData["weather"].([]interface{})[0].(map[string]interface{})["main"].(string),
	}

	// Convert weather struct to JSON and send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}

// Handler for POST /city
func postCityHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON input, error: %s", err), http.StatusBadRequest)
		return
	}

	// Call getCityHandler with the provided city name
	r.URL.Query().Set("name", input.Name)
	getCityHandler(w, r)
}

func main() {
	http.HandleFunc("/city", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCityHandler(w, r)
		case http.MethodPost:
			postCityHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
