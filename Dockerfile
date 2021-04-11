
FROM golang:1.16 AS build
WORKDIR /src
COPY . .
RUN rm .env
RUN go get -u github.com/swaggo/swag/cmd/swag

RUN swag init -g main.go
RUN go build
EXPOSE 8100
RUN chmod +x ./LocalGuideUserService
CMD ["./LocalGuideUserService"]