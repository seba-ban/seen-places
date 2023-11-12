from pathlib import Path

import pika.channel
import pika.spec

from .logger import logger
from .processing import garmin_data_to_points, parse_activity
from .proto.events_pb2 import (
    ExtractedFilePoint,
    ExtractedFilePoints,
    ProcessFileRequest,
)
from .settings import broker_settings, settings


def callback(
    ch: pika.channel.Channel,
    method: pika.spec.Basic.Deliver,
    properties: pika.spec.BasicProperties,
    body: bytes,
):
    req = ProcessFileRequest()
    try:
        req.ParseFromString(body)
    except Exception as e:
        logger.error(f"Failed to parse message: {e}")
        return

    fit_file = Path(settings.local_storage_files_dir) / req.filepath
    if not fit_file.exists():
        logger.error(f"File {fit_file} does not exist")
        return

    parsed = parse_activity(fit_file)

    if parsed is None:
        logger.error(f"Failed to parse file {fit_file}")
        return

    points: list[ExtractedFilePoint] = []
    for point in garmin_data_to_points(parsed):
        point.filepath = req.filepath
        points.append(point)
        if len(points) >= broker_settings.points_buffer_size:
            points_msg = ExtractedFilePoints()
            points_msg.points.extend(points)
            ch.basic_publish(
                exchange="",
                routing_key=broker_settings.target_queue,
                body=points_msg.SerializeToString(),
            )
            points = []

    if len(points) > 0:
        points_msg = ExtractedFilePoints()
        points_msg.points.extend(points)
        ch.basic_publish(
            exchange="",
            routing_key=broker_settings.target_queue,
            body=points_msg.SerializeToString(),
        )
