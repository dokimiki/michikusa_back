FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go getNeaestStation.go getStationList.go

FROM gcr.io/distroless/base-debian10
COPY --from=builder /app/main /
CMD ["/main"]
