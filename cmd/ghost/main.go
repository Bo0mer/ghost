// Program ghost is a command-line interface that interacts with ghost handlers.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Bo0mer/ghost"
)

var (
	remote string
)

func init() {
	flag.StringVar(&remote, "g", "", "Remote Ghost address")
}

const usage = `Usage: ghost is a tool to enable/disable monitoring targets.

Commands:
	monitor		Activate monitoring of a target.
	unmonitor	Deactivate monitoring of a target.
	targets		List available monitoring targets and their status.

Flags:
	-g		Specify remote ghost address.`

func main() {
	flag.Parse()
	flag.Usage = func() { printUsage("") }

	action := flag.Arg(0)
	if action == "" {
		printUsage("missing action")
	}

	if remote == "" {
		printUsage("missing remote ghost address")
	}

	cmd, ok := commands[action]
	if !ok {
		printUsage("unknown action")
	}

	if err := cmd(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

var commands = map[string]func() error{
	"monitor":   monitor,
	"unmonitor": unmonitor,
	"targets":   targets,
}

func monitor() error {
	target := flag.Arg(1)
	if target == "" {
		return errors.New("monitor: missing target")
	}

	return ghost.EnableMonitor(remote, target)
}
func unmonitor() error {
	target := flag.Arg(1)
	if target == "" {
		return errors.New("monitor: missing target")
	}

	return ghost.DisableMonitor(remote, target)
}
func targets() error {
	monitors, err := ghost.Monitors(remote)
	if err != nil {
		return err
	}
	var status string
	fmt.Println("Name\t\tStatus")
	for monitor, enabled := range monitors {
		status = "disabled"
		if enabled {
			status = "enabled"
		}

		fmt.Printf("%v\t\t%v\n", monitor, status)
	}
	return nil
}

func printUsage(msg string) {
	if msg != "" {
		fmt.Printf("ghost: %v\n", msg)
	}
	fmt.Fprintf(os.Stderr, "%v\n", usage)
	os.Exit(1)
}
