package main

import (
	"./G"
	"./UI"
	"./Util"
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type funcMap map[string]func(interface{})string

type cmd struct {
	cmds  funcMap
	graphs map[string]*UI.Graph
	sortTables map[string]*UI.SortTable
	nameFilter map[string]string
	reportChan chan G.ReportPackage
}

func CreateCmd() *cmd {
	c := cmd{
		graphs: map[string]*UI.Graph{},
		sortTables: map[string]*UI.SortTable{},
		nameFilter: map[string]string{},
		reportChan: make(chan G.ReportPackage, 1000),
	}
	for env, _ := range G.Env {
		c.graphs[env] = UI.CreateGraph(env)
		c.sortTables[env] = UI.CreateSortTable(env)
	}
	c.cmds = funcMap{
		"report": c.Report,
		"graph": c.Graph,
		"sortTable": c.SortTable,
		"setting": c.Setting,
	}
	go report_core(c)
	return &c
}

func report_core(c cmd) {
	for {
		cell := <-c.reportChan
		// 非法上报过滤
		for i := len(cell.Bodys)-1; i >= 0; i-- {
			a, b := len(cell.Bodys[i].From), len(cell.Bodys[i].To)
			if a == 0 || a >= 64 || b == 0 || b >= 64 {
				cell.Bodys = append(cell.Bodys[:i], cell.Bodys[i+1:]...)
			} else {
				cell.Bodys[i].From = strings.Replace(cell.Bodys[i].From, "\"", "", -1)
				cell.Bodys[i].To = strings.Replace(cell.Bodys[i].To, "\"", "", -1)
			}
		}
		// 遍历上报bodys，生成节点和关系
		var toNameIp []string
		var allCall, allCost, allTpc = 0, 0, 0
		for _, body := range cell.Bodys {
			body.FromIp = strings.Split(cell.RemoteAddr, ":")[0]
			body.FromId = fmt.Sprintf("%s@%s", body.From, body.FromIp)
			body.ToId = body.To
			toNameIp = strings.Split(body.To, "@")
			body.To = toNameIp[0]
			body.ToIp = toNameIp[1]
			body.Tpc = int(body.Cost) / int(body.Call)
			allCall += int(body.Call)
			allCost += int(body.Cost)
			c.nameFilter[body.FromId] = "node"
			c.nameFilter[body.ToId] = "node"
			c.nameFilter[body.Code] = "code"
		}
		if allCall > 0 {
			allTpc = allCost / allCall
		}
		graph := c.graphs[cell.Env]
		sortTable := c.sortTables[cell.Env]
		// 清理sortTable
		for _, body := range cell.Bodys {
			sortTable.Clear(body.FromId)
		}
		// 生成节点
		record := map[string]map[string]map[string]*G.ReportBody{}	// fromId -> toId -> code -> reportBody
		level1, level2 := G.DefaultLevel, G.DefaultLevel
		for _, body := range cell.Bodys {
			for key, level := range G.LevelMap {
				if strings.Contains(body.From, key) {
					level1 = level
				} else if strings.Contains(body.To, key) {
					level2 = level
				}
			}
			graph.AddNode(body.FromId, body.From, body.FromIp, level1, allCall, allCost, allTpc)
			graph.AddNode(body.ToId, body.To, body.ToIp, level2, 0, 0, 0)
			sortTable.Insert(body)
			if _, ok := record[body.FromId]; !ok { record[body.FromId] = map[string]map[string]*G.ReportBody{} }
			a := record[body.FromId]
			if _, ok := a[body.ToId]; !ok { a[body.ToId] = map[string]*G.ReportBody{} }
			b := a[body.ToId]
			b[body.Code] = body
			Util.Alarm.Report_node(cell.Env, body.FromId)
			Util.Alarm.Check_tpc(cell.Env, body)
		}
		// 生成连线
		for _, m1 := range record {
			for _, m2 := range m1 {
				graph.AddRelations(m2)		// m2: map[code]reportBody
			}
		}
		graph.FlushTime = time.Now().Unix() + G.NodeAliveTime + 1
	}
}

func (this *cmd) Report(data interface{}) string {
	req := data.(*http.Request)
	var bodys []*G.ReportBody
	err := json.NewDecoder(req.Body).Decode(&bodys)
	if err != nil {
		fmt.Println(G.PARAM_ERR + "decode body failed.")
	}
	cell := G.ReportPackage{Env: req.Form["env"][0], RemoteAddr: req.RemoteAddr, Bodys: bodys}
	this.reportChan <- cell
	return G.OK
}

func (this *cmd) Graph(data interface{}) string {
	req := data.(*http.Request)
	_, isDebug := req.Form["debug"]
	return this.graphs[req.Form["env"][0]].Output(isDebug)
}

func (this *cmd) SortTable(data interface{}) string {
	req := data.(*http.Request)
	_, isDebug := req.Form["debug"]
	return this.sortTables[req.Form["env"][0]].Output(isDebug)
}

func (this *cmd) Setting(data interface{}) string {
	req := data.(*http.Request)
	env := req.Form["env"][0]
	if len(req.Form["name"][0]) > 0 {
		name := req.Form["name"][0]
		a, err := strconv.Atoi(req.Form["tpcTop"][0])
		if err != nil { panic(err) }
		G.TpcTop[env][name] = a
		a, err = strconv.Atoi(req.Form["alarmTop"][0])
		if err != nil { panic(err) }
		G.AlarmTop[env][name] = a
		a, err = strconv.Atoi(req.Form["alarmInterval"][0])
		if err != nil { panic(err) }
		G.AlarmInterval[env][name] = a
		G.AlarmTimes[env][name] = a * 2
	} else {
		a, err := strconv.Atoi(req.Form["tpcTop"][0])
		if err != nil { panic(err) }
		G.DefaultTpcTop[env] = a
		a, err = strconv.Atoi(req.Form["alarmTop"][0])
		if err != nil { panic(err) }
		G.DefaultAlarmTop[env] = a
		a, err = strconv.Atoi(req.Form["alarmInterval"][0])
		if err != nil { panic(err) }
		G.DefaultAlarmInterval[env] = a
		G.DefaultAlarmTimes[env] = a * 2
	}

	this.graphs[env].Setting(G.SerializeSetting(env))

	tmp := map[string]map[string]map[string]string{} // env -> name -> code
	for env, tcpTop := range G.DefaultTpcTop {
		tmp[env] = map[string]map[string]string{}
		tmp[env]["all"] = map[string]string{}
		tmp[env]["all"]["tcpTop"] = strconv.Itoa(tcpTop)
	}
	for env, alarmTop := range G.DefaultAlarmTop {
		tmp[env]["all"]["alarmTop"] = strconv.Itoa(alarmTop)
	}
	for env, alarmInterval := range G.DefaultAlarmInterval {
		tmp[env]["all"]["alarmInterval"] = strconv.Itoa(alarmInterval)
	}
	for env, m := range G.TpcTop {
		for name, value := range m {
			if _, ok := tmp[env][name]; !ok { tmp[env][name] = map[string]string{} }
			tmp[env][name]["tcpTop"] = strconv.Itoa(value)
		}
	}
	for env, m := range G.AlarmTop {
		for name, value := range m {
			if _, ok := tmp[env][name]; !ok { tmp[env][name] = map[string]string{} }
			tmp[env][name]["alarmTop"] = strconv.Itoa(value)
		}
	}
	for env, m := range G.AlarmInterval {
		for name, value := range m {
			if _, ok := tmp[env][name]; !ok { tmp[env][name] = map[string]string{} }
			tmp[env][name]["alarmInterval"] = strconv.Itoa(value)
		}
	}

	if _, err := os.Stat("../conf"); err != nil {
		return G.SYS_ERR + "no conf folder"
	}
	f, err := os.OpenFile("../conf/custom.conf", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm|os.ModeTemporary)
	if err != nil { panic(err) }
	defer f.Close()
	writer := bufio.NewWriter(f)
	for env, m1 := range tmp {
		writer.WriteString("env=" + env + "\n")
		for key, value := range m1["all"] {
			writer.WriteString(key + "=" + value + "\n")
		}
		delete(m1, "all")
		for name, m2 := range m1 {
			writer.WriteString(this.nameFilter[name] + "=" + name + "\n")
			for key, value := range m2 {
				writer.WriteString(key + "=" + value + "\n")
			}
		}
	}
	writer.Flush()

	return G.OK
}