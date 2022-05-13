package UI
//
//import (
//	"../Module"
//)
//
//type Tree struct {
//	Root *TreeNode
//	pos position
//	draw *draw
//}
//
//type TreeNode struct {
//	Node
//	parent *TreeNode
//	children []*TreeNode
//}
//
//func CreateTree(name string) *Tree {
//	pos := position{200, 200}
//	t := Tree{
//		Root: &TreeNode{
//			Node:Node{name: name, pos: pos},
//			children:[]*TreeNode{}},
//		pos: pos,
//		draw: createDraw()}
//	return &t
//}
//
//func (this *Tree) Serialize() string {
//	return this.draw.output()
//}
//
//func (this *TreeNode) AddChild(name, ip string) *TreeNode {
//	child := TreeNode{
//		Node:Node{name: name, ip:ip},
//		parent: this,
//		children: []*TreeNode{}}
//	this.children = append(this.children, &child)
//	return &child
//}
//
//func (this *Tree) Flush() {
//	leaf := []*TreeNode{}
//	stack := Module.Stack{}
//	que := Module.Queue{}
//	que.Push(this.Root)
//	for que.Size() > 0 {
//		node := que.Pop().(*TreeNode)
//		if len(node.children) == 0 {
//			leaf = append(leaf, node)
//		} else {
//			stack.Push(node)
//			for _, child := range node.children {
//				que.Push(child)
//			}
//		}
//	}
//
//	que.Clear()
//	base := 0
//	offsetX := map[*TreeNode]int{}
//	for i, cell := range leaf {
//		size_x := cell.width()
//		if i == 0 {
//			offsetX[cell] = base
//		} else {
//			offsetX[cell] = base + size_x / 2
//		}
//		base += size_x + alignX
//	}
//
//	for stack.Size() > 0 {
//		node := stack.Pop().(*TreeNode)
//		offsetX[node] = (offsetX[node.children[0]] + offsetX[node.children[len(node.children)-1]]) / 2
//	}
//
//	que.Clear()
//	que.Push(this.Root)
//	baseX, baseY := this.pos.x, this.pos.y
//	for que.Size() > 0 {
//		node := que.Pop().(*TreeNode)
//		height := 0
//		if node.parent == nil {
//			height = baseY
//		} else {
//			height = node.parent.pos.y + alignY
//		}
//		node.pos.x = baseX + offsetX[node]
//		node.pos.y = height
//		this.draw.text(node.name, node.ip, node.pos)
//		if node.parent != nil {
//			this.draw.line("", node.parent.pos, node.pos, 0)
//		}
//		if len(node.children) > 0 {
//			for _, child := range node.children {
//				que.Push(child)
//			}
//		}
//	}
//}
