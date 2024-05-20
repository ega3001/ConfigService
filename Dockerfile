FROM golang:1.22.0 as build

WORKDIR /build

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN go build -o app

FROM golang:1.22.0

WORKDIR /app

COPY --from=build /build/app /app

ENTRYPOINT [ "/app/app" ]