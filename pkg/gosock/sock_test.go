// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package gosock

import (
	"encoding/json"
	"net"
	"testing"
	"time"
)

func TestSock_Listen(t *testing.T) {
	var sock = NewSock("test.sock")
	sock.OnCommand(func(cmd *Command) {
		cmdJSON, err := json.Marshal(cmd)
		if err != nil {
			t.Fatal(err)
		}
		if cmd.Code == "sleep" {
			time.Sleep(1 * time.Second)
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

func TestSock_SendTimeout(t *testing.T) {
	var sock = NewSock("test.sock")
	reply, err := sock.SendTimeout(&Command{
		Code:   "sleep",
		Params: map[string]interface{}{},
	}, 1*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("reply:", reply)
}

func TestSock_IsListening(t *testing.T) {
	var sock = NewSock("test.sock")
	t.Log(sock.IsListening())
}

func TestSock_Close(t *testing.T) {
	var sock = NewSock("test.sock")
	go func() {
		t.Log("listening ...")
		t.Log(sock.Listen())
		t.Log("end listening")
	}()
	time.Sleep(1 * time.Second)
	err := sock.Close()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(sock.Listen())
}
