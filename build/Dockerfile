FROM golang:1.20.4-alpine3.17 as build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o simple-ws ./cmd/main.go


FROM alpine:3.17 as app

COPY --from=build /build/simple-ws /app/simple-ws

EXPOSE 8000

ENTRYPOINT [ "/app/simple-ws" ,"-port=8000"]