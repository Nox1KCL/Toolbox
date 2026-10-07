package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/Nox1KCL/Toolbox/internal/portcheck"
)

type Flags struct {
	port int
	kill bool
}

func main() {
	// TODO: зробити систему історії
	var flags Flags
	flags.Parse() // Моя комплексна функція для створення і парсингу флагів

	strPort := fmt.Sprintf(":%d", flags.port)
	cmd := exec.Command("ss", "-tlnp", "sport", "=", strPort)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("trying execute command %v: %v", cmd, err)
	}
	text := string(output)

	var ports []portcheck.PortData
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" || i == 0 {
			continue
		}
		var port portcheck.PortData
		elements := strings.Fields(line)

		port.Fill(elements, strPort)

		if flags.kill {
			if err := portcheck.KillProcess(port.PID); err != nil {
				port.Killed = "No"
			} else {
				port.Killed = "Yes"
			}
		} else {
			port.Killed = "Not requested"
		}

		ports = append(ports, port)
	}
	if len(ports) == 0 {
		fmt.Printf("No process is listening on port %d\n", flags.port)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	portcheck.FormMessage(w, ports)
	w.Flush()
}

func (f *Flags) Parse() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.IntVar(&f.port, "port", 0, "port that you need to find")
	flag.BoolVar(&f.kill, "kill", false, "kill port process you entered")
	flag.Parse()

	if f.port <= 0 || f.port > 65535 {
		fmt.Fprintln(os.Stderr, "Error: valid port (1-65535) is required")
		flag.Usage()
		os.Exit(1)
	}
}
