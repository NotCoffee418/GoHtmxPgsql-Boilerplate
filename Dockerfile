# Set up Go
FROM golang:latest AS builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o main .

# Set up Node
FROM node:latest AS node-builder
WORKDIR /app
COPY package*.json ./
RUN npm ci

# Prepare final image
FROM alpine:latest
WORKDIR /app
ENV GIN_MODE=release
COPY --from=builder /app /app
RUN chmod +x /app/main
COPY --from=node-builder /app/node_modules /app/node_modules

# Expose port and run Go binary
EXPOSE 8080
CMD ["./main"]
