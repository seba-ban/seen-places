FROM golang:1.21 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

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