# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: beneath/proto/gateway.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from beneath.proto import engine_pb2 as beneath_dot_proto_dot_engine__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='beneath/proto/gateway.proto',
  package='proto',
  syntax='proto3',
  serialized_options=_b('\n\025network.beneath.protoB\014BeneathProtoP\001'),
  serialized_pb=_b('\n\x1b\x62\x65neath/proto/gateway.proto\x12\x05proto\x1a\x1a\x62\x65neath/proto/engine.proto\"V\n\x12ReadRecordsRequest\x12\x13\n\x0binstance_id\x18\x01 \x01(\x0c\x12\r\n\x05limit\x18\x02 \x01(\x05\x12\r\n\x05where\x18\x03 \x01(\t\x12\r\n\x05\x61\x66ter\x18\x04 \x01(\t\"5\n\x13ReadRecordsResponse\x12\x1e\n\x07records\x18\x01 \x03(\x0b\x32\r.proto.Record\"N\n\x18ReadLatestRecordsRequest\x12\x13\n\x0binstance_id\x18\x01 \x01(\x0c\x12\r\n\x05limit\x18\x02 \x01(\x05\x12\x0e\n\x06\x62\x65\x66ore\x18\x03 \x01(\x03\"\x16\n\x14WriteRecordsResponse\"A\n\x14StreamDetailsRequest\x12\x14\n\x0cproject_name\x18\x01 \x01(\t\x12\x13\n\x0bstream_name\x18\x02 \x01(\t\"\xdd\x01\n\x15StreamDetailsResponse\x12\x1b\n\x13\x63urrent_instance_id\x18\x01 \x01(\x0c\x12\x12\n\nproject_id\x18\x02 \x01(\x0c\x12\x14\n\x0cproject_name\x18\x03 \x01(\t\x12\x13\n\x0bstream_name\x18\x04 \x01(\t\x12\x12\n\nkey_fields\x18\x05 \x03(\t\x12\x13\n\x0b\x61vro_schema\x18\x06 \x01(\t\x12\x0e\n\x06public\x18\x07 \x01(\x08\x12\x10\n\x08\x65xternal\x18\x08 \x01(\x08\x12\r\n\x05\x62\x61tch\x18\t \x01(\x08\x12\x0e\n\x06manual\x18\n \x01(\x08\"7\n\nClientPing\x12\x11\n\tclient_id\x18\x01 \x01(\t\x12\x16\n\x0e\x63lient_version\x18\x02 \x01(\t\"P\n\nClientPong\x12\x15\n\rauthenticated\x18\x01 \x01(\x08\x12\x0e\n\x06status\x18\x02 \x01(\t\x12\x1b\n\x13recommended_version\x18\x03 \x01(\t2\xfb\x02\n\x07Gateway\x12\x46\n\x0bReadRecords\x12\x19.proto.ReadRecordsRequest\x1a\x1a.proto.ReadRecordsResponse\"\x00\x12R\n\x11ReadLatestRecords\x12\x1f.proto.ReadLatestRecordsRequest\x1a\x1a.proto.ReadRecordsResponse\"\x00\x12I\n\x0cWriteRecords\x12\x1a.proto.WriteRecordsRequest\x1a\x1b.proto.WriteRecordsResponse\"\x00\x12O\n\x10GetStreamDetails\x12\x1b.proto.StreamDetailsRequest\x1a\x1c.proto.StreamDetailsResponse\"\x00\x12\x38\n\x0eSendClientPing\x12\x11.proto.ClientPing\x1a\x11.proto.ClientPong\"\x00\x42\'\n\x15network.beneath.protoB\x0c\x42\x65neathProtoP\x01\x62\x06proto3')
  ,
  dependencies=[beneath_dot_proto_dot_engine__pb2.DESCRIPTOR,])




_READRECORDSREQUEST = _descriptor.Descriptor(
  name='ReadRecordsRequest',
  full_name='proto.ReadRecordsRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='instance_id', full_name='proto.ReadRecordsRequest.instance_id', index=0,
      number=1, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=_b(""),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='limit', full_name='proto.ReadRecordsRequest.limit', index=1,
      number=2, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='where', full_name='proto.ReadRecordsRequest.where', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='after', full_name='proto.ReadRecordsRequest.after', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=66,
  serialized_end=152,
)


