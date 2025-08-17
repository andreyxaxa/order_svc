# Step 1: Modules caching
FROM golang:1.24-alpine AS modules

WORKDIR /modules

COPY go.mod go.sum /modules/
RUN go mod download

# Step 2: Builder
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY --from=modules /go/pkg /go/pkg
COPY . /app/
COPY .env /

RUN CGO_ENABLED=0 go build -o /bin/app ./cmd/app

# Step 3: Final
FROM scratch

COPY --from=builder /app/config /config
COPY --from=builder /bin/app /app
COPY --from=builder /.env /
COPY --from=builder /app/docs/html /docs/html

CMD [ "/app" ]