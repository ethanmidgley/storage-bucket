FROM golang:1.17.7-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o storage-bucket ./cmd/main.go

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/storage-bucket ./storage-bucket
COPY bucket.yaml ./

EXPOSE 5000

USER nonroot:nonroot

ENTRYPOINT ["/storage-bucket"]