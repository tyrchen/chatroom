package chat

import (
	"errors"
	"fmt"
	"github.com/tyrchen/goutil/regex"
	"regexp"
)

type Command struct {
	cmd string
	arg string
}

type Run func(server *Server, client *Client, arg string)

const (
	CMD_REGEX = `:(?P<cmd>\w+)\s*(?P<arg>.*)`
)

var (
	commands map[string]Run
)

func init() {
	commands = map[string]Run{
		"name": changeName,
		"quit": quit,
	}
}

func parseCommand(msg string) (cmd Command, err error) {
	r := regexp.MustCompile(CMD_REGEX)
	if values, ok := regex.MatchAll(r, msg); ok {
		cmd.cmd, _ = values[0]["cmd"]
		cmd.arg, _ = values[0]["arg"]
		return
	}
	err = errors.New("Unparsed message: " + msg)
	return
}

func (self *Server) executeCommand(client *Client, cmd Command) (err error) {
	if f, ok := commands[cmd.cmd]; ok {
		f(self, client, cmd.arg)
		return
	}

	err = errors.New("Unsupported command: " + cmd.cmd)
	return
}

// commands

func changeName(server *Server, client *Client, arg string) {
	oldname := client.GetName()
	client.SetName(arg)
	server.broadcast(fmt.Sprintf("Notification: %s changed its name to %s", oldname, arg))
}

func quit(server *Server, client *Client, arg string) {
	client.quit()
	server.broadcast(fmt.Sprintf("Notification: %s quit the chat room.", client.GetName()))
}
