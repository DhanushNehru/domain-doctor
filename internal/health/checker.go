package health

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dhanushnehru/domain-doctor/internal/resolver"
)

// DomainHealth contains the results of the health checks
type DomainHealth struct {
	Domain      string
	HasARecord  bool
	HasMXRecord bool
	HasSPF      bool
	HasDMARC    bool
	Issues      []string
	Warnings    []string
}

func SetResolver(dnsServer string) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 5 * time.Second,
			}

			conn, err := d.DialContext(ctx, network, dnsServer)
			if err != nil {
				return nil, fmt.Errorf("dns dial failed: %w", err)
			}

			return conn, nil
		},
	}
}

// Check performs DNS-based health checks on a domain
func Check(domain string, r resolver.NetResolver) *DomainHealth {
	health := &DomainHealth{
		Domain: domain,
	}

	var wg sync.WaitGroup
	wg.Add(3)

	// Check A/AAAA records (Liveness)
	go func() {
		defer wg.Done()
		ips, err := r.LookupHost(context.Background(), domain)
		if err == nil && len(ips) > 0 {
			health.HasARecord = true
		} else {
			health.Issues = append(health.Issues, "No A or AAAA records found. The domain might not be live.")
		}
	}()

	// Check MX records
	go func() {
		defer wg.Done()
		mxs, err := r.LookupMX(context.Background(), domain)
		if err == nil && len(mxs) > 0 {
			health.HasMXRecord = true
		} else {
			health.Issues = append(health.Issues, "No MX records found. The domain cannot receive email.")
		}
	}()

	// Check TXT records for SPF and DMARC
	go func() {
		defer wg.Done()
		txts, err := r.LookupTXT(context.Background(), domain)
		if err == nil {
			for _, txt := range txts {
				if strings.HasPrefix(txt, "v=spf1") {
					health.HasSPF = true
					if strings.Contains(txt, "+all") {
						health.Warnings = append(health.Warnings, "SPF record uses '+all', which allows anyone to send email on your behalf. Use '~all' or '-all'.")
					}
				}
			}
		} else {
			health.Issues = append(health.Issues, "Could not fetch TXT records for SPF.")
		}

		// Check DMARC (on _dmarc.domain)
		dmarcDomain := "_dmarc." + domain
		dmarcTxts, err := r.LookupTXT(context.Background(), dmarcDomain)
		if err == nil {
			for _, txt := range dmarcTxts {
				if strings.HasPrefix(txt, "v=DMARC1") {
					health.HasDMARC = true
					if strings.Contains(txt, "p=none") {
						health.Warnings = append(health.Warnings, "DMARC policy is set to 'p=none'. This does not prevent spoofing. Consider upgrading to 'p=quarantine' or 'p=reject'.")
					}
				}
			}
		} else {
			health.Warnings = append(health.Warnings, "No DMARC record found or could not be fetched.")
		}
	}()

	wg.Wait()

	if !health.HasSPF {
		health.Issues = append(health.Issues, "No SPF record found. Your emails are highly likely to go to spam.")
	}

	return health
}

func (h *DomainHealth) RenderReport(showHint bool) string {
	header := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Left, checkLabel.Render("Domain:"), domainStyle.Render(h.Domain)),
	)

	status := lipgloss.JoinVertical(
		lipgloss.Left,
		renderStatus("Live (A/AAAA)", h.HasARecord),
		renderStatus("Email Receiver (MX)", h.HasMXRecord),
		renderStatus("Sender Policy (SPF)", h.HasSPF),
		renderStatus("Domain Auth (DMARC)", h.HasDMARC),
	)

	passed := 0
	for _, c := range []bool{h.HasARecord, h.HasMXRecord, h.HasSPF, h.HasDMARC} {
		if c {
			passed++
		}
	}

	score := renderScoreBar(passed, 4)

	sections := []string{
		header,
		status,
		dividerStyle.Render(strings.Repeat("─", 50)),
		score,
	}

	if len(h.Warnings) > 0 {
		warnList := []string{warnTitleStyle.Render("⚠  Warnings")}
		for _, w := range h.Warnings {
			row := lipgloss.JoinHorizontal(
				lipgloss.Left,
				warnBullet.Render("."),
				w,
			)
			warnList = append(warnList, row)
		}
		sections = append(sections, lipgloss.JoinVertical(lipgloss.Left, warnList...))
	}

	if len(h.Issues) > 0 {
		issueList := []string{critTitleStyle.Render("✖  Critical Issues")}
		for _, i := range h.Issues {
			row := lipgloss.JoinHorizontal(
				lipgloss.Left,
				critBullet.Render("."),
				i,
			)
			issueList = append(issueList, row)
		}
		sections = append(sections, lipgloss.JoinVertical(lipgloss.Left, issueList...))
	}

	if len(h.Issues) == 0 && len(h.Warnings) == 0 {
		sections = append(
			sections,
			okStyle.Render("✅  All checks passed. Your domain is healthy!"),
		)
	}

	if showHint {
		sections = append(sections,
			hintStyle.Render("Press Enter or Esc to exit"),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		sections...,
	)
}

func renderStatus(name string, passed bool) string {
	var badge string
	if passed {
		badge = passBadge.Render("✓ PASS")
	} else {
		badge = failBadge.Render("✗ FAIL")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		badge,
		checkLabel.Render(name),
	)
}

func renderScoreBar(passed, total int) string {
	filled := scoreBarFill.Render(strings.Repeat("█", passed))
	empty := scoreBarEmpty.Render(strings.Repeat("░", total-passed))

	bar := lipgloss.JoinHorizontal(lipgloss.Left, filled, empty)
	label := scoreBarLabel.Render(fmt.Sprintf(" %d/%d checks passed", passed, total))

	return lipgloss.JoinHorizontal(lipgloss.Left, bar, label)
}
