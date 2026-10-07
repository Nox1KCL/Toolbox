package portcheck

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var userRegex = regexp.MustCompile(`users:\(\("([^"]+)",pid=(\d+),fd=(\d+)\)`)

type PortData struct {
	State            string
	LocalAddressPort string
	ProcessName      string
	PID              string
	FD               string
	Killed           string
}

func FormMessage(w io.Writer, ports []PortData) {
	fmt.Fprintln(w, "№\tState\tProcess\tLocal Address\tPID\tFD\tKILLED")
	for i, info := range ports {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n",
			i+1,
			info.State,
			info.ProcessName,
			info.LocalAddressPort,
			info.PID,
			info.FD,
			info.Killed,
		)
	}
	fmt.Fprintln(w)

}

func (p *PortData) Fill(elements []string, port string) {
	if len(elements) == 0 {
		return
	}

	p.State = elements[0]
	for _, v := range elements {
		if strings.Contains(v, port) {
			p.LocalAddressPort = v
		}
		if matches := userRegex.FindStringSubmatch(v); len(matches) == 4 {
			p.ProcessName = matches[1]
			p.PID = matches[2]
			p.FD = matches[3]
		}
	}

}

func KillProcess(processPid string) error {
	pid, err := strconv.Atoi(processPid)
	if err != nil {
		return err
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return proc.Kill()
}
