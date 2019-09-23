# build-stage
FROM golang:1.13 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/

# container-stage
FROM golang:1.13
WORKDIR /app

COPY --from=build /app/main .
COPY --from=build /app/configuration/config.yaml .

EXPOSE 5678
CMD ["./main", "-config", "."]