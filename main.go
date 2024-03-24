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
	apiKey := "WeatherCheckAssignment"
	apiUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", cityName, apiKey)
	resp, err := http.Get(apiUrl)
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}

	var weatherData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		http.Error(w, "Failed to decode weather data", http.StatusInternalServerError)
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
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
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
