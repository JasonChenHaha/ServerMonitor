package Util

type Stack struct {
	data []interface{}
}

func (this *Stack) Push(cell interface{}) {
	this.data = append(this.data, cell)
}

func (this *Stack) Pop() interface{} {
	cell := this.data[len(this.data)-1]
	this.data = this.data[:len(this.data)-1]
	return cell
}

func (this *Stack) Size() int {
	return len(this.data)
}

func (this *Stack) Clear() {
	this.data = this.data[0:0]
}