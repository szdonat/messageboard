FROM golang:1.26.5-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /message-board ./cmd

FROM scratch

COPY --from=builder /message-board /message-board

EXPOSE 8080

ENTRYPOINT ["/message-board"]
