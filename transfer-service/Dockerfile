FROM golang:1.18-alpine
RUN apk add build-base

EXPOSE ${TRANSFERSERVICE_API_BIND_ADDR}
WORKDIR /app
COPY ./transfer-service /app/transfer-service
WORKDIR /app/transfer-service/

RUN go mod download && go mod verify && go build -o ./service
CMD ["./service"]