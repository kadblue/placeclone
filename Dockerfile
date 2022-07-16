FROM golang:alpine AS builder

WORKDIR /go/src/app

COPY / .

RUN go build -o server

FROM alpine
WORKDIR /app
COPY --from=builder /go/src/app/ /app/

CMD ["./server"]
