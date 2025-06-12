# Stage 1: Build the Go application
FROM golang:1.22.6-alpine3.20 AS build
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Stage 2: Create the final image
FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/main .
COPY app.env .

EXPOSE 8080
ENTRYPOINT [ "/app/main" ]