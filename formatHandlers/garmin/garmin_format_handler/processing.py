from datetime import datetime, timezone

from garmin_fit_sdk import Decoder, Stream
from pydantic import BaseModel

from .proto.events_pb2 import ExtractedFilePoint


def normalize_garmin_lat_lon(val: int) -> float:
    # https://gis.stackexchange.com/questions/371656/garmin-fit-coordinate-system
    # TODO: consider using a library for this, maybe corner cases?
    return val / (2**32 / 360)


def get_activity_type(messages: dict) -> str | None:
    try:
        return messages["sport_mesgs"][0]["sport"].lower()
    except:
        return None


def get_activity_time(messages: dict) -> datetime | None:
    try:
        return messages["file_id_mesgs"][0]["time_created"]
    except:
        return None


def process_datetime(dt: datetime) -> str:
    if dt.tzinfo is None:
        dt = dt.replace(tzinfo=timezone.utc)
    return dt.isoformat()


class Point(BaseModel):
    latitude: float
    longitute: float
    timestamp: datetime

    @classmethod
    def from_record(cls, record: dict):
        try:
            return cls(
                latitude=normalize_garmin_lat_lon(record["position_lat"]),
                longitute=normalize_garmin_lat_lon(record["position_long"]),
                timestamp=process_datetime(record["timestamp"]),
            )
        except:
            # TODO: log?
            return None


class GarminData(BaseModel):
    activity_type: str
    activity_time: datetime
    points: list[Point]


def parse_activity(
    fit_path: str,
) -> GarminData | None:
    stream = Stream.from_file(fit_path)
    decoder = Decoder(stream)
    messages, errors = decoder.read()

    if errors:
        print(errors)

    activity_type = get_activity_type(messages)
    if activity_type is None:
        return None

    activity_time = get_activity_time(messages)
    if activity_time is None:
        return None

    points = [
        parsed
        for record in messages["record_mesgs"]
        if (parsed := Point.from_record(record)) is not None
    ]

    return GarminData(
        activity_type=activity_type,
        activity_time=activity_time,
        points=points,
    )


def garmin_data_to_points(data: GarminData):
    for point in data.points:
        yield ExtractedFilePoint(
            latitude=point.latitude,
            longitude=point.longitute,
            timestamp=process_datetime(point.timestamp),
        )
