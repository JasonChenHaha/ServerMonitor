package UI

import (
	"io/ioutil"
	"os"
	"unsafe"
)

type draw struct {
	format string
	js string
	jsDebug string
	env string
}

func createDraw(fileName, env string) *draw {
	file, err := os.Open(fileName)
	if err != nil { panic(err) }
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil { panic(err) }
	s1 := *(*string)(unsafe.Pointer(&content))
	return &draw{
		env: env,
		format: s1,
	}
}

func (this *draw) output(isDebug bool) string {
	if isDebug {
		return this.jsDebug
	} else {
		return this.js
	}
}
