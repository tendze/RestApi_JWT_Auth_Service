FROM golang:1.21

RUN go version
ENV GOPATH=/

COPY ./ ./
RUN go mod download
RUN go build -o jwt-auth-service ./cmd/jwt-auth-service/main.go

cmd ["./jwt-auth-service"]