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
4. Set up configuration.
- Create `app.conf` file inside ***conf*** directory.
```bash
touch conf/app.conf
```
- Copy the variables from **app.conf.sample** and enter your configurations. Example is provided below.
```bash
appname = flight-api
httpport = 8080
runmode = dev
autorender = false
copyrequestbody = true
EnableDocs = true
ES_LOCAL_API_KEY=QnZnMjdwUUJfZXVoNWRBbE1MaTg6c19PM0hWUVFRay1QM0QyLXNuWE1fZw==
ES_LOCAL_URL=http://elasticsearch:9200
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
