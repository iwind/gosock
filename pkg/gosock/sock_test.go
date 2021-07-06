// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package gosock

import (
	"encoding/json"
	"net"
	"testing"
)

func TestSock_Listen(t *testing.T) {
	var sock = NewSock("test.sock")
	sock.OnCommand(func(cmd *Command) {
		cmdJSON, err := json.Marshal(cmd)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(cmdJSON))
		err = cmd.ReplyOk()
		if err != nil {
			t.Fatal(err)
		}
	})
	sock.On("reload", func(params map[string]interface{}) {
		t.Log("onReload:", params)
	})
	err := sock.Listen()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("end")
}

func TestSock_Client(t *testing.T) {
	var sock = NewSock("test.sock")
	conn, err := net.Dial("unix", sock.path)
	if err != nil {
		t.Fatal(err)
	}

	cmd := &Command{
		Code: "reload",
		Params: map[string]interface{}{
			"a":       1,
			"b":       false,
			"newLine": "a\nb\nc",
		},
	}
	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		t.Fatal(err)
	}

	cmdJSON = append(cmdJSON, '\n')
	_, err = conn.Write(cmdJSON)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}

func TestSock_Send(t *testing.T) {
	var sock = NewSock("test.sock")
	reply, err := sock.Send(&Command{
		Code:   "stop",
		Params: map[string]interface{}{},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("reply:", reply)
}
