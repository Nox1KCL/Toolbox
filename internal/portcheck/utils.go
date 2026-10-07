package portcheck

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"text/tabwriter"
)

type PortData struct {
	State            string
	LocalAddressPort string
	ProcessName      string
	PID              string
	FD               string
	Killed           string
}

func FormingMessage(w *tabwriter.Writer, ports []PortData) {
	fmt.Fprintln(w, "№\tState\tProcess\tLocal Address\tPID\tFD\tKILLED")
	for i, info := range ports {
		column := fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%s\t%s\n",
			i+1,
			info.State,
			info.ProcessName,
			info.LocalAddressPort,
			info.PID,
			info.FD,
			info.Killed,
		)
		fmt.Fprint(w, column)

		if i < len(ports)-1 {
			fmt.Fprint(w, "")
		}
	}
	fmt.Fprintln(w)

}

func (p *PortData) Fill(column []string, port string) {
	p.State = column[0]

	for i, v := range column {
		if strings.Contains(v, port) {
			p.LocalAddressPort = column[i]
		} else if strings.Contains(v, "users") {
			users := strings.Split(column[i], "(")
			data := strings.Split(users[2], ")")
			complete := strings.Split(data[0], ",")

			p.ProcessName = complete[0]
			p.PID = complete[1]
			p.FD = complete[2]
		}
	}
}

func (p *PortData) Kill(kill bool) {
	if kill {
		strPID := fmt.Sprintf("%s", strings.Split(p.PID, "=")[1])
		cmd := exec.Command("kill", strPID)
		err := cmd.Run()
		if err != nil {
			log.Print(err)
			p.Killed = "No"
		} else {
			p.Killed = "Yes"
		}
	} else {
		p.Killed = "Not requested"
	}

}
