package Util

type Heap struct {
	data []int64
	kind byte
}

func CreateHeap(k byte) *Heap {
	h := Heap{data: []int64{}, kind: k}
	return &h
}

func (this *Heap) Size() int {
	return len(this.data)
}

func (this *Heap) Push(cell int64) {
	this.data = append(this.data, cell)
	data := this.data
	c := len(this.data) - 1
	p := (c-1) / 2
	if this.kind == '<' {
		for p >= 0 {
			if data[p] > data[c] {
				data[p], data[c] = data[c], data[p]
				c = p;
				p = (c-1) / 2;
			} else {
				break
			}
		}
	} else {
		for p >= 0 {
			if data[p] < data[c] {
				data[p], data[c] = data[c], data[p]
				c = p;
				p = (c-1) / 2;
			} else {
				break
			}
		}
	}
}

func (this *Heap) Top() int64 {
	return this.data[0]
}

func (this *Heap) Pop() (res int64) {
	if len(this.data) == 0 {
		return -1
	}
	last := len(this.data) - 1
	res = this.data[0]
	this.data[0] = this.data[last]
	this.data = this.data[:last]
	if last > 0 {
		heapify(this.data, 0, this.kind)
	}
	return
}

func heapify(data []int64, index int, kind byte) {
	c1 := index * 2 + 1
	c2 := c1 + 1
	index2 := index
	size := len(data)
	if kind == '<' {
		for c1 < size {
			if c2 < size && data[c1] > data[c2] && data[index] > data[c2] {
				index2 = c2
			} else if (data[index] > data[c1]) {
				index2 = c1
			}
			if index != index2 {
				data[index], data[index2] = data[index2], data[index]
				index = index2
				c1 = index * 2 + 1
				c2 = c1 + 1
			} else {
				break
			}
		}
	} else {
		for c1 < size {
			if c2 < size && data[c1] < data[c2] && data[index] < data[c2] {
				index2 = c2
			} else if (data[index] < data[c1]) {
				index2 = c1
			}
			if index != index2 {
				data[index], data[index2] = data[index2], data[index]
				index = index2
				c1 = index * 2 + 1
				c2 = c1 + 1
			} else {
				break
			}
		}
	}
}
