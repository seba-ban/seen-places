from typing import ClassVar as _ClassVar
from typing import Iterable as _Iterable
from typing import Mapping as _Mapping
from typing import Optional as _Optional
from typing import Union as _Union

from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf.internal import containers as _containers

DESCRIPTOR: _descriptor.FileDescriptor

class ExtractedFilePoint(_message.Message):
    __slots__ = ["filepath", "latitude", "longitude", "metadata", "timestamp"]
    FILEPATH_FIELD_NUMBER: _ClassVar[int]
    LATITUDE_FIELD_NUMBER: _ClassVar[int]
    LONGITUDE_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    filepath: str
    latitude: float
    longitude: float
    metadata: str
    timestamp: str

    def __init__(
        self,
        filepath: _Optional[str] = ...,
        timestamp: _Optional[str] = ...,
        latitude: _Optional[float] = ...,
        longitude: _Optional[float] = ...,
        metadata: _Optional[str] = ...,
    ) -> None: ...

class ExtractedFilePoints(_message.Message):
    __slots__ = ["points"]
    POINTS_FIELD_NUMBER: _ClassVar[int]
    points: _containers.RepeatedCompositeFieldContainer[ExtractedFilePoint]

    def __init__(
        self, points: _Optional[_Iterable[_Union[ExtractedFilePoint, _Mapping]]] = ...
    ) -> None: ...

class ProcessFileRequest(_message.Message):
    __slots__ = ["filepath"]
    FILEPATH_FIELD_NUMBER: _ClassVar[int]
    filepath: str

    def __init__(self, filepath: _Optional[str] = ...) -> None: ...
