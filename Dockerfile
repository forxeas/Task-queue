FROM golang:1.24-bookworm

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o /app/cmd/bin ./cmd/api

CMD ["/app/cmd/bin"]