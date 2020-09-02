package gocli

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
)

var (
	commandNotExistsMessage = "command not exists"
	helpDescription         = "is this help message"
	returnDescription       = "is return upper to menu"
	exitDescription         = "is exit this cli"
	interruptMessage        = "input exit"
)

func SetCommandNotExistsMessage(msg string) { commandNotExistsMessage = msg }
func SetHelpDescription(msg string)         { helpDescription = msg }
func SetReturnDescription(msg string)       { returnDescription = msg }
func SetExitDescription(msg string)         { exitDescription = msg }
func SetInterruptMessage(msg string)        { interruptMessage = msg }
func ThroughInterrupt() {
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt)
		for range ch {
			fmt.Println(interruptMessage)
		}
	}()
}

func NewGocli() Gocli {
	bs := bufio.NewScanner(os.Stdin)
	bs.Split(bufio.ScanLines)
	cmd := NewCommand("", "")

	return Gocli{Command: cmd, current: cmd, bs: bs, path: []string{""}}
}

type Gocli struct {
	Command
	current Command
	path    []string
	bs      *bufio.Scanner
}

func (g *Gocli) Run() error {
	for {
		fmt.Printf("%s >>> ", strings.Join(g.path, "/"))
		g.bs.Scan()
		cmd := strings.TrimSpace(g.bs.Text())
		c := g.current.getCommand(cmd)
		if c == nil {
			fmt.Println(commandNotExistsMessage)
			continue
		}
		g.current = c
		g.path = append(g.path, g.current.getCmd())

		if g.current.hasAction() {
			switch g.current.run(g.bs) {
			case AfterActionReturn:
				if g.current.getParent() != nil && g.current.getCmd() != "" {
					g.current = g.current.getParent()
					g.path = g.path[:len(g.path)-1]
				}
			case AfterActionReturnTwice:
				if g.current.getParent() != nil && g.current.getCmd() != "" {
					g.current = g.current.getParent()
					g.path = g.path[:len(g.path)-1]
				}
				if g.current.getParent() != nil && g.current.getCmd() != "" {
					g.current = g.current.getParent()
					g.path = g.path[:len(g.path)-1]
				}
			case AfterActionExit:
				return nil
			}
		}
		fmt.Println()
	}
}

func NewCommand(name string, desc string) Command {
	c := &command{cmd: name, desc: desc, commands: map[string]Command{}}

	// helpを追加
	c.AddSubCommand(&command{cmd: "help", desc: helpDescription, action: func(bs *bufio.Scanner) AfterAction {
		c.printHelp()
		return AfterActionReturn
	}})

	// returnを追加
	c.AddSubCommand(&command{cmd: "return", desc: returnDescription, action: func(bs *bufio.Scanner) AfterAction {
		return AfterActionReturnTwice
	}})

	// exitを追加
	c.AddSubCommand(&command{cmd: "exit", desc: exitDescription, action: func(bs *bufio.Scanner) AfterAction {
		return AfterActionExit
	}})

	return c
}

type Command interface {
	SetAction(action Action) Command
	AddSubCommand(command Command) Command
	getCmd() string
	getDescription() string
	getCommand(cmd string) Command
	setParent(parent Command)
	getParent() Command
	hasAction() bool
	run(bs *bufio.Scanner) AfterAction
}

type command struct {
	cmd       string
	desc      string
	action    Action
	cmds      []string
	maxCmdLen int
	commands  map[string]Command
	parent    Command
}

func (c *command) getCmd() string         { return c.cmd }
func (c *command) getDescription() string { return c.desc }
func (c *command) getCommand(cmd string) Command {
	if c, ok := c.commands[cmd]; ok {
		return c
	} else {
		return nil
	}
}
func (c *command) setParent(parent Command) {
	c.parent = parent
}
func (c *command) getParent() Command {
	return c.parent
}
func (c *command) hasAction() bool {
	return c.action != nil
}
func (c *command) run(bs *bufio.Scanner) AfterAction {
	return c.action(bs)
}
func (c *command) SetAction(action Action) Command {
	c.action = action
	return c
}
func (c *command) AddSubCommand(command Command) Command {
	command.setParent(c)
	if _, ok := c.commands[command.getCmd()]; !ok {
		c.cmds = append(c.cmds, command.getCmd())
		if c.maxCmdLen < len(command.getCmd()) {
			c.maxCmdLen = len(command.getCmd())
		}
	}
	c.commands[command.getCmd()] = command
	return c
}
func (c *command) printHelp() {
	for _, cmd := range c.cmds {
		fmt.Printf("%s %s %s\n", cmd, strings.Repeat(" ", c.maxCmdLen-len(cmd)), c.commands[cmd].getDescription())
	}
}

type Action func(bs *bufio.Scanner) AfterAction

type AfterAction string

const (
	AfterActionNone        AfterAction = ""
	AfterActionKeep        AfterAction = "keep"
	AfterActionReturn      AfterAction = "return"
	AfterActionReturnTwice AfterAction = "return twice"
	AfterActionExit        AfterAction = "exit"
)
