FROM golang:alpine AS build

RUN apk add --update git nodejs npm
WORKDIR /go/src/api
COPY . .

RUN npm install
RUN npm run build
RUN CGO_ENABLED=0  go run github.com/a-h/templ/cmd/templ@latest generate
RUN CGO_ENABLED=0 go build -o /go/bin/app-services-go cmd/api/main.go

# Building image with the binary
#FROM scratch
#COPY --from=build /go/bin/app-services-go /go/bin/app-services-go
#ENTRYPOINT ["/go/bin/app-services-go"]
