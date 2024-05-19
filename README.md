# Exchange rates api
This is a test case for second step of selection for Software Engineering School 4.0

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Installing

1. Clone the repository
```sh
git clone https://github.com/myroslavve/genesis-test-case.git
```
2. Navigate to the project directory
```sh
cd genesis-test-case
```

### Setting Up the Environment

1. Copy the `.env.example` file to a new file named `.env`:
```sh
cp .env.example .env
```
2. Update the `.env` file with your own values for SMTP server and database.

### Running the Application

To build and run the application using Docker Compose:

```sh
docker-compose up --build
```

This command will start the MongoDB container, build the Go application, run the migrations, and start the application.

The application will be available on `localhost:8080`

## Endpoints

### 1. Get Current Exchange Rate

- **Endpoint**: `/api/rate`
- **Method**: `GET`
- **Description**: Retrieves the current exchange rate from USD to UAH.
- **Response**:
  - `200 OK`: Returns the current exchange rate as a JSON number.
  - `400 Bad Request`: Invalid request.

### 2. Subscribe Email

- **Endpoint**: `/api/subscribe`
- **Method**: `POST`
- **Description**: Subscribes an email to receive current exchange rate updates.
- **Request**:
  - Content-Type: `application/x-www-form-urlencoded`
  - Parameters:
    - `email` (string, required): The email address to subscribe.
- **Response**:
  - `200 OK`: Email added successfully.
  - `409 Conflict`: Email already exists in the database.
