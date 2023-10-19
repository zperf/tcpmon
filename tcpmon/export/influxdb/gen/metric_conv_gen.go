package main

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/zperf/tcpmon/tcpmon/gproto"
)

const NetTemplate = `p = write.NewPoint("net",
map[string]string{"Hostname": c.Hostname},
map[string]interface{}{"%s": metric.%s},
ts)
points = append(points, p)
`

func GenNet() {
	t := reflect.TypeOf(gproto.NetstatMetric{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if unicode.IsUpper(rune(field.Name[0])) {
			if field.Name == "Timestamp" || field.Name == "Type" {
				continue
			}
			fmt.Printf(NetTemplate, field.Name, field.Name)
		}
	}
}

const IfaceTemplate = `p = write.NewPoint("nic",
	map[string]string{"Hostname": c.Hostname, "Name": iface.Name},
	map[string]interface{}{"%s": iface.%s},
	ts)
points = append(points, p)
`

func GenNic() {
	t := reflect.TypeOf(gproto.IfaceMetric{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if unicode.IsUpper(rune(field.Name[0])) {
			if field.Name == "Timestamp" || field.Name == "Type" {
				continue
			}
			if field.Name == "Name" {
				continue
			}
			fmt.Printf(IfaceTemplate, field.Name, field.Name)
		}
	}
}

const TcpTemplate = `p = write.NewPoint("tcp", tags,
map[string]interface{}{"%s": s.%s},
ts)
points = append(points, p)
`

func GenTcp() {
	t := reflect.TypeOf(gproto.SocketMetric{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if unicode.IsUpper(rune(field.Name[0])) {
			if field.Name == "Timestamp" || field.Name == "Type" {
				continue
			}
			if field.Name == "Processes" || field.Name == "Timers" {
				continue
			}
			fmt.Printf(TcpTemplate, field.Name, field.Name)
		}
	}
}

func main() {
	//GenNet()
	//GenNic()
	GenTcp()
}
