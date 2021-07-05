package protoquery

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/sorcererxw/protoquery/internal/genid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Encoder struct{}

// Encode converts proto message to query.
func (e *Encoder) Encode(msg proto.Message) url.Values {
	query := make(url.Values)
	msg.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, val protoreflect.Value) bool {
		if !val.IsValid() {
			return true
		}
		switch {
		case fd.IsMap():
			s := e.encodeMap(fd, val.Map())
			if s != "" {
				query.Set(fd.TextName(), s)
			}
		case fd.IsList():
			s := e.encodeList(fd, val.List())
			if s != "" {
				query.Set(fd.TextName(), s)
			}
		default:
			s := e.encodeSingular(fd, val)
			if s != "" {
				query.Set(fd.TextName(), s)
			}
		}
		return true
	})
	return query
}

func (e *Encoder) encodeMap(fd protoreflect.FieldDescriptor, m protoreflect.Map) string {
	return ""
}

func (e *Encoder) encodeList(fd protoreflect.FieldDescriptor, list protoreflect.List) string {
	arr := make([]string, 0, list.Len())
	for i := 0; i < list.Len(); i++ {
		s := e.encodeSingular(fd, list.Get(i))
		if s != "" {
			arr = append(arr, s)
		}
	}
	return strings.Join(arr, ",")
}

func (e *Encoder) encodeSingular(fd protoreflect.FieldDescriptor, val protoreflect.Value) string {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		return strconv.FormatBool(val.Bool())
	case protoreflect.DoubleKind, protoreflect.FloatKind:
		return strconv.FormatFloat(val.Float(), 'f', 6, 64)
	case protoreflect.BytesKind:
		return base64.URLEncoding.EncodeToString(val.Bytes())
	case protoreflect.StringKind:
		return val.String()
	case protoreflect.Int64Kind, protoreflect.Int32Kind,
		protoreflect.Sint64Kind, protoreflect.Sint32Kind,
		protoreflect.Sfixed64Kind, protoreflect.Sfixed32Kind:
		return strconv.FormatInt(val.Int(), 10)
	case protoreflect.Uint64Kind, protoreflect.Uint32Kind,
		protoreflect.Fixed64Kind, protoreflect.Fixed32Kind:
		return strconv.FormatUint(val.Uint(), 10)
	case protoreflect.EnumKind:
		return string(fd.Enum().Values().ByNumber(val.Enum()).Name())
	case protoreflect.MessageKind:
		return e.encodeMessage(val.Message())
	}
	return ""
}

func (e *Encoder) encodeMessage(m protoreflect.Message) string {
	switch m.Descriptor().FullName() {
	case genid.Timestamp_message_fullname:
		return e.encodeTimestamp(m)
	case genid.Duration_message_fullname:
		return e.encodeDuration(m)
	}
	return ""
}

func (e *Encoder) encodeTimestamp(m protoreflect.Message) string {
	fds := m.Descriptor().Fields()
	fdSeconds := fds.ByNumber(genid.Timestamp_Seconds_field_number)
	fdNanos := fds.ByNumber(genid.Timestamp_Nanos_field_number)
	return fmt.Sprintf("%d.%09d", m.Get(fdSeconds).Int(), m.Get(fdNanos).Int())
}

func (e *Encoder) encodeDuration(m protoreflect.Message) string {
	fds := m.Descriptor().Fields()
	fdSeconds := fds.ByNumber(genid.Duration_Seconds_field_number)
	fdNanos := fds.ByNumber(genid.Duration_Nanos_field_number)
	return fmt.Sprintf("%d.%09d", m.Get(fdSeconds).Int(), m.Get(fdNanos).Int())
}
