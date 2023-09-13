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
	"strings"
)

func boolToUint32(x bool) uint32 {
	if !x {
		return 0
	} else {
		return 1
	}
}

func replaceStar(s string) string {
	s = strings.Replace(s, ":", "_", -1)
	s = strings.Replace(s, "*", "all", -1)
	return s
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
			s := fmt.Sprintf("\tfmt.Printf(\"%s %%d %%d type=net hostname=%%s\\n\", m.GetTimestamp(), m.Get%s(), hostname)\n", name, name)
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
			s := fmt.Sprintf("\t\tfmt.Printf(\"%s %%d %%d type=nic hostname=%%s name=%%s\\n\", m.GetTimestamp(), iface.Get%s(), hostname, iface.GetName())\n", name, name)
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
				s := "\n\t\tfor _, process := range socket.GetProcesses() {\n"
				s += "\t\t\tfmt.Printf(\"Pid %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s ProcessName=%s\\n\", m.GetTimestamp(), process.GetPid(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()), process.GetName())\n"
				s += "\t\t\tfmt.Printf(\"Fd %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s ProcessName=%s\\n\", m.GetTimestamp(), process.GetFd(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()), process.GetName())\n"
				s += "\t\t}\n\n"
				_, err = goFile.WriteString(s)
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
				continue
			} else if number == 9 {
				s := "\t\tfor _, timer := range socket.GetTimers() {\n"
				s += "\t\t\tfmt.Printf(\"ExpireTimeUs %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s TimerName=%s\\n\", m.GetTimestamp(), timer.GetExpireTimeUs(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()), timer.GetName())\n"
				s += "\t\t\tfmt.Printf(\"Retrans %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s TimerName=%s\\n\", m.GetTimestamp(), timer.GetRetrans(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()), timer.GetName())\n"
				s += "\t\t}\n\n"
				_, err = goFile.WriteString(s)
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
				continue
			} else if number == 10 {
				s := "\t\tfmt.Printf(\"RmemAlloc %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\\n\", m.GetTimestamp(), socket.GetSkmem().GetRmemAlloc(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))\n"
				s += "\t\tfmt.Printf(\"RcvBuf %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\\n\", m.GetTimestamp(), socket.GetSkmem().GetRcvBuf(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))\n"
				s += "\t\tfmt.Printf(\"WmemAlloc %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\\n\", m.GetTimestamp(), socket.GetSkmem().GetWmemAlloc(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))\n"
				s += "\t\tfmt.Printf(\"SndBuf %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\\n\", m.GetTimestamp(), socket.GetSkmem().GetSndBuf(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))\n"
				s += "\t\tfmt.Printf(\"FwdAlloc %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\\n\", m.GetTimestamp(), socket.GetSkmem().GetFwdAlloc(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))\n"
				s += "\t\tfmt.Printf(\"WmemQueued %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\\n\", m.GetTimestamp(), socket.GetSkmem().GetWmemQueued(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))\n"
				s += "\t\tfmt.Printf(\"OptMem %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\\n\", m.GetTimestamp(), socket.GetSkmem().GetOptMem(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))\n"
				s += "\t\tfmt.Printf(\"BackLog %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\\n\", m.GetTimestamp(), socket.GetSkmem().GetBackLog(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))\n"
				s += "\t\tfmt.Printf(\"SockDrop %d %d type=tcp hostname=%s LocalAddr=%s PeerAddr=%s\\n\", m.GetTimestamp(), socket.GetSkmem().GetSockDrop(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))\n"
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
				s = fmt.Sprintf("\t\tfmt.Printf(\"%s %%d %%d type=tcp hostname=%%s LocalAddr=%%s PeerAddr=%%s\\n\", m.GetTimestamp(), socket.Get%s(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))\n", name, name)
			} else if dataType == "double" {
				s = fmt.Sprintf("\t\tfmt.Printf(\"%s %%d %%f type=tcp hostname=%%s LocalAddr=%%s PeerAddr=%%s\\n\", m.GetTimestamp(), socket.Get%s(), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))\n", name, name)
			} else {
				s = fmt.Sprintf("\t\tfmt.Printf(\"%s %%d %%d type=tcp hostname=%%s LocalAddr=%%s PeerAddr=%%s\\n\", m.GetTimestamp(), boolToUint32(socket.Get%s()), hostname, replaceStar(socket.GetLocalAddr()), replaceStar(socket.GetPeerAddr()))\n", name, name)
			}
			_, err = goFile.WriteString(s)
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
		}
	}
}