_READRECORDSRESPONSE = _descriptor.Descriptor(
  name='ReadRecordsResponse',
  full_name='proto.ReadRecordsResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='records', full_name='proto.ReadRecordsResponse.records', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=154,
  serialized_end=207,
)


_READLATESTRECORDSREQUEST = _descriptor.Descriptor(
  name='ReadLatestRecordsRequest',
  full_name='proto.ReadLatestRecordsRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='instance_id', full_name='proto.ReadLatestRecordsRequest.instance_id', index=0,
      number=1, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=_b(""),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='limit', full_name='proto.ReadLatestRecordsRequest.limit', index=1,
      number=2, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='before', full_name='proto.ReadLatestRecordsRequest.before', index=2,
      number=3, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=209,
  serialized_end=287,
)


_WRITERECORDSRESPONSE = _descriptor.Descriptor(
  name='WriteRecordsResponse',
  full_name='proto.WriteRecordsResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=289,
  serialized_end=311,
)


_STREAMDETAILSREQUEST = _descriptor.Descriptor(
  name='StreamDetailsRequest',
  full_name='proto.StreamDetailsRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='project_name', full_name='proto.StreamDetailsRequest.project_name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='stream_name', full_name='proto.StreamDetailsRequest.stream_name', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=313,
  serialized_end=378,
)


_STREAMDETAILSRESPONSE = _descriptor.Descriptor(
  name='StreamDetailsResponse',
  full_name='proto.StreamDetailsResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='current_instance_id', full_name='proto.StreamDetailsResponse.current_instance_id', index=0,
      number=1, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=_b(""),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='project_id', full_name='proto.StreamDetailsResponse.project_id', index=1,
      number=2, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=_b(""),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='project_name', full_name='proto.StreamDetailsResponse.project_name', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='stream_name', full_name='proto.StreamDetailsResponse.stream_name', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='key_fields', full_name='proto.StreamDetailsResponse.key_fields', index=4,
      number=5, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='avro_schema', full_name='proto.StreamDetailsResponse.avro_schema', index=5,
      number=6, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='public', full_name='proto.StreamDetailsResponse.public', index=6,
      number=7, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='external', full_name='proto.StreamDetailsResponse.external', index=7,
      number=8, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='batch', full_name='proto.StreamDetailsResponse.batch', index=8,
      number=9, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='manual', full_name='proto.StreamDetailsResponse.manual', index=9,
      number=10, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=381,
  serialized_end=602,
)


_CLIENTPING = _descriptor.Descriptor(
  name='ClientPing',
  full_name='proto.ClientPing',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='client_id', full_name='proto.ClientPing.client_id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='client_version', full_name='proto.ClientPing.client_version', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=604,
  serialized_end=659,
)


_CLIENTPONG = _descriptor.Descriptor(
  name='ClientPong',
  full_name='proto.ClientPong',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='authenticated', full_name='proto.ClientPong.authenticated', index=0,
      number=1, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='status', full_name='proto.ClientPong.status', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='recommended_version', full_name='proto.ClientPong.recommended_version', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=661,
  serialized_end=741,
)

_READRECORDSRESPONSE.fields_by_name['records'].message_type = beneath_dot_proto_dot_engine__pb2._RECORD
DESCRIPTOR.message_types_by_name['ReadRecordsRequest'] = _READRECORDSREQUEST
DESCRIPTOR.message_types_by_name['ReadRecordsResponse'] = _READRECORDSRESPONSE
DESCRIPTOR.message_types_by_name['ReadLatestRecordsRequest'] = _READLATESTRECORDSREQUEST
DESCRIPTOR.message_types_by_name['WriteRecordsResponse'] = _WRITERECORDSRESPONSE
DESCRIPTOR.message_types_by_name['StreamDetailsRequest'] = _STREAMDETAILSREQUEST
DESCRIPTOR.message_types_by_name['StreamDetailsResponse'] = _STREAMDETAILSRESPONSE
DESCRIPTOR.message_types_by_name['ClientPing'] = _CLIENTPING
DESCRIPTOR.message_types_by_name['ClientPong'] = _CLIENTPONG
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

ReadRecordsRequest = _reflection.GeneratedProtocolMessageType('ReadRecordsRequest', (_message.Message,), {
  'DESCRIPTOR' : _READRECORDSREQUEST,
  '__module__' : 'beneath.proto.gateway_pb2'
  # @@protoc_insertion_point(class_scope:proto.ReadRecordsRequest)
  })
