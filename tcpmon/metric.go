package tcpmon

import (
	"time"
)

type Metric struct {
	Timestamp time.Time   `bson:"timestamp"`
	Type      string      `bson:"type"`
	Record    interface{} `bson:"record"`
	Raw       string      `bson:"raw"`
}
