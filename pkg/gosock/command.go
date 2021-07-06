// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package gosock

import (
	"encoding/json"
	"net"
)

type Command struct {
	Code   string                 `json:"code"`
	Params map[string]interface{} `json:"params"`

	conn net.Conn
}

func (this *Command) Reply(reply *Command) error {
	if this.conn == nil {
		return nil
	}
	replyJSON, err := json.Marshal(reply)
	if err != nil {
		return err
	}
	replyJSON = append(replyJSON, '\n')
	_, err = this.conn.Write(replyJSON)
	if err != nil {
		return err
	}
	return nil
}

func (this *Command) ReplyOk() error {
	return this.Reply(&Command{
		Code: "ok",
	})
}
