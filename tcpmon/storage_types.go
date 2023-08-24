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

// PrefixMetadata is the metadata of the database
const PrefixMetadata = "mdt"

const MetadataTypeCount = 4

var KeyTotalCount = PrefixMetadata + "/count/total"
var KeyTcpCount = PrefixMetadata + "/count/" + PrefixTcpMetric
var KeyNicCount = PrefixMetadata + "/count/" + PrefixNicMetric
var KeyNetCount = PrefixMetadata + "/count/" + PrefixNetMetric

func ValidCountKind(s string) bool {
	return s == "total" || s == PrefixTcpMetric || s == PrefixNicMetric || s == PrefixNetMetric
}

func ValidMetricPrefix(s string) bool {
	return s == PrefixNicMetric || s == PrefixTcpMetric || s == PrefixNetMetric
}

func ValidPrefix(s string) bool {
	return ValidMetricPrefix(s) || s == PrefixMember
}

func GetPrefixType(key string) string {
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
	Callback func()
}

func NewKVPair(key string, value []byte) *KVPair {
	return &KVPair{
		Key:   key,
		Value: value,
	}
}

func (p KVPair) IsSignal() bool {
	return strings.HasPrefix(p.Key, PrefixSignal)
}

func (p KVPair) ToProto() (proto.Message, error) {
	kind := GetPrefixType(p.Key)
	if kind == "" {
		return nil, errors.Newf("invalid kind: '%v'", kind)
	}

	switch kind {
	case PrefixTcpMetric:
		var m TcpMetric
		err := proto.Unmarshal(p.Value, &m)
		if err != nil {
			return nil, err
		}
		return &m, nil

	case PrefixNetMetric:
		var m NetstatMetric
		err := proto.Unmarshal(p.Value, &m)
		if err != nil {
			return nil, err
		}
		return &m, nil

	case PrefixNicMetric:
		var m NicMetric
		err := proto.Unmarshal(p.Value, &m)
		if err != nil {
			return nil, err
		}
		return &m, nil

	default:
		return nil, errors.Newf("Should not reach here")
	}
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
