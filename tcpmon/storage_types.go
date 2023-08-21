package tcpmon

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const MetricTypeCount = 3
const PrefixTcpMetric = "tcp"
const PrefixNicMetric = "nic"
const PrefixNetMetric = "net"

// PrefixSignal inject a signal into storage. Only for testing
const PrefixSignal = "sig"

// PrefixMember is the member in the cluster
const PrefixMember = "mbr"

func ValidMetricPrefix(s string) bool {
	return s == PrefixNicMetric || s == PrefixTcpMetric || s == PrefixNetMetric
}

func GetType(key string) string {
	p := strings.IndexRune(key, '/')
	if p < 0 {
		return ""
	}
	return key[:p]
}

func KeyJoin(elems ...string) string {
	return strings.Join(elems, "/")
}

func NewKey(kind string) string {
	return fmt.Sprintf("%s/%v/", kind, time.Now().UnixMilli())
}

type KVPair struct {
	Key      string
	Value    []byte
	Callback func(err error)
}

func (p KVPair) ToProto() (proto.Message, error) {
	kind := GetType(p.Key)
	if kind == "" {
		return nil, errors.Newf("invalid kind: '%v'", kind)
	}

	var msg proto.Message
	switch kind {
	case PrefixTcpMetric:
		var m TcpMetric
		err := proto.Unmarshal(p.Value, &m)
		if err != nil {
			return nil, err
		}
		msg = &m

	case PrefixNetMetric:
		var m NetstatMetric
		err := proto.Unmarshal(p.Value, &m)
		if err != nil {
			return nil, err
		}
		msg = &m

	case PrefixNicMetric:
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

func (p KVPair) ToJSONString() string {
	h := p.ToJSON()
	s, err := json.Marshal(h)
	if err != nil {
		return fmt.Sprintf(`{"error":"%s"}`, err.Error())
	}

	return string(s)
}
