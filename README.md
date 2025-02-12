# Flight API Service

This is a Go-based REST API service that provides flight-related information through various endpoints.

## API Endpoints

- `GET /v1/api/flights/all_params/search` - Search flights using multiple parameters
- `GET /v1/api/flights/:id` - Get detailed information about a specific flight based on `flight_id`
- `GET /v1/api/flights/dest_time/search` - Search flights by destination and time

## Installation

1. Make sure you have Go installed (version 1.16 or higher)
2. Clone this repository
```bash
git clone https://github.com/uzzalcse/flight-details.git
cd flight-details
```
3. Install dependencies
```bash
go mod tidy
```

## Running the Application

1. Start the server
```bash
go run main.go
```
The server will start on `localhost:8080` by default.

## API Documentation

This project includes Swagger UI for API documentation and testing.

To access the Swagger UI:
1. Start the server
2. Visit `http://localhost:8080/swagger/` in your browser

## Development

To update the Swagger documentation, run:
```bash
swag init
```

## Technologies Used

- Go
- Beego Framework
- Swagger UI
- Elasticsearch
