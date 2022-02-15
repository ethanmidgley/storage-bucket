FROM golang:1.17.7-buster AS build

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o storage-bucket ./cmd/main.go

FROM ubuntu

WORKDIR /

RUN mkdir -p /data

COPY --from=build /build/storage-bucket .

EXPOSE 5000

VOLUME ["/data"]

ENTRYPOINT ["/storage-bucket", "-pathPrefix", "/data"]