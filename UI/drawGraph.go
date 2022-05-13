package UI

import (
	"../G"
	"bytes"
	"fmt"
)

type drawGraph struct {
	*draw
	textParams bytes.Buffer
	lineParams bytes.Buffer
	lineDetailParams bytes.Buffer
	settingParams string
	allWidth int
}

func createDrawGraph(env string) *drawGraph {
	s := G.SerializeSetting(env)
	g := &drawGraph{draw: createDraw("./UI/graph.html", env)}
	g.js = fmt.Sprintf(g.format, "false", env, "", "", "", "", 0, G.Ip+":"+G.Port)
	g.jsDebug = fmt.Sprintf(g.format, "true", env, "", "", "", s, 0, G.Ip+":"+G.Port)
	g.settingParams = s
	return g
}

func (this *drawGraph) flush() {
	var str1, str2, str3 string
	if this.textParams.Len() > 0 {
		str1 = this.textParams.String()
		str1 = str1[:len(str1)-1]
	}
	if this.lineParams.Len() > 0 {
		str2 = this.lineParams.String()
		str2 = str2[:len(str2)-1]
	}
	if this.lineDetailParams.Len() > 0 {
		str3 = this.lineDetailParams.String()
		str3 = str3[:len(str3)-1]
	}
	this.js = fmt.Sprintf(this.format, "false", this.env, str1, str2, str3, this.settingParams, this.allWidth, G.Ip+":"+G.Port)
	this.jsDebug = fmt.Sprintf(this.format, "true", this.env, str1, str2, str3, this.settingParams, this.allWidth, G.Ip+":"+G.Port)
}

func (this *drawGraph) width(width int) {
	this.allWidth = width
}

func (this *drawGraph) node(id, text string, subText string, pos position) {
	this.textParams.WriteString(fmt.Sprintf("\"%s\",\"%s\",\"%s\",%d,%d,", id, text, subText, pos.x, pos.y))
}

func (this *drawGraph) lineDetail(id, code string, call, cost int) {
	this.lineDetailParams.WriteString(fmt.Sprintf("\"%s\",\"%s\",%d,%d,", id, code, call, cost))
}

func (this *drawGraph) line(id string) {
	this.lineParams.WriteString(fmt.Sprintf("\"%s\",", id))
}

func (this *drawGraph) setting(data string) {
	this.settingParams = data
}

func (this *drawGraph) clear() {
	this.textParams.Reset()
	this.lineParams.Reset()
	this.lineDetailParams.Reset()
}