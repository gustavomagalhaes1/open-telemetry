# Builder
FROM golang:1.22.2 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o goapp-b ./cmd/server

FROM scratch
EXPOSE 8081
WORKDIR /app
COPY --from=build /app/goapp-b .
ENTRYPOINT ["./goapp-b"]
