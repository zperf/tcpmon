package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

func main() {
	protoFilePath := "../proto/tcpmon.proto"
	goFilePath := "../tcpmon/print_metric_tsdb.go"

	protoFile, err := os.Open(protoFilePath)
	if err != nil {
		fmt.Printf("Error opening .proto file: %v\n", err)
		return
	}
	defer protoFile.Close()

	goFile, err := os.Create(goFilePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer goFile.Close()

	s := `package tcpmon

import (
	"fmt"
)

func boolToUint32(x bool) uint32 {
	if !x {
		return 0
	} else {
		return 1
	}
}

type TSDBMetricPrinter struct {}

`
	_, err = goFile.WriteString(s)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	// 0: default value
	// 1: in NetstatMetric
	// 2: in IfaceMetric
	// 3: in SocketMetric
	flag := 0

	scanner := bufio.NewScanner(protoFile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if flag == 0 { // default value
			if line == "message NetstatMetric {" {
				flag = 1
				s := "func (tsdb TSDBMetricPrinter) PrintNetstatMetric(m *NetstatMetric, hostname string) {\n"
				_, err = goFile.WriteString(s)
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
			} else if line == "message IfaceMetric {" {
				flag = 2
				s := "func (tsdb TSDBMetricPrinter) PrintNicMetric(m *NicMetric, hostname string) {\n"
				s += "\tfor _, iface := range m.GetIfaces() {\n"
				_, err = goFile.WriteString(s)
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
			} else if line == "message SocketMetric {" {
				flag = 3
				s := "func (tsdb TSDBMetricPrinter) PrintTcpMetric(m *TcpMetric, hostname string) {\n"
				s += "\tfor _, socket := range m.GetSockets() {\n"
				_, err = goFile.WriteString(s)
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
			}
		} else if flag == 1 { // NetstatMetric
			if line == "" || strings.HasPrefix(line, "//") {
				continue
			} else if line == "}" {
				flag = 0
				_, err = goFile.WriteString("}\n")
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
				continue
			}
			line = strings.TrimSuffix(line, ";")
			lineSplit := strings.Split(line, " ")
			if len(lineSplit) != 4 {
				panic("Wrong length of lineSplit")
			}
			name := strcase.ToCamel(lineSplit[1])
			number, _ := strconv.Atoi(lineSplit[3])
			if number <= 2 {
				continue
			}
			s := fmt.Sprintf("\tfmt.Printf(\"net,hostname=%%s %s=%%d %%d\\n\", hostname, m.Get%s(), m.GetTimestamp())\n", name, name)
			_, err = goFile.WriteString(s)
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
		} else if flag == 2 { // IfaceMetric
			if line == "" || strings.HasPrefix(line, "//") {
				continue
			} else if line == "}" {
				flag = 0
				_, err = goFile.WriteString("\t}\n}\n\n")
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
				continue
			}
			line = strings.TrimSuffix(line, ";")
			lineSplit := strings.Split(line, " ")
			if len(lineSplit) != 4 {
				panic("Wrong length of lineSplit")
			}
			name := strcase.ToCamel(lineSplit[1])
			number, _ := strconv.Atoi(lineSplit[3])
			if number <= 1 {
				continue
			}
			s := fmt.Sprintf("\t\tfmt.Printf(\"nic,name=%%s,hostname=%%s %s=%%d %%d\\n\", iface.GetName(), hostname, iface.Get%s(), m.GetTimestamp())\n", name, name)
			_, err = goFile.WriteString(s)
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
		} else if flag == 3 { // SocketMetric
			if line == "" || strings.HasPrefix(line, "//") {
				continue
			} else if line == "}" {
				flag = 0
				_, err = goFile.WriteString("\t}\n}\n\n")
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
				continue
			}
			line = strings.Split(line, "//")[0]
			line = strings.TrimSpace(line)
			line = strings.TrimSuffix(line, ";")
			lineSplit := strings.Split(line, " ")
			if !(len(lineSplit) == 4 || lineSplit[0] == "repeated" && len(lineSplit) == 5) {
				panic("Wrong length of lineSplit")
			}
			var name string
			var number int
			var dataType string
			if len(lineSplit) == 4 {
				name = strcase.ToCamel(lineSplit[1])
				number, _ = strconv.Atoi(lineSplit[3])
				dataType = lineSplit[0]
			} else {
				name = strcase.ToCamel(lineSplit[2])
				number, _ = strconv.Atoi(lineSplit[4])
				dataType = lineSplit[1]
			}
			if number == 6 || number == 7 {
				continue
			} else if number == 8 {
				processMetric := []string{"Pid", "Fd"}
				s := "\n\t\tfor _, process := range socket.GetProcesses() {\n"
				for _, tmp := range processMetric {
					s += fmt.Sprintf("\t\t\tfmt.Printf(\"tcp,LocalAddr=%%s,PeerAddr=%%s,ProcessName=%%s,hostname=%%s %s=%%d %%d\\n\", socket.GetLocalAddr(), socket.GetPeerAddr(), process.GetName(), hostname, process.Get%s(), m.GetTimestamp())\n", tmp, tmp)
				}
				s += "\t\t}\n\n"
				_, err = goFile.WriteString(s)
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
				continue
			} else if number == 9 {
				timerMetric := []string{"ExpireTimeUs", "Retrans"}
				s := "\t\tfor _, timer := range socket.GetTimers() {\n"
				for _, tmp := range timerMetric {
					s += fmt.Sprintf("\t\t\tfmt.Printf(\"tcp,LocalAddr=%%s,PeerAddr=%%s,TimerName=%%s,hostname=%%s %s=%%d %%d\\n\", socket.GetLocalAddr(), socket.GetPeerAddr(), timer.GetName(), hostname, timer.Get%s(), m.GetTimestamp())\n", tmp, tmp)
				}
				s += "\t\t}\n\n"
				_, err = goFile.WriteString(s)
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
				continue
			} else if number == 10 {
				skmemMetric := []string{"RmemAlloc", "RcvBuf", "WmemAlloc", "SndBuf", "FwdAlloc", "WmemQueued", "OptMem", "BackLog", "SockDrop"}
				s := ""
				for _, tmp := range skmemMetric {
					s += fmt.Sprintf("\t\tfmt.Printf(\"tcp,LocalAddr=%%s,PeerAddr=%%s,hostname=%%s %s=%%d %%d\\n\", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.GetSkmem().Get%s(), m.GetTimestamp())\n", tmp, tmp)
				}
				s += "\n"
				_, err = goFile.WriteString(s)
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
				continue
			}
			if dataType != "bool" && dataType != "double" && dataType[0:4] != "uint" && dataType != "SocketState" {
				panic("Wrong data type")
			}
			var s string
			if dataType[0:4] == "uint" || dataType == "SocketState" {
				s = fmt.Sprintf("\t\tfmt.Printf(\"tcp,LocalAddr=%%s,PeerAddr=%%s,hostname=%%s %s=%%d %%d\\n\", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.Get%s(), m.GetTimestamp())\n", name, name)
			} else if dataType == "double" {
				s = fmt.Sprintf("\t\tfmt.Printf(\"tcp,LocalAddr=%%s,PeerAddr=%%s,hostname=%%s %s=%%f %%d\\n\", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, socket.Get%s(), m.GetTimestamp())\n", name, name)
			} else {
				s = fmt.Sprintf("\t\tfmt.Printf(\"tcp,LocalAddr=%%s,PeerAddr=%%s,hostname=%%s %s=%%d %%d\\n\", socket.GetLocalAddr(), socket.GetPeerAddr(), hostname, boolToUint32(socket.Get%s()), m.GetTimestamp())\n", name, name)
			}
			_, err = goFile.WriteString(s)
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
		}
	}
}
