# gosock
Send json command through sock in golang

# Server
~~~go
var sock = NewSock("test.sock")
sock.OnCommand(func(cmd *Command) {
	log.Println(cmd)
	_ = cmd.ReplyOk()
})
err := sock.Listen()
if err != nil {
	log.Fatal(err)
}
~~~

# Client
~~~go
var sock = NewSock("test.sock")
_, err := sock.Send(&Command{
	Code:   "stop",
	Params: map[string]interface{}{},
})
if err != nil {
	log.Fatal(err)
}
~~~