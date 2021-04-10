
FROM golang:1.16 AS build
WORKDIR /src
COPY . .
RUN rm .env
RUN go build
EXPOSE 8100
RUN chmod +x ./LocalGuideUserService
CMD ["./LocalGuideUserService"]