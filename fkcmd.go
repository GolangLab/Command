package command

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

type FkCommand interface {
	Name() string
	Help() string
	Call(args []string) string
}

type CommandCPUProf struct{}

func (c *CommandCPUProf) Name() string {
	return "cpuprof"
}

func profileName() string {
	now := time.Now()
	return fmt.Sprintf("%d%02d%02d_%02d_%02d_%02d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second())
}

func (c *CommandCPUProf) Call(args []string) string {
	if len(args) == 0 {
		return c.Help()
	}
	command := args[0]
	switch command {
	case "start":
		fn := profileName() + ".cpuprof"
		f, err := os.Create(fn)
		if err != nil {
			return err.Error()
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			f.Close()
			return err.Error()
		}
		return fmt.Sprintln("profiling begin with", fn)
	case "stop":
		pprof.StopCPUProfile()
		return fmt.Sprintln("profiling stoped")
	default:
		return c.Help()
	}
}

func (c *CommandCPUProf) Help() string {
	return fmt.Sprintln("cpuprof start (to begin)") + fmt.Sprintln("cpuprof stop (to end)")
}

type CommandProf struct{}

func (c *CommandProf) Name() string {
	return "prof"
}

func (c *CommandProf) Call(args []string) string {
	if len(args) == 0 {
		return c.usage()
	}

	var (
		p  *pprof.Profile
		fn string
	)
	switch args[0] {
	case "goroutine":
		p = pprof.Lookup("goroutine")
		fn = profileName() + ".gprof"
	case "heap":
		p = pprof.Lookup("heap")
		fn = profileName() + ".hprof"
	case "thread":
		p = pprof.Lookup("threadcreate")
		fn = profileName() + ".tprof"
	case "block":
		p = pprof.Lookup("block")
		fn = profileName() + ".bprof"
	default:
		return c.usage()
	}

	f, err := os.Create(fn)
	if err != nil {
		return err.Error()
	}
	defer f.Close()
	err = p.WriteTo(f, 0)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintln("saving", fn)
}

func (c *CommandProf) Help() string {
	return fmt.Sprintln("writes a pprof-formatted snapshot")
}

func (c *CommandProf) usage() string {
	return "" +
		fmt.Sprintln("prof writes runtime profiling data in the format expected by") +
		fmt.Sprintln("the pprof visualization tool") +
		fmt.Sprintln("Usage: prof goroutine|heap|thread|block") +
		fmt.Sprintln("  goroutine - stack traces of all current goroutines") +
		fmt.Sprintln("  heap      - a sampling of all heap allocations") +
		fmt.Sprintln("  thread    - stack traces that led to the creation of new OS threads") +
		fmt.Sprintln("  block     - stack traces that led to blocking on synchronization primitives")
}
