// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package gosock

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"os"
	"time"
)

var ErrSockInUse = errors.New("bind: address already in use")

type Sock struct {
	path     string
	listener net.Listener

	onCommandFunc func(cmd *Command)
	onErrorFunc   func(err error)

	callbackMap map[string]func(params map[string]interface{})
	ticker      *time.Ticker
}

func NewSock(path string) *Sock {
	var sock = &Sock{
		path:        path,
		callbackMap: map[string]func(params map[string]interface{}){},
		ticker:      time.NewTicker(5 * time.Second),
	}
	go func() {
		for range sock.ticker.C {
			sock.fix()
		}
	}()
	return sock
}

func NewTmpSock(name string) *Sock {
	return NewSock(os.TempDir() + "/" + name)
}

func (this *Sock) Listen() error {
	// 是否正在使用
	_, err := this.Dial()
	if err == nil {
		return ErrSockInUse
	}

	f, err := os.Stat(this.path)
	if err == nil && !f.IsDir() {
		_ = os.Remove(this.path)
	}

	listener, err := net.Listen("unix", this.path)
	if err != nil {
		return err
	}

	this.listener = listener
	for {
		conn, err := listener.Accept()
		if err != nil {
			// 重新监听
			listener, err = net.Listen("unix", this.path)
			if err != nil {
				return err
			}
			this.listener = listener
			continue
		}
		go this.handle(conn, func(cmd *Command) {
			if this.onCommandFunc != nil {
				this.onCommandFunc(cmd)
			}

			callback, ok := this.callbackMap[cmd.Code]
			if ok {
				callback(cmd.Params)
			}
		})
	}
}

func (this *Sock) Dial() (net.Conn, error) {
	return net.Dial("unix", this.path)
}

func (this *Sock) OnCommand(f func(cmd *Command)) {
	this.onCommandFunc = f
}

func (this *Sock) On(cmd string, f func(params map[string]interface{})) {
	this.callbackMap[cmd] = f
}

func (this *Sock) OnError(f func(err error)) {
	this.onErrorFunc = f
}

func (this *Sock) Send(cmd *Command) (reply *Command, err error) {
	conn, err := this.Dial()
	if err != nil {
		return nil, err
	}
	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	cmdJSON = append(cmdJSON, '\n')
	_, err = conn.Write(cmdJSON)
	if err != nil {
		return nil, err
	}

	this.handle(conn, func(cmd *Command) {
		reply = cmd
		_ = conn.Close()
	})

	return reply, nil
}

func (this *Sock) IsListening() bool {
	conn, err := this.Dial()
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func (this *Sock) Close() error {
	if this.listener == nil {
		return nil
	}
	return this.listener.Close()
}

func (this *Sock) handle(conn net.Conn, callback func(cmd *Command)) {
	var buf = make([]byte, 1024)
	var cmdBuf = []byte{}
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			newLine := bytes.Index(buf[:n], []byte{'\n'})
			if newLine < 0 {
				cmdBuf = append(cmdBuf, buf[:n]...)
			} else {
				cmdBuf = append(cmdBuf, buf[:n][:newLine]...)

				var cmdObj = &Command{}
				err = json.Unmarshal(cmdBuf, cmdObj)
				if err != nil {
					if this.onErrorFunc != nil {
						this.onErrorFunc(err)
					}
				} else {
					cmdObj.conn = conn
					callback(cmdObj)
				}

				cmdBuf = buf[:n][newLine+1:]
			}
		}
		if err != nil {
			return
		}
	}
}

// fix sock file
func (this *Sock) fix() {
	if this.listener == nil {
		return
	}
	_, err := os.Stat(this.path)
	if os.IsNotExist(err) {
		_ = this.listener.Close()
	}
}
