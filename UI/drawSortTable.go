package UI

import (
	"bytes"
	"fmt"
	"../G"
)

type drawSortTable struct {
	*draw
	title string
	cellParams bytes.Buffer
}

func createDrawSortTable(env string) *drawSortTable {
	g := &drawSortTable{draw:createDraw("./UI/sortTable.html", env)}
	g.js = fmt.Sprintf(g.format, env, G.Ip+":"+G.Port, "")
	//g.jsDebug = g.format
	return g
}

func (this *drawSortTable) insert(ty, serve, code string, call, cost, tpc int) {
	this.cellParams.WriteString(fmt.Sprintf("[\"%s\",\"%s\",\"%s\",%d,%d,%d],", ty, serve, code, call, cost, tpc))
}

func (this *drawSortTable) flush() {
	var str string
	if this.cellParams.Len() > 0 {
		str = this.cellParams.String()
		str = str[:len(str)-1]
	}
	this.js = fmt.Sprintf(this.format, this.env, G.Ip+":"+G.Port, str)
	//this.jsDebug = fmt.Sprintf(this.format, this.cellParams)
}

func (this *drawSortTable) clear() {
	this.cellParams.Reset()
}