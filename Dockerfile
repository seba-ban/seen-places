FROM golang:1.21 AS init

RUN go install github.com/minio/mc@latest
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM golang:1.21 AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev

ENV LOCAL_STORAGE_FILES_DIR /data/storage
RUN go install github.com/cosmtrek/air@latest
ENTRYPOINT [ "air", "-c", ".air.toml" ]

FROM base AS build

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /seen

FROM scratch AS production

ENV GIN_MODE=release

WORKDIR /

COPY --from=build /seen /seen
COPY templates /templates

VOLUME [ "/data" ]
ENV LOCAL_STORAGE_DIR /data/storage
ENV TMP_DIR /data/tmp

EXPOSE 8080

ENTRYPOINT ["/seen"]