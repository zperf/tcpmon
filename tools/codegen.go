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
	goFilePath := "../tcpmon/parse.go"

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

	s := "package tcpmon\n\n"
	s += "import \"fmt\"\n\n"
	s += "func boolToUint32(x bool) uint32 {\n"
	s += "\tif x == false {\n"
	s += "\t\treturn 0\n"
	s += "\t} else {\n"
	s += "\t\treturn 1\n"
	s += "\t}\n"
	s += "}\n\n"
	_, err = goFile.WriteString(s)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	// 0: 默认值
	// 1: 在NetstatMetric内
	// 2: 在IfaceMetric内
	// 3: 在SocketMetric内
	flag := 0

	// 逐行读取文件内容
	scanner := bufio.NewScanner(protoFile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if flag == 0 { // 默认值
			if line == "message NetstatMetric {" {
				flag = 1
				s := "func (m *NetstatMetric) TSDBPrintNetstatMetric(hostname string) {\n"
				_, err = goFile.WriteString(s)
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
			} else if line == "message IfaceMetric {" {
				flag = 2
				s := "func (m *NicMetric) TSDBPrintNicMetric(hostname string) {\n"
				s += "\tfor _, iface := range m.GetIfaces() {\n"
				_, err = goFile.WriteString(s)
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
			} else if line == "message SocketMetric {" {
				flag = 3
				s := "func (m *TcpMetric) TSDBPrintTcpMetric(hostname string) {\n"
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
			s := fmt.Sprintf("\tfmt.Printf(\"%s type=net,hostname=%%s %%d %%d\\n\", hostname, m.GetTimestamp(), m.Get%s())\n", name, name)
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
			s := fmt.Sprintf("\t\tfmt.Printf(\"%s type=nic,hostname=%%s,name=%%s %%d %%d\\n\", hostname, iface.GetName(), m.GetTimestamp(), iface.Get%s())\n", name, name)
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
				s += "\t\t\tfmt.Printf(\"Pid type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s,ProcessName=%s %d %d\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), process.GetName(), m.GetTimestamp(), process.GetPid())\n"
				s += "\t\t\tfmt.Printf(\"Fd type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s,ProcessName=%s %d %d\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), process.GetName(), m.GetTimestamp(), process.GetFd())\n"
				s += "\t\t}\n\n"
				_, err = goFile.WriteString(s)
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
				continue
			} else if number == 9 {
				s := "\t\tfor _, timer := range socket.GetTimers() {\n"
				s += "\t\t\tfmt.Printf(\"ExpireTimeUs type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s,TimerName=%s %d %d\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), timer.GetName(), m.GetTimestamp(), timer.GetExpireTimeUs())\n"
				s += "\t\t\tfmt.Printf(\"Retrans type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s,TimerName=%s %d %d\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), timer.GetName(), m.GetTimestamp(), timer.GetRetrans())\n"
				s += "\t\t}\n\n"
				_, err = goFile.WriteString(s)
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
					return
				}
				continue
			} else if number == 10 {
				s := "\t\tfmt.Printf(\"RmemAlloc type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetRmemAlloc())\n"
				s += "\t\tfmt.Printf(\"RcvBuf type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetRcvBuf())\n"
				s += "\t\tfmt.Printf(\"WmemAlloc type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetWmemAlloc())\n"
				s += "\t\tfmt.Printf(\"SndBuf type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetSndBuf())\n"
				s += "\t\tfmt.Printf(\"FwdAlloc type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetFwdAlloc())\n"
				s += "\t\tfmt.Printf(\"WmemQueued type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetWmemQueued())\n"
				s += "\t\tfmt.Printf(\"OptMem type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetOptMem())\n"
				s += "\t\tfmt.Printf(\"BackLog type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetBackLog())\n"
				s += "\t\tfmt.Printf(\"SockDrop type=tcp,hostname=%s,LocalAddr=%s,PeerAddr=%s %d %d\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.GetSkmem().GetSockDrop())\n"
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
				s = fmt.Sprintf("\t\tfmt.Printf(\"%s type=tcp,hostname=%%s,LocalAddr=%%s,PeerAddr=%%s %%d %%d\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.Get%s())\n", name, name)
			} else if dataType == "double" {
				s = fmt.Sprintf("\t\tfmt.Printf(\"%s type=tcp,hostname=%%s,LocalAddr=%%s,PeerAddr=%%s %%d %%f\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), socket.Get%s())\n", name, name)
			} else {
				s = fmt.Sprintf("\t\tfmt.Printf(\"%s type=tcp,hostname=%%s,LocalAddr=%%s,PeerAddr=%%s %%d %%d\\n\", hostname, socket.GetLocalAddr(), socket.GetPeerAddr(), m.GetTimestamp(), boolToUint32(socket.Get%s()))\n", name, name)
			}
			_, err = goFile.WriteString(s)
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
		}
	}
}
