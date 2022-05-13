package UI

import (
	"../G"
	"math"
	"time"
)

var graphAlignX = 800
var graphAlignY = [2]int{50, 400}

type Graph struct {
	env string
	nodeByLevelId map[int]map[string]*node // level:id
	nodeByLevelArr map[int][]*node // level
	relation map[string]map[string]map[string]*G.ReportBody	// from:to:code
	relationAll map[string]map[string]*relation // from:to:relation
	draw *drawGraph
	needFlush bool
	FlushTime int64
}

func CreateGraph(env string) *Graph {
	return &Graph{
		env: env,
		nodeByLevelId: map[int]map[string]*node{},
		nodeByLevelArr: map[int][]*node{},
		relation: map[string]map[string]map[string]*G.ReportBody{},
		relationAll: map[string]map[string]*relation{},
		FlushTime: -1,
		draw: createDrawGraph(env),
	}
}

func (this *Graph) Output(isDebug bool) string {
	if this.needFlush {
		this.needFlush = false;
		this.flush()
	}
	if this.FlushTime > 0 && this.FlushTime < time.Now().Unix() {
		this.FlushTime = -1
		this.flush()
	}
	return this.draw.output(isDebug)
}

func (this *Graph) AddNode(id, name, ip string, level, call, cost, tpc int) {
	if _, ok := this.nodeByLevelId[level]; !ok {
		this.nodeByLevelId[level] = map[string]*node{}
		this.nodeByLevelArr[level] = make([]*node, 0, 666)
	}
	if _, ok := this.nodeByLevelId[level][id]; !ok {
		n := &node{id: id, name: name, ip: ip, width: 0}
		this.nodeByLevelArr[level] = append(this.nodeByLevelArr[level], n)
		size := len(this.nodeByLevelArr[level])
		// 插入排序
		for i := 0; i < size; i++ {
			if n.name <= this.nodeByLevelArr[level][i].name {
				copy(this.nodeByLevelArr[level][i+1:size], this.nodeByLevelArr[level][i:size-1])
				this.nodeByLevelArr[level][i] = n
				break
			}
		}
		this.nodeByLevelId[level][id] = n
	}
	n := this.nodeByLevelId[level][id]
	n.call = call
	n.cost = cost
	n.tpc = tpc
	n.freshTime = time.Now().Unix()
}

func (this *Graph) AddRelations(data map[string]*G.ReportBody) {
	var fromId, toId string
	for _, body := range data {
		fromId = body.FromId
		toId = body.ToId
		break
	}

	if _, ok := this.relation[fromId]; !ok {
		this.relation[fromId] = map[string]map[string]*G.ReportBody{}
	}
	this.relation[fromId][toId] = data
	if _, ok := this.relationAll[fromId]; !ok {
		this.relationAll[fromId] = map[string]*relation{}
	}
	if _, ok := this.relationAll[fromId][toId]; !ok {
		this.relationAll[fromId][toId] = &relation{id: fromId + "->" + toId}
	}
	a := this.relationAll[fromId][toId]

	worstTpc, worstRatio, worstTpcTop := 0, 0.0, G.DefaultTpcTop[this.env]
	for _, body := range data {
		tmpTop := G.DefaultTpcTop[this.env]
		if b, ok := G.TpcTop[this.env][body.Code]; ok {
			tmpTop = b
		} else if b, ok := G.TpcTop[this.env][body.ToId]; ok {
			tmpTop = b
		}
		ratio := float64(body.Tpc) / float64(tmpTop)
		if worstRatio < ratio {
			worstRatio = ratio
			worstTpc = body.Tpc
			worstTpcTop = tmpTop
		}
	}
	a.tpc = worstTpc
	a.tpcTop = worstTpcTop
	a.deadTime = time.Now().Unix() + 60
	this.needFlush = true
}

func (this *Graph) Setting(data string) {
	this.draw.setting(data)
	this.needFlush = true
}

// 打开全部注释变成第二种布局方案
func (this *Graph) flush() {
	var maxY, minLv = 0, 999
	var tmpY, size_fy int
	var pos position
	var help = map[int]int{}
	var now = time.Now().Unix()
	var alignY = map[int]int{}
	//var maxNode, minAlignY, tmpAlignY int

	this.draw.clear()
	this.draw.width((len(this.nodeByLevelArr)-1)*graphAlignX)

	//for i, s := range this.nodeByLevelArr {
	//	tmpAlignY = int(math.Max(float64(graphAlignY[0]), float64(graphAlignY[1] - (graphAlignY[1] - graphAlignY[0]) / 10 * len(s))))
	//	if i < minLv { minLv = i }
	//	tmpY = 0
	//	for _, node := range s {
	//		if now <= node.freshTime + G.NodeAliveTime {
	//			if tmpY != 0 { tmpY += tmpAlignY }
	//			tmpY += node.getHeight()
	//		}
	//	}
	//	if tmpY > maxY {
	//		maxY = tmpY
	//		maxNode = len(s)
	//		minAlignY = tmpAlignY
	//	}
	//}

	for lv, s := range this.nodeByLevelArr {
		//alignY[i] = int(math.Min(float64(graphAlignY[1]), float64(minAlignY) * (float64(maxNode) / float64(len(s)))))
		alignY[lv] = int(math.Max(float64(graphAlignY[0]), float64(graphAlignY[1] - (graphAlignY[1] - graphAlignY[0]) / 10 * len(s))))
		if lv < minLv { minLv = lv }
		tmpY = 0
		for _, node := range s {
			if now <= node.freshTime + G.NodeAliveTime {
				if tmpY != 0 { tmpY += alignY[lv] }
				tmpY += node.getHeight()
			}
		}
		help[lv] = tmpY
		if tmpY > maxY {
			maxY = tmpY
		}
	}
	for i, s := range this.nodeByLevelArr {
		pos = position{(i-minLv)*graphAlignX, (maxY-help[i])/2}
		for _, node := range s {
			if now <= node.freshTime + G.NodeAliveTime {
				size_fy = node.getHeight() / 2
				pos.y += size_fy
				node.pos = pos
				pos.y += size_fy + alignY[i]
				this.draw.node(node.id, node.name, node.ip, node.pos)
			}
		}
	}

	for fromId, m := range this.relation {
		for toId, m2 := range m {
			tmp2 := this.relationAll[fromId][toId]
			if tmp2.deadTime < now {
				continue
			}
			for _, body := range m2 {
				this.draw.lineDetail(tmp2.id, body.Code, int(body.Call), int(body.Cost))
			}
			this.draw.line(tmp2.id)
		}
	}
	this.draw.flush()
}