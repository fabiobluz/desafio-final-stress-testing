# build stage
FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /gload ./cmd/gload

# runtime
FROM gcr.io/distroless/base-debian12
COPY --from=build /gload /gload
ENTRYPOINT ["/gload"]
