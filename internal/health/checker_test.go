package health

import (
	"strings"
	"testing"
)

func TestRenderStatus_Pass(t *testing.T) {
	out := renderStatus("Live (A/AAAA)", true)
	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS in output, got: %q", out)
	}
	if !strings.Contains(out, "Live (A/AAAA)") {
		t.Errorf("expected label in output, got: %q", out)
	}
}

func TestRenderStatus_Fail(t *testing.T) {
	out := renderStatus("Sender Policy (SPF)", false)
	if !strings.Contains(out, "FAIL") {
		t.Errorf("expected FAIL in output, got: %q", out)
	}
	if !strings.Contains(out, "Sender Policy (SPF)") {
		t.Errorf("expected label in output, got: %q", out)
	}
}

func TestRenderScoreBar(t *testing.T) {
	tests := []struct {
		passed, total int
		wantLabel     string
	}{
		{4, 4, "4/4"},
		{0, 4, "0/4"},
		{2, 4, "2/4"},
		{3, 4, "3/4"},
	}

	for _, tt := range tests {
		out := renderScoreBar(tt.passed, tt.total)
		if !strings.Contains(out, tt.wantLabel) {
			t.Errorf("renderScoreBar(%d, %d): expected %q in output, got: %q",
				tt.passed, tt.total, tt.wantLabel, out)
		}
	}
}

func TestRenderReport_AllPassing(t *testing.T) {
	h := &DomainHealth{
		Domain:      "example.com",
		HasARecord:  true,
		HasMXRecord: true,
		HasSPF:      true,
		HasDMARC:    true,
	}
	out := h.RenderReport(false)

	if !strings.Contains(out, "example.com") {
		t.Error("expected domain name in report")
	}
	if !strings.Contains(out, "4/4") {
		t.Error("expected 4/4 score in report")
	}
	if !strings.Contains(out, "All checks passed") {
		t.Error("expected all-clear message in report")
	}
	if strings.Contains(out, "Warnings") {
		t.Error("unexpected Warnings section in all-passing report")
	}
	if strings.Contains(out, "Critical Issues") {
		t.Error("unexpected Critical Issues section in all-passing report")
	}
}

func TestRenderReport_WithIssues(t *testing.T) {
	h := &DomainHealth{
		Domain:  "broken.example",
		Issues:  []string{"No MX records found. The domain cannot receive email."},
	}
	out := h.RenderReport(false)

	if !strings.Contains(out, "Critical Issues") {
		t.Error("expected Critical Issues section")
	}
	if !strings.Contains(out, "No MX records found") {
		t.Error("expected issue message in report")
	}
	if !strings.Contains(out, "0/4") {
		t.Error("expected 0/4 score in report")
	}
	if strings.Contains(out, "All checks passed") {
		t.Error("unexpected all-clear message when issues present")
	}
}

func TestRenderReport_WithWarnings(t *testing.T) {
	h := &DomainHealth{
		Domain:      "warn.example",
		HasARecord:  true,
		HasMXRecord: true,
		HasSPF:      true,
		HasDMARC:    true,
		Warnings:    []string{"DMARC policy is set to 'p=none'."},
	}
	out := h.RenderReport(false)

	if !strings.Contains(out, "Warnings") {
		t.Error("expected Warnings section")
	}
	if !strings.Contains(out, "p=none") {
		t.Error("expected warning message in report")
	}
	if strings.Contains(out, "All checks passed") {
		t.Error("unexpected all-clear message when warnings present")
	}
}

func TestRenderReport_ContainsHint(t *testing.T) {
	h := &DomainHealth{Domain: "example.com"}
	out := h.RenderReport(true)
	if !strings.Contains(out, "Press Enter or Esc to exit") {
		t.Error("expected exit hint in report")
	}
}
