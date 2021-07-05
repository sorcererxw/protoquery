package protoquery

import (
	"net/url"
	"testing"
	"time"

	test "github.com/sorcererxw/protoquery/testdata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDecoder_Decode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    url.Values
		want    proto.Message
		wantErr bool
	}{
		{
			name: "",
			args: func() url.Values {
				v := url.Values{}
				return v
			}(),
			want: &test.Request{},
		},
		{
			name: "int64",
			args: func() url.Values {
				v := url.Values{}
				v.Set("int64", "64")
				return v
			}(),
			want: &test.Request{
				Int64: 64,
			},
		},
		{
			name: "enum",
			args: func() url.Values {
				v := url.Values{}
				v.Set("enum", "ENUM")
				return v
			}(),
			want: &test.Request{},
		},
		{
			name: "enum",
			args: func() url.Values {
				v := url.Values{}
				v.Set("enum", "ENUM1")
				return v
			}(),
			want: &test.Request{
				Enum: test.Request_ENUM1,
			},
		},
		{
			name: "repeated",
			args: func() url.Values {
				v := url.Values{}
				v.Set("repeatedString", "12,34")
				return v
			}(),
			want: &test.Request{
				RepeatedString: []string{"12", "34"},
			},
		},
		{
			name: "timestamp",
			args: func() url.Values {
				v := url.Values{}
				v.Set("timestamp", "12.678")
				return v
			}(),
			want: &test.Request{
				Timestamp: timestamppb.New(time.Unix(12, 678000000)),
			},
		},
		{
			name: "timestamp",
			args: func() url.Values {
				v := url.Values{}
				v.Set("timestamp", "-12.678")
				return v
			}(),
			want: &test.Request{
				Timestamp: timestamppb.New(time.Unix(-12, 678000000)),
			},
		},
		{
			name: "duration",
			args: func() url.Values {
				v := url.Values{}
				v.Set("duration", "12.678")
				return v
			}(),
			want: &test.Request{
				Duration: durationpb.New(12678000000),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Decoder{}
			dest := &test.Request{}
			if err := d.Decode(tt.args, dest); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !proto.Equal(tt.want, dest) {
				t.Errorf("Decode() = %v, want = %v", dest, tt.want)
			}
		})
	}
}

func TestDecoder_decodeUnix(t *testing.T) {
	tests := []struct {
		name     string
		args     string
		wantSec  int64
		wantNsec int64
		wantErr  bool
	}{
		{
			name:    "",
			args:    "123",
			wantSec: 123,
		},
		{
			name:     "",
			args:     "123.1",
			wantSec:  123,
			wantNsec: 100000000,
		},
		{
			name:     "",
			args:     ".1",
			wantNsec: 100000000,
		},
		{
			name:    "",
			args:    "1.",
			wantSec: 1,
		},
		{
			name:    "",
			args:    "-1.",
			wantSec: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Decoder{}
			gotSec, gotNsec, err := d.decodeUnix(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeUnix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSec != tt.wantSec {
				t.Errorf("decodeUnix() gotSec = %v, want %v", gotSec, tt.wantSec)
			}
			if gotNsec != tt.wantNsec {
				t.Errorf("decodeUnix() gotNsec = %v, want %v", gotNsec, tt.wantNsec)
			}
		})
	}
}
