package UI

import (
	"../G"
	"strings"
)

type SortTable struct {
	env string
	bodys map[string][]*G.ReportBody // fromId
	needFlush bool
	draw *drawSortTable
}

func CreateSortTable(env string) *SortTable {
	return &SortTable{
		env: env,
		bodys: map[string][]*G.ReportBody{},
		draw: createDrawSortTable(env),
	}
}

func (this *SortTable) Clear(fromId string) {
	this.bodys[fromId] = nil
}

func (this *SortTable) Insert(body *G.ReportBody) {
	if _, ok := this.bodys[body.FromId]; !ok {
		this.bodys[body.FromId] = []*G.ReportBody{}
	}
	this.bodys[body.FromId] = append(this.bodys[body.FromId], body)
	this.needFlush = true
}

func (this *SortTable) Output(isDebug bool) string {
	if this.needFlush {
		this.needFlush = false
		this.flush()
	}
	return this.draw.output(isDebug)
}

func (this *SortTable) flush() {
	this.draw.clear()

	record := map[string]*G.ReportBody{}
	for _, bodys := range this.bodys {
		for _, body := range bodys {
			if _, ok := record[body.Code]; !ok {
				record[body.Code] = &G.ReportBody{
					To: body.To,
					Call: body.Call,
					Code: body.Code,
					Cost: body.Cost,
					Tpc: body.Tpc,
				}
			} else {
				record[body.Code].Call += body.Call
				record[body.Code].Cost += body.Cost
			}
		}
	}
	for _, body := range record {
		body.Tpc = int(body.Cost / body.Call)
	}
	for _, body := range record {
		ty := "act"
		if strings.Contains(body.To, "cmd") {
			ty = "cmd"
		}
		this.draw.insert(ty, body.To, body.Code, int(body.Call), int(body.Cost), int(body.Tpc))
	}
	this.draw.flush()
}