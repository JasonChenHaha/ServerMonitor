package Util

import (
	"../G"
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var Alarm = CreateAlarm()
var notis = map[string][]string {
	"0":[]string{"pandongfang@mini1.cn", "chenjian@mini1.cn", "chenjirui@mini1.cn"},
	"1":[]string{"pandongfang@mini1.cn", "chenjian@mini1.cn"},
	"10":[]string{"wangyongfu@mini1.cn", "longjianxin@mini1.cn", "pandongfang@mini1.cn", "lihanfeng@mini1.cn"},
}

type alarm struct {
	nodeRecord map[string]map[string]int64		// env->id->freshTime
	nodeReportCounter map[string]map[string]int		// env->id->int
	tpcRecord map[string]map[string]map[string]map[string][]*G.ReportBody	// env->from->to->code->[]body
}

func CreateAlarm() *alarm {
	a := &alarm{
		nodeRecord: map[string]map[string]int64{},
		nodeReportCounter: map[string]map[string]int{},
		tpcRecord: map[string]map[string]map[string]map[string][]*G.ReportBody{},
	}
	for env, _ := range G.Env {
		a.nodeRecord[env] = map[string]int64{}
		a.nodeReportCounter[env] = map[string]int{}
		a.tpcRecord[env] = map[string]map[string]map[string][]*G.ReportBody{}
	}
	if !G.Debug {
		go a.tick()
	}
	return a
}

func (this *alarm) Report_node(env, id string) {
	this.nodeRecord[env][id] = time.Now().Unix()
	if _, ok := this.nodeReportCounter[env][id]; !ok {
		this.nodeReportCounter[env][id] = 0
	}
}

func (this *alarm) Check_tpc(env string, body *G.ReportBody) {
	a := this.tpcRecord[env]
	if _, ok := a[body.FromId]; !ok { a[body.FromId] = map[string]map[string][]*G.ReportBody{} }
	b := a[body.FromId]
	if _, ok := b[body.ToId]; !ok { b[body.ToId] = map[string][]*G.ReportBody{} }
	c := b[body.ToId]
	if _, ok := c[body.Code]; !ok { c[body.Code] = []*G.ReportBody{} }

	if !G.Debug && (env == "0" || env == "10") {
		alarmTop := G.DefaultAlarmTop[env]
		if _, ok := G.AlarmTop[env][body.Code]; ok {
			alarmTop = G.AlarmTop[env][body.Code]
		} else if _, ok := G.AlarmTop[env][body.ToId]; ok {
			alarmTop = G.AlarmTop[env][body.ToId]
		}
		if body.Tpc >= alarmTop {
			c[body.Code] = append(c[body.Code], body)
			alarmTimes := G.DefaultAlarmTimes[env]
			if t, ok := G.AlarmTimes[env][body.Code]; ok {
				alarmTimes = t
			} else if t, ok := G.AlarmTimes[env][body.ToId]; ok {
				alarmTimes = t
			}
			if len(c[body.Code]) >= alarmTimes {
				this.sound_tpc_alarm(env, c[body.Code])
				c[body.Code] = c[body.Code][0:0]
			}
		} else {
			c[body.Code] = c[body.Code][0:0]
		}
	}
}

// 当节点正常上报超过30分钟之后，发现为常规节点
// 如果常规节点在一个tick时间内没有收到心跳则会告警
// 在若干次告警之后，系统将其视为非常规节点，并不再告警
// 直到正常上报若干次后，再次将其发现为常规节点
func (this *alarm) tick() {
	c := time.Tick(time.Second * 301)
	for {
		now := (<- c).Unix()
		for env, m := range this.nodeRecord {
			if env == "0" || env == "10" {
				for id, freshTime := range m {
					if freshTime + G.NodeAliveTime < now {
						if this.nodeReportCounter[env][id] > 0 {
							this.nodeReportCounter[env][id]--
							if this.nodeReportCounter[env][id] >= 6 {
								this.sound_node_alarm(env, id)
							}
						}
					} else {
						if this.nodeReportCounter[env][id] < 7 {
							this.nodeReportCounter[env][id]++
						}
					}
				}
			}
		}
	}
}

func (this *alarm) sound_node_alarm(env string, id string) {
	now := time.Now().Unix()
	sign := fmt.Sprintf("%x", md5.Sum([]byte(strconv.FormatInt(now, 10) + "miniw_0112f(123dsKsuqbfY")))
	msg := url.QueryEscape(fmt.Sprintf("Reason: 节点心跳丢失\nName:%s\nEnv:%s\n", id, env))
	for _, mail := range notis[env] {
		rsp, err := http.Get(fmt.Sprintf("http://120.24.64.132:8080/miniw/feishu?act=send_normal&title=ServerMonitor告警&msg=%s&role=ServerMonitor&time=%d&sign=%s&email=%s&msg_type=text",
			msg, now, sign, mail))
		if err != nil {
			panic(err)
		}
		defer rsp.Body.Close()
	}
}

func (this *alarm) sound_tpc_alarm(env string, bodys []*G.ReportBody) {
	now := time.Now().Unix()
	sign := fmt.Sprintf("%x", md5.Sum([]byte(strconv.FormatInt(now, 10) + "miniw_0112f(123dsKsuqbfY")))
	buff := bytes.Buffer{}
	buff.WriteString(fmt.Sprintf("Reason: tpc超过警戒值\nFrom: %s\nTo: %s\nCode: %s\nTpc: [", bodys[0].FromId, bodys[0].ToId, bodys[0].Code))
	for i, body := range bodys {
		if i+1 == len(bodys) {
			buff.WriteString(fmt.Sprintf("%d]", body.Tpc))
		} else {
			buff.WriteString(fmt.Sprintf("%d,", body.Tpc))
		}
	}
	alarmTop := G.DefaultAlarmTop[env]
	if a, ok := G.AlarmTop[env][bodys[0].ToId]; ok {
		alarmTop = a
	} else if a, ok := G.AlarmTop[env][bodys[0].Code]; ok {
		alarmTop = a
	}
	interval := G.DefaultAlarmInterval[env]
	if a, ok := G.AlarmInterval[env][bodys[0].ToId]; ok {
		interval = a
	} else if a, ok := G.AlarmInterval[env][bodys[0].Code]; ok {
		interval = a
	}
	buff.WriteString(fmt.Sprintf("\nEnv: %s\nTpcTop: %d\nInterval: %dmin", env, alarmTop, interval))
	msg := url.QueryEscape(buff.String())
	for _, mail := range notis[env] {
		rsp, err := http.Get(fmt.Sprintf("http://120.24.64.132:8080/miniw/feishu?act=send_normal&title=ServerMonitor告警&msg=%s&role=ServerMonitor&time=%d&sign=%s&email=%s&msg_type=text",
			msg, now, sign, mail))
		if err != nil {
			panic(err)
		}
		defer rsp.Body.Close()
	}
}