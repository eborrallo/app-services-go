# Services in Go

<p align="center">
<img src="https://external-preview.redd.it/creating-a-more-sustainable-model-for-oapi-codegen-in-the-v0-grxYBfrB_TY75WBP_OwdunQQkWeXgRpOCdwti_qRaGA.jpg?auto=webp&s=ee0a2ee35e1b55a3f70f7d0d3a57b07fd527720b" align="center"
alt="golang-logo"></p>

> ⚠️ IMPORTANT NOTE:
>The main goal of this project is how to implement different services in golang in a single repo (Monolit)
but with a clear structure and a good practices to split it in Microservices if was necessary . Some of the services are not already finalized, because in a real project the implementation part
will change to much and the idea is to do only the common part here but the main idea is to show how to implement it.

The architecture chose and implemented here was DDD  (Domain Driven Design) with a CQRS pattern (Command Query Responsibility Segregation) and Event Sourcing.
If you know DDD the scaffolding of the project will be very familiar to you. If not I will recommend you to read about it
before to continue with this project.

Quick list of the stack used in this project:

Golang, RabbitMQ , Redis , MySql
***

## What services we have here?

- **Auth Service**: This service is responsible to manage the users and the authentication process

## Auth Service
This main porpouse of this service is had all the users and the clasic authentication proces login , forgot , email validations ... , also has a DomainEvent to notify the other services when a user is created or token updated or user validated .All the persistence is done with a Mysql  database.
We have the following endpoints using REST API:
- [POST] /api/auth/user : This endpoint is used to create a new user. It expects a JSON payload with the user details. Upon successful creation, it returns a response with the created user's information.
- [GET] /api/auth/validate/:token : This endpoint is used to validate a user's email address. It expects a token as a URL parameter. Upon successful validation, it returns a response with the user's information.
- [POST] /api/auth/login : This endpoint is used to log in a user. It expects a JSON payload with the user's email and password. Upon successful login, it returns a response with the user's information and a JWT token.

Example request payload:

Run:
```bash
air
```
Grpc services:
 TODO

## Todo

- [x] Auth Service
    - [x] REST API
    - [x] Tests
    - [x] Docker
    - [x] JWT
    - [x] Domain Events
    - [x] Email Validation 
    - [ ] Forgot Password
    - [ ] Refresh Token

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
npm run test
```
