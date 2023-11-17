#!/bin/bash

set -e

# TODO: for now used only in local development inside docker-compose
mc alias set local "${S3_ENDPOINT_URL}" "${S3_ACCESS_KEY}" "${S3_SECRET_KEY}"
mc mb --ignore-existing "local/${S3_RAW_FILES_BUCKET}"
   