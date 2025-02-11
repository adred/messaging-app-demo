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

These features outline how I would finish the project:

- **File Sharing:**  
  Add support for file sharing in chats, with validation for PDF, JPEG, and PNG files.

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

### Running tests

```bash
go tests ./...
```

### Access the Service

- Messaging API: http://localhost:3000
- RabbitMQ Management Dashboard: http://localhost:15672
  (Default credentials: guest / guest)

## Technical Choices & System Architecture

### Language & Frameworks

- Go (Golang):
  Chosen for its performance, simplicity, and ease of deployment.

- Chi:
  A lightweight HTTP router that provides a simple and idiomatic way to build REST APIs.

- RabbitMQ:
  Used for asynchronous messaging to decouple message creation from downstream processing.

- Google Wire:
  Provides dependency injection, enabling a clean separation of concerns and easier testing.

- In-Memory Repositories:
  Used in this version for simplicity. In production, these would be replaced by persistent storage (e.g., MySQL or PostgreSQL).

### Architecture Overview

- Domain-Driven Design (DDD):
  The project is organized into multiple layers:

  - Domain: Contains core business entities (User, Chat, Message) and related logic.
  - Application: Contains the business logic (e.g., sending messages, creating chats, updating statuses).
  - Infrastructure: Provides integrations with external systems (API, repositories, RabbitMQ).
  - Configuration: Manages environment configuration.

- RESTful API:
  The API exposes endpoints for creating chats, sending messages, updating message statuses, retrieving chat messages, and listing user chats.

- Asynchronous Messaging:
  RabbitMQ is used to publish events asynchronously (e.g., when a message is sent), enabling future decoupled processing such as notifications or logging.

- Middleware:
  Basic authentication and rate limiting are applied via middleware to secure and protect API endpoints.

- Containerization:
  Docker and Docker Compose ensure a consistent deployment environment across development and production.

## Potential Improvements

- **Structured Logging:**  
  Integrate a logging library (e.g., [zap](https://github.com/uber-go/zap)) for structured and leveled logs.
- **Tracing:**  
  Implement distributed tracing using OpenTelemetry to track requests across microservices or internal components.
- **Hot Reload:**  
  Use [air](https://github.com/cosmtrek/air) for hot reloading during development to speed up the feedback loop.
- **Task Automation:**  
  Utilize a `Makefile` or a tool like [Task](https://github.com/go-task/task) to automate common tasks (e.g., running Wire for dependency injection code generation, starting Air, running tests, etc.).
- **Health Check Endpoint:**  
  Expose a health check endpoint (e.g., `/health`) to monitor the status of the application and its dependencies.

## API Documentation

Swagger UI / OpenAPI Spec:

- OpenAPI YAML: http://localhost:3000/docs/openapi.yaml
- Swagger UI: http://localhost:3000/openapi/index.html
