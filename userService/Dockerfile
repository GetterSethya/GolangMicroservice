#Build
FROM golang:latest AS build

WORKDIR /app

COPY library/ /app/library/
COPY userService/ .

RUN go mod edit -replace=github.com/GetterSethya/library=/app/library
RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o myapp .

#Final
FROM debian:latest

WORKDIR /app

COPY --from=build /app/myapp .

RUN chmod +x myapp

EXPOSE 3002

CMD ["./myapp"]

