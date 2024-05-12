# Go Services Implementation

<p align="center">
    <img src="https://external-preview.redd.it/creating-a-more-sustainable-model-for-oapi-codegen-in-the-v0-grxYBfrB_TY75WBP_OwdunQQkWeXgRpOCdwti_qRaGA.jpg?auto=webp&s=ee0a2ee35e1b55a3f70f7d0d3a57b07fd527720b" align="center" alt="golang-logo">
</p>

> ⚠️ **IMPORTANT NOTE:**
> The primary aim of this project is to demonstrate the implementation of various services in Go within a single repository (Monolith), while adhering to clear structure and good practices that facilitate splitting into Microservices if necessary. Some services are not finalized as the implementation details may vary significantly in real-world scenarios. However, the core objective is to illustrate common implementation approaches.

The architecture adopted and implemented here follows Domain Driven Design (DDD) with a CQRS pattern (Command Query Responsibility Segregation) and Event Sourcing. Familiarity with DDD will ease understanding the project structure. If unfamiliar, it's recommended to explore DDD before proceeding.

Quick overview of the technology stack used:

- Golang
- RabbitMQ
- Redis
- MySQL

***

## Services Overview

- **Auth Service**: Responsible for user management and authentication processes.

## Auth Service

The Auth Service primarily manages user-related functionalities such as user creation, authentication (login), password recovery, and email validation. It incorporates Domain Events to notify other services about user-related events. Persistence is handled using a MySQL database. The service exposes the following REST API endpoints:

- [POST] `/api/auth/user`: Creates a new user. Expects a JSON payload with user details. Upon success, returns user information.
- [GET] `/api/auth/validate/:token`: Validates a user's email address using a token provided as a URL parameter. Returns user information upon successful validation.
- [POST] `/api/auth/login`: Logs in a user. Expects a JSON payload with email and password. Returns user information and a JWT token upon successful login.

Run:
```bash
air
```
Grpc services:
 TODO

## Todo

- [x] WEB-2 Auth Service
    - [x] REST API
    - [x] Tests
    - [x] JWT
    - [x] Domain Events
    - [x] Email Validation 
    - [x] Forgot Password
    - [x] Refresh Token
    - [ ] Middleware
- [x] WEB-3 Auth Service
    - [x] REST API
    - [x] Tests
    - [x] JWT
    - [x] Domain Events
    - [x] Integrate go-ethereum
    - [x] Validate signature
    - [x] Integrate Reddis
    - [ ] Middleware
- [x] Integrate Domain events RabbitMq bus
- [x] Docker
- [ ] UI
    - [ ] Web3 Auth Service
    - [ ] Web2 Auth Service
***

## Docker
All the project is dockerized , in the Dockerfile you can see the configuration of each service and also in the docker-compose.yml file

you can run the following command to start all the services:
```bash
docker-compose up
```

## Tests
All the services have tests implemented with Jest and dockerTest, you can run the following command to run the tests:
```bash
go test ./...
```
