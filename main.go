package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhanushnehru/domain-doctor/internal/health"
	"github.com/dhanushnehru/domain-doctor/internal/resolver"
	"github.com/dhanushnehru/domain-doctor/pkg/tui"
)

func main() {
	domain := flag.String("domain", "", "Domain to analyze directly, skipping the interactive prompt")
	resolverFlag := flag.String("resolver", "", "Custom DNS resolver (e.g. 1.1.1.1:53)")
	flag.Parse()

	r := resolver.New(*resolverFlag)

	if *domain != "" {
		report := health.Check(*domain, *r)
		fmt.Print(report.RenderReport(true))
		return
	}

	p := tea.NewProgram(tui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
