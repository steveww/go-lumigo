FROM golang:1.19 AS build

WORKDIR /app

COPY go.* ./
COPY *.go ./
RUN go mod tidy
RUN go build -o go-lumigo .

FROM debian:bullseye-slim
WORKDIR /app
EXPOSE 8080
COPY --from=build /app/go-lumigo ./

CMD ["./go-lumigo"]
