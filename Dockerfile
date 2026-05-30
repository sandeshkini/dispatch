FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o dispatch .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /src/dispatch .
EXPOSE 8888
ENTRYPOINT ["./dispatch"]
