FROM golang:1.23-alpine
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o wallet-topup .
EXPOSE 8080
CMD ["./wallet-topup"]