FROM golang:1.22.1-alpine AS stage1
WORKDIR /app
COPY . .
RUN go mod tidy
WORKDIR /app/cmd/server
RUN go build -o /app/main .

FROM alpine:latest
WORKDIR /app
COPY --from=stage1 /app/main /app/main
EXPOSE 80
CMD ["./main"]

