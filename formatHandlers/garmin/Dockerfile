FROM python:3.11 AS dev

WORKDIR /app
COPY README.md pyproject.toml poetry.lock ./
RUN pip install poetry \
    && poetry config virtualenvs.in-project false \
    && poetry install
COPY --from=cosmtrek/air /go/bin/air /usr/local/bin/air

CMD [ "air", "-c", ".air.toml" ]

FROM python:3.11 AS builder

WORKDIR /build

COPY ./garmin_format_handler ./garmin_format_handler
COPY README.md pyproject.toml poetry.lock ./
RUN pip install poetry pex \
    && poetry config virtualenvs.create false \
    && poetry install --without dev --no-root \
    && poetry build \
    && pex $(pip freeze) dist/*.whl -o app.pex

FROM python:3.11-slim AS production

VOLUME [ "/data" ]
ENV LOCAL_STORAGE_FILES_DIR /data/storage

COPY --from=builder /build/app.pex /app.pex

CMD ["/app.pex", "-m", "garmin_format_handler"]