_sym_db.RegisterMessage(ReadRecordsRequest)

ReadRecordsResponse = _reflection.GeneratedProtocolMessageType('ReadRecordsResponse', (_message.Message,), {
  'DESCRIPTOR' : _READRECORDSRESPONSE,
  '__module__' : 'beneath.proto.gateway_pb2'
  # @@protoc_insertion_point(class_scope:proto.ReadRecordsResponse)
  })
_sym_db.RegisterMessage(ReadRecordsResponse)

ReadLatestRecordsRequest = _reflection.GeneratedProtocolMessageType('ReadLatestRecordsRequest', (_message.Message,), {
  'DESCRIPTOR' : _READLATESTRECORDSREQUEST,
  '__module__' : 'beneath.proto.gateway_pb2'
  # @@protoc_insertion_point(class_scope:proto.ReadLatestRecordsRequest)
  })
_sym_db.RegisterMessage(ReadLatestRecordsRequest)

WriteRecordsResponse = _reflection.GeneratedProtocolMessageType('WriteRecordsResponse', (_message.Message,), {
  'DESCRIPTOR' : _WRITERECORDSRESPONSE,
  '__module__' : 'beneath.proto.gateway_pb2'
  # @@protoc_insertion_point(class_scope:proto.WriteRecordsResponse)
  })
_sym_db.RegisterMessage(WriteRecordsResponse)

StreamDetailsRequest = _reflection.GeneratedProtocolMessageType('StreamDetailsRequest', (_message.Message,), {
  'DESCRIPTOR' : _STREAMDETAILSREQUEST,
  '__module__' : 'beneath.proto.gateway_pb2'
  # @@protoc_insertion_point(class_scope:proto.StreamDetailsRequest)
  })
_sym_db.RegisterMessage(StreamDetailsRequest)

StreamDetailsResponse = _reflection.GeneratedProtocolMessageType('StreamDetailsResponse', (_message.Message,), {
  'DESCRIPTOR' : _STREAMDETAILSRESPONSE,
  '__module__' : 'beneath.proto.gateway_pb2'
  # @@protoc_insertion_point(class_scope:proto.StreamDetailsResponse)
  })
_sym_db.RegisterMessage(StreamDetailsResponse)

ClientPing = _reflection.GeneratedProtocolMessageType('ClientPing', (_message.Message,), {
  'DESCRIPTOR' : _CLIENTPING,
  '__module__' : 'beneath.proto.gateway_pb2'
  # @@protoc_insertion_point(class_scope:proto.ClientPing)
  })
_sym_db.RegisterMessage(ClientPing)

ClientPong = _reflection.GeneratedProtocolMessageType('ClientPong', (_message.Message,), {
  'DESCRIPTOR' : _CLIENTPONG,
  '__module__' : 'beneath.proto.gateway_pb2'
  # @@protoc_insertion_point(class_scope:proto.ClientPong)
  })
_sym_db.RegisterMessage(ClientPong)


DESCRIPTOR._options = None

_GATEWAY = _descriptor.ServiceDescriptor(
  name='Gateway',
  full_name='proto.Gateway',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=744,
  serialized_end=1123,
  methods=[
  _descriptor.MethodDescriptor(
    name='ReadRecords',
    full_name='proto.Gateway.ReadRecords',
    index=0,
    containing_service=None,
    input_type=_READRECORDSREQUEST,
    output_type=_READRECORDSRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='ReadLatestRecords',
    full_name='proto.Gateway.ReadLatestRecords',
    index=1,
    containing_service=None,
    input_type=_READLATESTRECORDSREQUEST,
    output_type=_READRECORDSRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='WriteRecords',
    full_name='proto.Gateway.WriteRecords',
    index=2,
    containing_service=None,
    input_type=beneath_dot_proto_dot_engine__pb2._WRITERECORDSREQUEST,
    output_type=_WRITERECORDSRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='GetStreamDetails',
    full_name='proto.Gateway.GetStreamDetails',
    index=3,
    containing_service=None,
    input_type=_STREAMDETAILSREQUEST,
    output_type=_STREAMDETAILSRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='SendClientPing',
    full_name='proto.Gateway.SendClientPing',
    index=4,
    containing_service=None,
    input_type=_CLIENTPING,
    output_type=_CLIENTPONG,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_GATEWAY)

DESCRIPTOR.services_by_name['Gateway'] = _GATEWAY

# @@protoc_insertion_point(module_scope)
