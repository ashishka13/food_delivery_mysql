FROM golang:1.18.2-alpine AS builder
RUN apk add git

WORKDIR /build
COPY . .

ENV CGO_ENABLED=0 
RUN go build -o main .

FROM alpine:latest
# COPY --from=builder /build/config.json config.json
COPY --from=builder /build/main main

EXPOSE 8080
ENTRYPOINT ["./main"]