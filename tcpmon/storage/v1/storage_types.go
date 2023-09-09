package v1

import (
	"fmt"
	"strings"
	"time"
)

const MetricTypeCount = 3
const PrefixTcpMetric = "tcp"
const PrefixNicMetric = "nic"
const PrefixNetMetric = "net"

// PrefixSignal inject a signal into storage. Only for testing
const PrefixSignal = "sig"

// PrefixMetadata is the metadata of the database
const PrefixMetadata = "mdt"

const MetadataTypeCount = MetricTypeCount + 1

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

func GetPrefixType(key string) string {
	p := strings.IndexRune(key, '/')
	if p < 0 {
		return ""
	}
	return key[:p]
}

func NewKey(kind string) string {
	return fmt.Sprintf("%s/%v/", kind, time.Now().UnixMilli())
}

type MetricContext struct {
	Key      string
	Value    []byte
	Callback func()
}

func NewMetricContext(key string, value []byte) *MetricContext {
	return &MetricContext{
		Key:   key,
		Value: value,
	}
}

func (p MetricContext) IsSignal() bool {
	return strings.HasPrefix(p.Key, PrefixSignal)
}
