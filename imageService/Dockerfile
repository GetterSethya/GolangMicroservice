#Build
FROM golang:latest AS build

WORKDIR /app

COPY library/ /app/library/
COPY userProto/ /app/userProto/
COPY imageProto/ /app/imageProto/
COPY imageService/ .
ENV GOPROXY https://proxy.golang.org

RUN go mod edit -replace=github.com/GetterSethya/library=/app/library
RUN go mod edit -replace=github.com/GetterSethya/userProto=/app/userProto
RUN go mod edit -replace=github.com/GetterSethya/imageProto=/app/imageProto
RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

#Final
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/myapp .

RUN chmod +x myapp

EXPOSE 3001

CMD ["./myapp"]

