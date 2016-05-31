package command

import (
	"testing"
)

var quit = make(chan int)

func Test_command(t *testing.T) {
	server := NewCommandServer()
	addr := "127.0.0.1:12345"
	err := server.OpenServer(addr)
	if err != nil {
		t.Log("server open failed")
		return
	}
	<-quit
}
