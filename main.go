package main

import (
	"fmt"
	"net/http"
	"./G"
)

var ctrl *control

func main() {
	ctrl = createControl()
	http.HandleFunc("/miniw/goserver", handler)
	http.ListenAndServe("0.0.0.0:"+G.Port, nil)
}

func handler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Fprint(w, ctrl.call(req))
}