FROM golang:1.23-alpine AS construction
WORKDIR /application
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /serveur ./cmd/serveur

FROM construction AS verification
RUN go vet ./...

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
COPY --from=construction /serveur /serveur
EXPOSE 8080
ENTRYPOINT ["/serveur"]
