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

	var cmd *exec.Cmd
	strPort := fmt.Sprintf(":%d", flags.port)
	cmd = exec.Command("ss", "-tlnp", "sport", "=", strPort)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("trying execute command %v: %v", cmd, err)
		os.Exit(1)
	}
	text := string(output)

	var ports []portcheck.PortData
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" || i == 0 {
			continue
		}
		var port portcheck.PortData
		column := strings.Fields(line)

		port.Fill(column, strPort)
		port.Kill(flags.kill)

		ports = append(ports, port)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	portcheck.FormingMessage(w, ports)
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

	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(0)
	}
}
