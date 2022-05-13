package UI

var (
	fontSize = 15
	fontSize2 = fontSize / 3
)

type node struct {
	id string
	name string
	ip string
	call int
	cost int
	tpc int
	pos position
	width int
	freshTime int64
}

type position struct {
	x int
	y int
}

func (this *node) getWidth() int {
	if this.width == 0 {
		a, b := len(this.name), len(this.ip)
		if a < b {
			this.width = b * fontSize2
		} else {
			this.width = a * fontSize2
		}
	}
	return this.width
}

func (this *node) getHeight() int {
	return fontSize
}