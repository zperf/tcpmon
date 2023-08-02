package tcpmon

import (
	"encoding/json"
	"strings"

	"github.com/cockroachdb/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const PrefixTcpRecord = "tcp"
const PrefixNicRecord = "nic"
const PrefixNetRecord = "net"

func ValidPrefix(s string) bool {
	return s == PrefixNicRecord || s == PrefixTcpRecord || s == PrefixNetRecord
}

func GetType(key string) string {
	p := strings.IndexRune(key, '/')
	if p < 0 {
		return ""
	}
	return key[:p]
}

type KVPair struct {
	Key   string
	Value []byte
}

func (p KVPair) ToProto() (proto.Message, error) {
	kind := GetType(p.Key)
	if kind == "" {
		return nil, errors.Newf("invalid kind: '%v'", kind)
	}

	var msg proto.Message
	switch kind {
	case PrefixTcpRecord:
		var m TcpMetric
		err := proto.Unmarshal(p.Value, &m)
		if err != nil {
			return nil, err
		}
		msg = &m

	case PrefixNetRecord:
		var m NetstatMetric
		err := proto.Unmarshal(p.Value, &m)
		if err != nil {
			return nil, err
		}
		msg = &m

	case PrefixNicRecord:
		var m NicMetric
		err := proto.Unmarshal(p.Value, &m)
		if err != nil {
			return nil, err
		}
		msg = &m
	}

	return msg, nil
}

func (p KVPair) ToJSON() map[string]any {
	if p.Value == nil {
		return nil
	}

	m, err := p.ToProto()
	if err != nil {
		return ErrorJSON(errors.WithStack(err))
	}

	buf, err := protojson.Marshal(m)
	if err != nil {
		return ErrorJSON(errors.WithStack(err))
	}

	var h map[string]any
	err = json.Unmarshal(buf, &h)
	if err != nil {
		return ErrorJSON(errors.WithStack(err))
	}

	return h
}
