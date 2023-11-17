PROTO_DIR := proto
GEN_PATH_TS := formatHandlers/gopro/src/proto
GEN_PATH_GO := protogo
GEN_PATH_PYTHON := formatHandlers/garmin/garmin_format_handler/proto

# TODO: add cleanup before regenerating
# TODO: make this more user friendly, for now it was just quickly thrown together

.PHONY: generate
generate: generate-go generate-python generate-ts

# # https://developers.google.com/protocol-buffers/docs/reference/go-generated
.PHONY: generate-go
generate-go:
	PATH="${PATH}:$(shell go env GOPATH)/bin" protoc \
		-I=$(PROTO_DIR) \
		--go_out=$(GEN_PATH_GO) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_PATH_GO) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

.PHONY: generate-ts
generate-ts:
	mkdir -p $(GEN_PATH_TS)
	protoc \
		-I=$(PROTO_DIR) \
		--plugin=./formatHandlers/gopro/node_modules/.bin/protoc-gen-ts_proto \
		--ts_proto_out=$(GEN_PATH_TS) \
		--ts_proto_opt=env=node \
		--ts_proto_opt=outputServices=grpc-js \
		--ts_proto_opt=useAsyncIterable=true \
		$(PROTO_DIR)/*.proto

.PHONY: generate-python
generate-python:
	protoc \
		-I=$(PROTO_DIR) \
		--pyi_out=$(GEN_PATH_PYTHON) \
		--python_out=$(GEN_PATH_PYTHON) \
		$(PROTO_DIR)/*.proto

	cd formatHandlers/garmin && poetry run python -m grpc_tools.protoc \
		-I=../../$(PROTO_DIR) \
		--grpc_python_out=../../$(GEN_PATH_PYTHON) \
		../../$(PROTO_DIR)/*.proto

	cd formatHandlers/garmin && poetry run protol \
		--create-package \
		--in-place \
		--python-out ../../$(GEN_PATH_PYTHON) \
		protoc --proto-path=../../$(PROTO_DIR) ../../$(PROTO_DIR)/*.proto

.PHONY: new-migration
new-migration:
	migrate create -ext sql -dir db/migrations $(name)

.PHONY: gen-sqlc
gen-sqlc:
	$(shell go env GOPATH)/bin/sqlc generate

.PHONY: migrate
migrate:
	migrate -path db/migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" up

.PHONY: down
down:
	migrate -path db/migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" down -all

.PHONY: init-tmp-dir
init-tmp-dir:
	mkdir -p .tmp/storage