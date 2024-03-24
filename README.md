# Weather API Server

This project is a RESTful API server developed in Go that provides information about the current weather in a specified city. It integrates with an external weather API (OpenWeatherMap) to fetch real-time weather data.

## Features

- Handles both GET and POST requests at the endpoint `/city`.
- Accepts a city name as input and returns the current weather in JSON format.
- Containerized using Docker for easy deployment.

## API Endpoints

### GET /city

Accepts a city name as a query parameter and returns the current weather in JSON format.

#### Example request:

GET /city?name=Toronto

#### Example response:
json
"{
  "city": "Toronto",
  "temperature": "5°C",
  "weather": "Cloudy"
}"

### POST /city

Accepts a JSON body with a city name and returns the current weather in JSON format.

#### Example request:

GET /city?name=Toronto

#### Example response:
```json
{
  "city": "Toronto",
  "temperature": "5°C",
  "weather": "Cloudy"
}


