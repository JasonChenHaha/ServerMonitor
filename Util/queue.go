package Util

type Queue struct {
	data []interface{}
}

func (this *Queue) Push(cell interface{}) {
	this.data = append(this.data, cell)
}

func (this *Queue) Pop() interface{} {
	cell := this.data[0]
	this.data = this.data[1:len(this.data)]
	return cell
}

func (this *Queue) Size() int {
	return len(this.data)
}

func (this *Queue) Clear() {
	this.data = this.data[0:0]
}
