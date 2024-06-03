FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/planny/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY configs/.env ./configs/.env
COPY scripts/start.sh .
COPY scripts/wait-for.sh .

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]