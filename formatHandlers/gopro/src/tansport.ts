import * as fs from "node:fs/promises";
import * as path from "node:path";

import { Channel, ConsumeMessage, connect } from "amqplib";
import { BufferWriter } from "protobufjs";

import { processRequest } from "./dataHandling";
import { ExtractedFilePoints, ProcessFileRequest } from "./proto/events";
import {
  amqpSettings,
  channelSettings,
  localStorageSettings,
} from "./settings";

export const startTransport = async () => {
  const queue = channelSettings.workQueue;
  const outQueue = channelSettings.targetQueue;

  const conn = await connect({
    protocol: amqpSettings.protocol,
    hostname: amqpSettings.hostname,
    port: amqpSettings.port,
    username: amqpSettings.username,
    password: amqpSettings.password,
    vhost: amqpSettings.vhost,
  });

  const inputChannel = await conn.createChannel();
  await inputChannel.assertQueue(queue, { durable: true });

  const outputChannel = await conn.createChannel();
  await outputChannel.assertQueue(outQueue, { durable: true });

  return { conn, inputChannel, outputChannel };
};

export const messageHandlerFactory = (
  inputChannel: Channel,
  outputChannel: Channel,
) => {
  return async (msg: ConsumeMessage | null) => {
    if (msg === null) {
      console.log("Consumer cancelled by server");
      return;
    }
    console.log("Received:", msg.content.toString());

    const req = ProcessFileRequest.decode(msg.content);

    if (!req.filepath) {
      console.log("invalid request");

      await inputChannel.ack(msg);
      return;
    }

    if (
      !(await fs.stat(path.join(localStorageSettings.filesDir, req.filepath)))
    ) {
      console.log("file not found");

      await inputChannel.ack(msg);
      return;
    }

    console.log("processing points");
    for await (const points of processRequest(req)) {
      console.log("sending points");
      // @ts-expect-error: I think typing is always assuming particular writer
      const buff: Buffer = ExtractedFilePoints.encode(
        points,
        new BufferWriter(),
      ).finish();
      await outputChannel.sendToQueue(channelSettings.targetQueue, buff);
    }

    await inputChannel.ack(msg);
  };
};

export async function startTransportHandling() {
  const { conn, inputChannel, outputChannel } = await startTransport();
  inputChannel.consume(
    channelSettings.workQueue,
    messageHandlerFactory(inputChannel, outputChannel),
  );

  console.log("transport started");

  process.on("SIGINT", async () => {
    console.log("SIGINT received");
    await conn.close();
    process.exit(0);
  });
}
