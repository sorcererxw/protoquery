package protoquery

import (
	"net/url"
	"reflect"
	"testing"
	"time"

test	"github.com/sorcererxw/protoquery/testdata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEncoder_Encode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args proto.Message
		want url.Values
	}{
		{
			name: "invalid",
			args: &test.Request{},
			want: url.Values{},
		},
		{
			name: "invalid",
			args: &test.Request{
				String_: "",
			},
			want: url.Values{},
		},
		{
			name: "int64",
			args: &test.Request{
				Int64: 64,
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("int64", "64")
				return v
			}(),
		},
		{
			name: "int32",
			args: &test.Request{
				Int32: 32,
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("int32", "32")
				return v
			}(),
		},
		{
			name: "sint32",
			args: &test.Request{
				Sint32: 32,
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("sint32", "32")
				return v
			}(),
		},
		{
			name: "sint64",
			args: &test.Request{
				Sint64: 64,
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("sint64", "64")
				return v
			}(),
		},
		{
			name: "fixed64",
			args: &test.Request{
				Fixed64: 64,
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("fixed64", "64")
				return v
			}(),
		},
		{
			name: "fixed32",
			args: &test.Request{
				Fixed32: 32,
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("fixed32", "32")
				return v
			}(),
		},
		{
			name: "sfixed64",
			args: &test.Request{
				Sfixed64: 64,
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("sfixed64", "64")
				return v
			}(),
		},
		{
			name: "sfixed32",
			args: &test.Request{
				Sfixed32: 32,
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("sfixed32", "32")
				return v
			}(),
		},
		{
			name: "unint64",
			args: &test.Request{
				Uint64: 64,
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("uint64", "64")
				return v
			}(),
		},
		{
			name: "unint32",
			args: &test.Request{
				Uint32: 32,
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("uint32", "32")
				return v
			}(),
		},
		{
			name: "bool",
			args: &test.Request{
				Bool: true,
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("bool", "true")
				return v
			}(),
		},
		{
			name: "string",
			args: &test.Request{
				String_: "str",
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("string", "str")
				return v
			}(),
		},
		{
			name: "double",
			args: &test.Request{
				Double: 1.1,
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("double", "1.100000")
				return v
			}(),
		},
		{
			name: "float",
			args: &test.Request{
				Float: 1.1,
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("float", "1.100000")
				return v
			}(),
		},
		{
			name: "enum",
			args: &test.Request{
				Enum: test.Request_ENUM1,
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("enum", "ENUM1")
				return v
			}(),
		},
		{
			name: "repeated enum",
			args: &test.Request{
				RepeatedEnum: []test.Request_Enum{test.Request_ENUM1, test.Request_ENUM2},
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("repeatedEnum", "ENUM1,ENUM2")
				return v
			}(),
		},
		{
			name: "repeated empty",
			args: &test.Request{
				RepeatedString: []string{},
			},
			want: func() url.Values {
				v := url.Values{}
				return v
			}(),
		},
		{
			name: "timestamp",
			args: &test.Request{
				Timestamp: timestamppb.New(time.Unix(123, 123)),
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("timestamp", "123.000000123")
				return v
			}(),
		},
		{
			name: "duration",
			args: &test.Request{
				Duration: durationpb.New(123000000123),
			},
			want: func() url.Values {
				v := url.Values{}
				v.Set("duration", "123.000000123")
				return v
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (&Encoder{}).Encode(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
