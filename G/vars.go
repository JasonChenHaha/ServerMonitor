package G

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

var Debug = false
var Ip = ""
var Port = ""

var SecretKey = "8ae5fc3f70d4edee5957695618c21cd8"

var Env = map[string]map[string]string{
	"0": {},
	"1": {},
	"2": {},
	"10": {},
	"11": {},
	"12": {},
}

var RequireParam = []string{
	"cmd",
	"env",
	"time",
	"auth",
}

var OK = "ok"
var SYS_ERR = "system error: "
var PARAM_ERR = "param error: "
var AUTH_ERR = "auth error"

var DefaultLevel = 1
var LevelMap = map[string]int {
	"gate": 0,
}

var NodeAliveTime int64 = 60

var DefaultTpcTop = map[string]int {"0": 200, "1": 200, "2": 200, "10": 200, "11": 200, "12": 200}
var TpcTop = map[string]map[string]int {"0": {}, "1": {}, "2": {}, "10": {}, "11": {}, "12": {}}

var DefaultAlarmTop = map[string]int {"0": 5000, "1": 5000, "2": 5000, "10": 5000, "11": 5000, "12": 5000}
var AlarmTop = map[string]map[string]int {"0": {}, "1": {}, "2": {}, "10": {}, "11": {}, "12": {}}

var DefaultAlarmInterval = map[string]int {"0": 2, "1": 2, "2": 2, "10": 2, "11": 2, "12": 2}
var AlarmInterval = map[string]map[string]int {"0": {}, "1": {}, "2": {}, "10": {}, "11": {}, "12": {}}

var DefaultAlarmTimes = map[string]int {}
var AlarmTimes = map[string]map[string]int {"0": {}, "1": {}, "2": {}, "10": {}, "11": {}, "12": {}}

type ReportBody struct {
	From string
	FromId string
	FromIp string
	To string
	ToId string
	ToIp string
	Code string
	Call float64
	Cost float64
	Tpc int
}


type ReportPackage struct {
	Env string
	RemoteAddr string
	Bodys []*ReportBody
}

func init() {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipNet, isIpNet := addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				str := ipNet.IP.String()
				if strings.HasPrefix(str, "10.") ||
					strings.HasPrefix(str, "172.") ||
					strings.HasPrefix(str, "192.") {
				} else {
					Ip = str
					break
				}
			}
		}
	}

	if _, err := os.Stat("../conf"); err != nil {
		return
	}
	f, err := os.OpenFile("../conf/setting.conf", os.O_RDONLY, os.ModePerm|os.ModeTemporary)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if len(line) == 0 { break }
			} else { panic(err) }
		}
		strs := strings.Split(line[:len(line)-1], "=")
		switch strs[0] {
		case "ip":
			Ip = strs[1]
		case "port":
			Port = strs[1]
		case "debug":
			if strs[1] != "0" && strs[1] != "false" {
				Debug = true
			}
		}
	}

	if _, err := os.Stat("../conf"); err != nil {
		return
	}
	f, err = os.OpenFile("../conf/custom.conf", os.O_RDONLY, os.ModePerm|os.ModeTemporary)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	reader = bufio.NewReader(f)
	var env, node, code string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if len(line) == 0 { break }
			} else { panic(err) }
		}
		strs := strings.Split(line[:len(line)-1], "=")
		switch strs[0] {
		case "env":
			env = strs[1]
			node = ""
			code = ""
		case "node":
			node = strs[1]
			code = ""
		case "code":
			node = ""
			code = strs[1]
		case "tcpTop":
			num, err := strconv.Atoi(strs[1])
			if err != nil { panic(err) }
			if code != "" {
				TpcTop[env][code] = num
			} else if node != "" {
				TpcTop[env][node] = num
			} else {
				DefaultTpcTop[env] = num
			}
		case "alarmTop":
			num, err := strconv.Atoi(strs[1])
			if err != nil { panic(err) }
			if code != "" {
				AlarmTop[env][code] = num
			} else if node != "" {
				AlarmTop[env][node] = num
			} else {
				DefaultAlarmTop[env] = num
			}
		case "alarmInterval":
			num, err := strconv.Atoi(strs[1])
			if err != nil { panic(err) }
			if code != "" {
				AlarmInterval[env][code] = num
				AlarmTimes[env][code] = num * 2
			} else if node != "" {
				AlarmInterval[env][node] = num
				AlarmTimes[env][node] = num * 2
			} else {
				DefaultAlarmInterval[env] = num
				DefaultAlarmTimes[env] = num * 2
			}
		}
	}
}