#Build
FROM alpine:edge AS build

RUN apk add --no-cache --update go gcc g++
WORKDIR /app

COPY library/ /app/library/
COPY userProto/ /app/userProto/
COPY imageProto/ /app/imageProto/
COPY userService/ .
ENV GOPROXY https://proxy.golang.org

RUN go mod edit -replace=github.com/GetterSethya/library=/app/library
RUN go mod edit -replace=github.com/GetterSethya/userProto=/app/userProto
RUN go mod edit -replace=github.com/GetterSethya/imageProto=/app/imageProto
RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o myapp .

#Final
FROM alpine:edge

WORKDIR /app

COPY --from=build /app/myapp .

RUN chmod +x myapp

EXPOSE 3002

CMD ["./myapp"]

