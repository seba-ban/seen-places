"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder

_sym_db = _symbol_database.Default()
DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(
    b'\n\x0cevents.proto"&\n\x12ProcessFileRequest\x12\x10\n\x08filepath\x18\x01 \x01(\t"p\n\x12ExtractedFilePoint\x12\x10\n\x08filepath\x18\x01 \x01(\t\x12\x11\n\ttimestamp\x18\x02 \x01(\t\x12\x10\n\x08latitude\x18\x03 \x01(\x02\x12\x11\n\tlongitude\x18\x04 \x01(\x02\x12\x10\n\x08metadata\x18\x05 \x01(\t":\n\x13ExtractedFilePoints\x12#\n\x06points\x18\x01 \x03(\x0b2\x13.ExtractedFilePointB(Z&github.com/seba-ban/seen-places/eventsb\x06proto3'
)
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, "events_pb2", globals())
if _descriptor._USE_C_DESCRIPTORS == False:
    DESCRIPTOR._options = None
    DESCRIPTOR._serialized_options = b"Z&github.com/seba-ban/seen-places/events"
    _PROCESSFILEREQUEST._serialized_start = 16
    _PROCESSFILEREQUEST._serialized_end = 54
    _EXTRACTEDFILEPOINT._serialized_start = 56
    _EXTRACTEDFILEPOINT._serialized_end = 168
    _EXTRACTEDFILEPOINTS._serialized_start = 170
    _EXTRACTEDFILEPOINTS._serialized_end = 228
