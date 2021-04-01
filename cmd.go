package goutil

import (
	"fmt"
	"os"
)

type CmdHandler map[string]func(args ...string)

func (h CmdHandler) handle(args []string) bool {
	if len(args) < 1 {
		fmt.Println("command entered with < 1 argument")
	}
	c := args[0]
	if f, ok := h[c]; ok {
		fmt.Println("handling cmd...")
		f(args[1:]...)
		return true
	}
	switch c {
	case "h", "-h":
		// printHelp()
	default:
		fmt.Printf("Unrecognized input: '%v'\n", args)
	}
	return false
}

func (h CmdHandler) Run() bool {
	args := os.Args[1:]
	if len(args) > 0 {
		return h.handle(args)
	}
	return false
}
