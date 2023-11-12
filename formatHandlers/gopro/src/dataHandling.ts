import { exec as execCb, spawn } from "node:child_process";
import * as fs from "node:fs/promises";
import { tmpdir } from "node:os";
import * as path from "node:path";
import * as util from "node:util";

import { goProTelemetry } from "gopro-telemetry";

import {
  ExtractedFilePoint,
  ExtractedFilePoints,
  ProcessFileRequest,
} from "./proto/events";
import { channelSettings, localStorageSettings } from "./settings";

const exec = util.promisify(execCb);

type FFProbeOutput = { [key: string]: any };

async function probeFile(filePath: string): Promise<FFProbeOutput> {
  const { stdout } = await exec(
    `ffprobe -v quiet -print_format json -show_format -show_streams ${filePath}`,
  );

  return JSON.parse(stdout);
}

/**
 *
 * @param ffprobeOutput
 * @returns -1 if not found, otherwise the index of the stream
 */
function getMetadataStreamIndex(ffprobeOutput: FFProbeOutput): number {
  if (!Array.isArray(ffprobeOutput.streams)) {
    return -1;
  }

  for (const stream of ffprobeOutput.streams) {
    if (!stream.tags || !Number.isInteger(stream.index)) {
      continue;
    }

    const handlerName = stream.tags.handler_name;
    if (typeof handlerName !== "string" || !handlerName) {
      continue;
    }

    if (handlerName.includes("GoPro MET")) return stream.index;
  }

  return -1;
}

async function extractTelemetry(filePath: string, streamIndex: number) {
  const tmp = await fs.mkdtemp(path.join(tmpdir(), "gopro-"));
  const tmpFile = path.join(tmp, path.basename(filePath));

  const proc = spawn(
    "ffmpeg",
    [
      "-y",
      "-i",
      filePath,
      "-codec",
      "copy",
      "-map",
      `0:${streamIndex}`,
      "-f",
      "rawvideo",
      tmpFile,
    ],
    {
      // stdio: "inherit",
    },
  );

  const wait = () =>
    new Promise((resolve, reject) => {
      proc.on("exit", (code) => {
        if (code === 0) {
          resolve(code);
        } else {
          reject(code);
        }
      });
    });
  await wait();

  const telemetry = await goProTelemetry(
    { rawData: await fs.readFile(tmpFile) },
    {},
  );

  await fs.rm(tmp, { recursive: true });

  return telemetry;
}

function* getPointsFromTelemetry(
  filepath: string,
  telemetry: { [key: string]: any },
) {
  let zeroes = 0;
  for (const obj of Object.values(telemetry)) {
    const samples = obj?.streams?.GPS9?.samples;
    if (!Array.isArray(samples)) {
      console.log("samples not an array");
      continue;
    }

    // "name": "GPS (Lat., Long., Alt., 2D, 3D, days, secs, DOP, fix)",
    // "units": [
    //   "deg",
    //   "deg",
    //   "m",
    //   "m/s",
    //   "m/s",
    //   "",
    //   "s",
    //   "",
    //   ""
    // ]
    // {
    //     "value": [
    //       52.2202675,
    //       21.0074384,
    //       102.311,
    //       1.047,
    //       1.16,
    //       8712,
    //       75025.8,
    //       2.02,
    //       3
    //     ],
    //     "cts": 0,
    //     "date": "2023-11-08T20:50:25.800Z",
    //     "sticky": {
    //       "altitude system": "MSLV"
    //     }
    //   },
    for (const sample of samples) {
      if (!Array.isArray(sample.value) || sample.value.length !== 9) {
        console.log("invalid sample value");
        continue;
      }

      // TODO: check if valid string with timezone
      if (!(sample.date instanceof Date)) {
        console.log("invalid sample date");
        continue;
      }

      const lat = sample.value[0];
      const lon = sample.value[1];

      if (!Number.isFinite(lat) || !Number.isFinite(lon)) {
        console.log("invalid values", lat, lon);
        continue;
      }

      if (lat === 0 && lon === 0) {
        // console.log("both lat and lon are 0, skipping...");
        zeroes++;
        continue;
      }

      yield ExtractedFilePoint.create({
        filepath,
        latitude: sample.value[0],
        longitude: sample.value[1],
        timestamp: sample.date.toISOString(),
      });
    }
  }
  console.log("zeroes", zeroes);
}

export async function* processRequest(req: ProcessFileRequest) {
  const filePath = path.join(localStorageSettings.filesDir, req.filepath);

  const meta = await probeFile(filePath);
  const telemetry = await extractTelemetry(
    filePath,
    getMetadataStreamIndex(meta),
  );

  let buff: ExtractedFilePoint[] = [];
  for (const point of getPointsFromTelemetry(req.filepath, telemetry)) {
    buff.push(point);
    if (buff.length >= channelSettings.pointsBufferSize) {
      console.log("yielding", buff.length);
      yield ExtractedFilePoints.create({ points: buff });
      buff = [];
    }
  }

  if (buff.length > 0) {
    console.log("yielding", buff.length);
    yield ExtractedFilePoints.create({ points: buff });
  }
}
