/* eslint-disable */
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "";

export interface ProcessFileRequest {
  filepath: string;
}

export interface ExtractedFilePoint {
  filepath: string;
  timestamp: string;
  latitude: number;
  longitude: number;
  /** json metadata */
  metadata: string;
}

export interface ExtractedFilePoints {
  points: ExtractedFilePoint[];
}

function createBaseProcessFileRequest(): ProcessFileRequest {
  return { filepath: "" };
}

export const ProcessFileRequest = {
  encode(message: ProcessFileRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.filepath !== "") {
      writer.uint32(10).string(message.filepath);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ProcessFileRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProcessFileRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.filepath = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<ProcessFileRequest, Uint8Array>
  async *encodeTransform(
    source:
      | AsyncIterable<ProcessFileRequest | ProcessFileRequest[]>
      | Iterable<ProcessFileRequest | ProcessFileRequest[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [ProcessFileRequest.encode(p).finish()];
        }
      } else {
        yield* [ProcessFileRequest.encode(pkt as any).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, ProcessFileRequest>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<ProcessFileRequest> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [ProcessFileRequest.decode(p)];
        }
      } else {
        yield* [ProcessFileRequest.decode(pkt as any)];
      }
    }
  },

  fromJSON(object: any): ProcessFileRequest {
    return { filepath: isSet(object.filepath) ? globalThis.String(object.filepath) : "" };
  },

  toJSON(message: ProcessFileRequest): unknown {
    const obj: any = {};
    if (message.filepath !== "") {
      obj.filepath = message.filepath;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ProcessFileRequest>, I>>(base?: I): ProcessFileRequest {
    return ProcessFileRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ProcessFileRequest>, I>>(object: I): ProcessFileRequest {
    const message = createBaseProcessFileRequest();
    message.filepath = object.filepath ?? "";
    return message;
  },
};

function createBaseExtractedFilePoint(): ExtractedFilePoint {
  return { filepath: "", timestamp: "", latitude: 0, longitude: 0, metadata: "" };
}

export const ExtractedFilePoint = {
  encode(message: ExtractedFilePoint, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.filepath !== "") {
      writer.uint32(10).string(message.filepath);
    }
    if (message.timestamp !== "") {
      writer.uint32(18).string(message.timestamp);
    }
    if (message.latitude !== 0) {
      writer.uint32(29).float(message.latitude);
    }
    if (message.longitude !== 0) {
      writer.uint32(37).float(message.longitude);
    }
    if (message.metadata !== "") {
      writer.uint32(42).string(message.metadata);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExtractedFilePoint {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExtractedFilePoint();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.filepath = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.timestamp = reader.string();
          continue;
        case 3:
          if (tag !== 29) {
            break;
          }

          message.latitude = reader.float();
          continue;
        case 4:
          if (tag !== 37) {
            break;
          }

          message.longitude = reader.float();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.metadata = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<ExtractedFilePoint, Uint8Array>
  async *encodeTransform(
    source:
      | AsyncIterable<ExtractedFilePoint | ExtractedFilePoint[]>
      | Iterable<ExtractedFilePoint | ExtractedFilePoint[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [ExtractedFilePoint.encode(p).finish()];
        }
      } else {
        yield* [ExtractedFilePoint.encode(pkt as any).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, ExtractedFilePoint>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<ExtractedFilePoint> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [ExtractedFilePoint.decode(p)];
        }
      } else {
        yield* [ExtractedFilePoint.decode(pkt as any)];
      }
    }
  },

  fromJSON(object: any): ExtractedFilePoint {
    return {
      filepath: isSet(object.filepath) ? globalThis.String(object.filepath) : "",
      timestamp: isSet(object.timestamp) ? globalThis.String(object.timestamp) : "",
      latitude: isSet(object.latitude) ? globalThis.Number(object.latitude) : 0,
      longitude: isSet(object.longitude) ? globalThis.Number(object.longitude) : 0,
      metadata: isSet(object.metadata) ? globalThis.String(object.metadata) : "",
    };
  },

  toJSON(message: ExtractedFilePoint): unknown {
    const obj: any = {};
    if (message.filepath !== "") {
      obj.filepath = message.filepath;
    }
    if (message.timestamp !== "") {
      obj.timestamp = message.timestamp;
    }
    if (message.latitude !== 0) {
      obj.latitude = message.latitude;
    }
    if (message.longitude !== 0) {
      obj.longitude = message.longitude;
    }
    if (message.metadata !== "") {
      obj.metadata = message.metadata;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ExtractedFilePoint>, I>>(base?: I): ExtractedFilePoint {
    return ExtractedFilePoint.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ExtractedFilePoint>, I>>(object: I): ExtractedFilePoint {
    const message = createBaseExtractedFilePoint();
    message.filepath = object.filepath ?? "";
    message.timestamp = object.timestamp ?? "";
    message.latitude = object.latitude ?? 0;
    message.longitude = object.longitude ?? 0;
    message.metadata = object.metadata ?? "";
    return message;
  },
};

function createBaseExtractedFilePoints(): ExtractedFilePoints {
  return { points: [] };
}

export const ExtractedFilePoints = {
  encode(message: ExtractedFilePoints, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.points) {
      ExtractedFilePoint.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExtractedFilePoints {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExtractedFilePoints();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.points.push(ExtractedFilePoint.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  // encodeTransform encodes a source of message objects.
  // Transform<ExtractedFilePoints, Uint8Array>
  async *encodeTransform(
    source:
      | AsyncIterable<ExtractedFilePoints | ExtractedFilePoints[]>
      | Iterable<ExtractedFilePoints | ExtractedFilePoints[]>,
  ): AsyncIterable<Uint8Array> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [ExtractedFilePoints.encode(p).finish()];
        }
      } else {
        yield* [ExtractedFilePoints.encode(pkt as any).finish()];
      }
    }
  },

  // decodeTransform decodes a source of encoded messages.
  // Transform<Uint8Array, ExtractedFilePoints>
  async *decodeTransform(
    source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>,
  ): AsyncIterable<ExtractedFilePoints> {
    for await (const pkt of source) {
      if (globalThis.Array.isArray(pkt)) {
        for (const p of (pkt as any)) {
          yield* [ExtractedFilePoints.decode(p)];
        }
      } else {
        yield* [ExtractedFilePoints.decode(pkt as any)];
      }
    }
  },

  fromJSON(object: any): ExtractedFilePoints {
    return {
      points: globalThis.Array.isArray(object?.points)
        ? object.points.map((e: any) => ExtractedFilePoint.fromJSON(e))
        : [],
    };
  },

  toJSON(message: ExtractedFilePoints): unknown {
    const obj: any = {};
    if (message.points?.length) {
      obj.points = message.points.map((e) => ExtractedFilePoint.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ExtractedFilePoints>, I>>(base?: I): ExtractedFilePoints {
    return ExtractedFilePoints.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ExtractedFilePoints>, I>>(object: I): ExtractedFilePoints {
    const message = createBaseExtractedFilePoints();
    message.points = object.points?.map((e) => ExtractedFilePoint.fromPartial(e)) || [];
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
