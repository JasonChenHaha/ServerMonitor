package main

import (
	"./G"
	"crypto/md5"
	"fmt"
	"net/http"
	"net/url"
)

type control struct {
	cmd *cmd
}

func createControl() *control {
	ctrl := control{cmd: CreateCmd()}
	return &ctrl
}

func (this *control) param_exist_check(param url.Values) error {
	for _, key := range G.RequireParam {
		if _, ok := param[key]; !ok {
			return &G.ServerMonitorErr{G.PARAM_ERR + "no " + key}
		}
	}
	return nil
}

func (this *control) auth_check(param url.Values) error {
	bytes := md5.Sum([]byte(fmt.Sprintf("%s%s%s", param["cmd"][0], param["time"][0], G.SecretKey)))
	str := fmt.Sprintf("%x", bytes)
	if param["auth"][0] != str {
		fmt.Println(str)
		return &G.ServerMonitorErr{G.AUTH_ERR}
	}
	return nil
}

func (this *control) call(req *http.Request) string {
	param := req.Form
	if err := this.param_exist_check(param); err != nil {
		return err.Error()
	}
	if err := this.auth_check(param); err != nil {
		return err.Error()
	}
	if _, ok := this.cmd.cmds[param["cmd"][0]]; !ok {
		return G.PARAM_ERR + "no cmd"
	}
	return this.cmd.cmds[param["cmd"][0]](req)
}