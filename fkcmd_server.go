package command

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

type FkCommandServer struct {
	commands map[string]FkCommand
}

type InnerCmd struct {
	name string
	call func([]string) string
	help func() string
}

func (i *InnerCmd) Name() string {
	return i.name
}

func (i *InnerCmd) Call(args []string) string {
	if i.call != nil {
		return i.call(args)
	}
	return fmt.Sprintln("done nothing")
}

func (i *InnerCmd) Help() string {
	if i.help != nil {
		return i.help()
	}
	return fmt.Sprintln("nothing to help")
}

func NewCommandServer() FkCommandServer {
	ret := FkCommandServer{commands: make(map[string]FkCommand)}
	ret.AddCommand(&InnerCmd{"help", ret.Useage, nil})
	ret.AddCommand(&InnerCmd{"quit", nil, nil})
	ret.AddCommand(&InnerCmd{"list", ret.ListCommnad, nil})

	ret.AddCommand(&CommandCPUProf{})
	ret.AddCommand(&CommandProf{})
	return ret
}

var SameNameError = errors.New("common with same name has added")

func (f *FkCommandServer) getCommand(name string) (c FkCommand) {
	return f.commands[name]
}

func (f *FkCommandServer) AddCommand(c FkCommand) error {
	old := f.getCommand(c.Name())
	if old != nil {
		return SameNameError
	}
	f.commands[c.Name()] = c
	return nil
}

func (f *FkCommandServer) ListCommnad(args []string) string {
	var ret string
	for _, c := range f.commands {
		ret += fmt.Sprintln(c.Name())
	}
	return ret
}

func (f *FkCommandServer) Useage(args []string) string {
	ret := "" +
		fmt.Sprintln("1. help (to print usage)") +
		fmt.Sprintln("2. list (to list all commands)") +
		fmt.Sprintln("3. quit (to exit)")
	return ret
}

func (f *FkCommandServer) handle_conn(c net.Conn) {
	c.Write([]byte(f.Useage(nil)))

	read := bufio.NewReader(c)
	defer c.Close()
	for {
		line, err := read.ReadString('\n')
		if err != nil {
			fmt.Println("read conn", c.RemoteAddr(), "with error", err)
			return
		}
		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}
		if args[0] == "quit" {
			c.Write([]byte("Bye!\n"))
			return
		}
		cm := f.getCommand(args[0])
		if cm == nil {
			c.Write([]byte(fmt.Sprintln("command", args[0], "not found")))
			continue
		}
		out := cm.Call(args[1:])
		c.Write([]byte(out))
	}
}

func (f *FkCommandServer) listen_server(l net.Listener) {
	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}
		go f.handle_conn(c)
		fmt.Println("handler command conn", c.RemoteAddr())
	}
}

func (f *FkCommandServer) OpenServer(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	go f.listen_server(l)
	fmt.Println("begin listen to server", addr)
	return nil
}
