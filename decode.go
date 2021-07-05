package protoquery

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sorcererxw/protoquery/internal/genid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var prefNil = protoreflect.ValueOf(nil)

type Decoder struct{}

// Decode converts query to proto message.
func (d *Decoder) Decode(query url.Values, msg proto.Message) error {
	pr := msg.ProtoReflect()
	for key := range query {
		val := query.Get(key)

		fd := pr.Descriptor().Fields().ByName(protoreflect.Name(key))
		if fd == nil {
			continue
		}
		switch {
		case fd.IsMap():
			if err := d.decodeMap(fd, val, pr.Mutable(fd).Map()); err != nil {
				return err
			}
		case fd.IsList():
			if err := d.decodeList(fd, val, pr.Mutable(fd).List()); err != nil {
				return err
			}
		default:
			val, err := d.decodeSingular(fd, val, pr.Get(fd))
			if err != nil {
				return err
			}
			if val.IsValid() {
				pr.Set(fd, val)
			}
		}
	}
	return nil
}

func (d *Decoder) decodeMap(fd protoreflect.FieldDescriptor, param string, m protoreflect.Map) error {
	return nil
}

func (d *Decoder) decodeList(fd protoreflect.FieldDescriptor, param string, list protoreflect.List) error {
	items := strings.Split(param, ",")
	for _, it := range items {
		val, err := d.decodeSingular(fd, it, list.NewElement())
		if err != nil {
			return err
		}
		if !val.IsValid() {
			continue
		}
		list.Append(val)
	}
	return nil
}

func (d *Decoder) decodeSingular(fd protoreflect.FieldDescriptor, param string, m protoreflect.Value) (protoreflect.Value, error) {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		v, err := strconv.ParseBool(param)
		if err != nil {
			return prefNil, err
		}
		return protoreflect.ValueOfBool(v), nil
	case protoreflect.DoubleKind:
		v, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return prefNil, err
		}
		return protoreflect.ValueOfFloat64(v), nil
	case protoreflect.FloatKind:
		v, err := strconv.ParseFloat(param, 32)
		if err != nil {
			return prefNil, err
		}
		return protoreflect.ValueOfFloat32(float32(v)), nil
	case protoreflect.BytesKind:
		v, err := base64.URLEncoding.DecodeString(param)
		if err != nil {
			return prefNil, err
		}
		return protoreflect.ValueOfBytes(v), nil
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(param), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		v, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return prefNil, err
		}
		return protoreflect.ValueOfInt64(v), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		v, err := strconv.ParseInt(param, 10, 32)
		if err != nil {
			return prefNil, err
		}
		return protoreflect.ValueOfInt32(int32(v)), nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		v, err := strconv.ParseUint(param, 10, 64)
		if err != nil {
			return prefNil, err
		}
		return protoreflect.ValueOfUint64(v), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		v, err := strconv.ParseUint(param, 10, 32)
		if err != nil {
			return prefNil, err
		}
		return protoreflect.ValueOfUint32(uint32(v)), nil
	case protoreflect.EnumKind:
		v := fd.Enum().Values().ByName(protoreflect.Name(param))
		if v == nil {
			return prefNil, nil
		}
		return protoreflect.ValueOfEnum(v.Number()), nil
	case protoreflect.MessageKind:
		v, err := d.decodeMessage(param, m.Message())
		if err != nil {
			return prefNil, nil
		}
		return v, nil
	}
	return prefNil, nil
}

func (d *Decoder) decodeMessage(param string, m protoreflect.Message) (protoreflect.Value, error) {
	fmt.Println(m.Descriptor().FullName())
	switch m.Descriptor().FullName() {
	case genid.Timestamp_message_fullname:
		return d.decodeTimestamp(param)
	case genid.Duration_message_fullname:
		return d.decodeDuration(param)
	}
	return prefNil, nil
}

func (d *Decoder) decodeTimestamp(param string) (protoreflect.Value, error) {
	sec, nsec, err := d.decodeUnix(param)
	if err != nil {
		return prefNil, err
	}
	t := time.Unix(sec, nsec)
	return protoreflect.ValueOfMessage(timestamppb.New(t).ProtoReflect()), nil
}

func (d *Decoder) decodeDuration(param string) (protoreflect.Value, error) {
	sec, nsec, err := d.decodeUnix(param)
	if err != nil {
		return prefNil, err
	}
	dur := time.Duration(sec)*time.Second + time.Duration(nsec)
	return protoreflect.ValueOfMessage(durationpb.New(dur).ProtoReflect()), nil
}

func (d *Decoder) decodeUnix(param string) (sec int64, nsec int64, err error) {
	f := strings.SplitN(param, ".", 2)
	if len(f) >= 1 && f[0] != "" {
		sec, err = strconv.ParseInt(f[0], 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}
	if len(f) >= 2 && f[1] != "" {
		s := f[1]
		switch {
		case len(s) > 9:
			return 0, 0, errors.New("invalid nsec format")
		case len(s) < 9:
			s += strings.Repeat("0", 9-len(s))
		}
		nsec, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}
	return
}
