#Build
FROM golang:latest AS build

WORKDIR /app

COPY library/ /app/library/
COPY imageService/ .

RUN go mod edit -replace=github.com/GetterSethya/library=/app/library
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

#Final
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/myapp .

RUN chmod +x myapp

EXPOSE 3001

CMD ["./myapp"]

