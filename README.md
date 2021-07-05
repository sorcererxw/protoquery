# protoquery

protoquery encodes and decodes protocol buffer messages as URL search param format.

## Type Mapping

Following table shows how protoquery encoding proto message to url query. 

|proto3|Example|Notes|
|---|---|---|
|message|/||
|enum|arg=ENUM||
|map<K,V>|/||
|repeated V|arg=v1,v2,v3||
|bool|arg=true||
|string|arg=hello||
|bytes|arg=YWJjMTIzIT8kKiYoKSctPUB+|bytes payload would be encoded to URL-safe base64 string.|
|Int32,fixed32,uint32,int64,fixed64,uint64|arg=12||
| float,double                              |arg=1.1||
|Any|/||
|Timestamp|arg=3463213.000123124|Unix timestamp|
|Duration|arg=1231.001231000| Seconds of the duration                                   |
|Struct|/||
|Wrapper types|/||
|Field Mask|/||
|ListValue|/||
|Value|/||
|NullValue|/||
|Empty|/||
