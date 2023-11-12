import { statSync } from "node:fs";

import * as env from "env-var";

export const amqpSettings = {
  protocol: env.get("RABBITMQ_PROTOCOL").default("amqp").asString(),
  hostname: env.get("RABBITMQ_HOST").required().asString(),
  port: env.get("RABBITMQ_PORT").default(5672).asPortNumber(),
  username: env.get("RABBITMQ_USERNAME").asString(),
  password: env.get("RABBITMQ_PASSWORD").asString(),
  vhost: env.get("RABBITMQ_VHOST").default("/").asString(),
};

export const channelSettings = {
  workQueue: env.get("RABBITMQ_WORK_QUEUE").default("gopro").asString(),
  targetQueue: env
    .get("RABBITMQ_TARGET_QUEUE")
    .default("extractedpoints")
    .asString(),
  prefetch: env.get("RABBITMQ_PREFETCH").default(1).asIntPositive(),
  pointsBufferSize: env
    .get("RABBITMQ_POINTS_BUFFER_SIZE")
    .default(10000)
    .asIntPositive(),
};

export const localStorageSettings = {
  filesDir: env.get("LOCAL_STORAGE_FILES_DIR").required().asString(),
};

if (!statSync(localStorageSettings.filesDir).isDirectory()) {
  throw new Error(`Directory ${localStorageSettings.filesDir} does not exist`);
}
