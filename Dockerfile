FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/battery-api ./cmd/server
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/battery-api /app/battery-api
COPY --from=build /src/migrations /app/migrations
EXPOSE 8080
ENTRYPOINT ["/app/battery-api"]
