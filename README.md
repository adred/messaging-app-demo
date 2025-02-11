# Messaging Service

A simple messaging service built in Go using Domain-Driven Design principles, RESTful APIs, and RabbitMQ for asynchronous messaging. This project was designed as an assessment and includes core features such as sending messages, retrieving message history, tracking message status, and managing private chats between hardcoded users.

## Table of Contents

- [Features](#features)
- [Setup Instructions](#setup-instructions)
- [Technical Choices & System Architecture](#technical-choices--system-architecture)
- [Potential Improvements](#potential-improvements)
- [API Documentation](#api-documentation)

## Features

### Core Features

- **REST API Endpoints:**
  - Send a message to an existing chat.
  - Retrieve message history for a chat.
  - Update message status (e.g., sent, delivered, read, failed).
  - List all chats a user participates in.
  - Create a chat by providing two user IDs.
- **Hardcoded Users:**  
  The system uses a hardcoded list of four users (Red, Jrue, Miro, Joann) as valid recipients.

### Additional (Planned) Features

These features outline how I would finish and enhance the project:

- **Structured Logging:**  
  Integrate a logging library (e.g., [zap](https://github.com/uber-go/zap)) for structured and leveled logs.
- **Tracing:**  
  Implement distributed tracing using OpenTelemetry to track requests across microservices or internal components.
- **File Sharing:**  
  Add support for file sharing in chats, with validation for PDF, JPEG, and PNG files.
- **Hot Reload:**  
  Use [air](https://github.com/cosmtrek/air) for hot reloading during development to speed up the feedback loop.
- **Task Automation:**  
  Utilize a `Makefile` or a tool like [Task](https://github.com/go-task/task) to automate common tasks (e.g., running Wire for dependency injection code generation, starting Air, running tests, etc.).
- **Health Check Endpoint:**  
  Expose a health check endpoint (e.g., `/health`) to monitor the status of the application and its dependencies.

## Setup Instructions

### Prerequisites

    - [Go 1.20+](https://golang.org/dl/)
    - [Docker & Docker Compose](https://www.docker.com/)

### Running with Docker Compose

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/adred/messaging-app-demo.git
   cd messaging-service

   ```

2. **Configure Environment Variables:**
   Create a .env file in the project root with variables such as:

   ```bash
   RABBITMQ_DSN=amqp://guest:guest@rabbitmq:5672/
   RABBITMQ_QUEUE=messages
   HTTP_PORT=3000
   AUTH_USERNAME=red
   AUTH_PASSWORD=abc123
   RATE_LIMIT=100
   ```

3. **Build and Run Containers:**

   ```bash
   docker-compose up --build
   ```

4. **Access the Service:**

- Messaging API: http://localhost:3000
- RabbitMQ Management Dashboard: http://localhost:15672
  (Default credentials: guest / guest)

## API Documentation

Swagger UI / OpenAPI Spec:

- OpenAPI YAML: http://localhost:3000/docs/openapi.yaml
- Swagger UI: http://localhost:3000/openapi/index.html
