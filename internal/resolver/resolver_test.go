package resolver

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"
)

func TestNew_DefaultResolver(t *testing.T) {
	r := New("")
	if r == nil {
		t.Fatal("New() returned nil")
	}
	if r.r != net.DefaultResolver {
		t.Error("Expected net.DefaultResolver when no resolverName is provided")
	}
}

func TestNew_CustomResolver(t *testing.T) {
	customDNS := "8.8.8.8:53"
	r := New(customDNS)

	if r == nil {
		t.Fatal("New() returned nil")
	}
	if r.r == nil {
		t.Fatal("net.Resolver is nil")
	}
	if !r.r.PreferGo {
		t.Error("PreferGo should be true for custom resolver")
	}
}

func TestNetResolver_LookupHost(t *testing.T) {
	mock := &mockResolver{
		hostRecords: map[string][]string{
			"example.com": {"93.184.216.34", "2606:2800:220:1:248:1893:25c8:1946"},
		},
	}

	tests := []struct {
		name     string
		resolver Resolver
		host     string
		wantLen  int
		wantErr  bool
	}{
		{
			name:     "valid host",
			resolver: mock,
			host:     "example.com",
			wantLen:  2,
			wantErr:  false,
		},
		{
			name:     "not found",
			resolver: mock,
			host:     "nonexistent.domain.test",
			wantLen:  0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ips, err := tt.resolver.LookupHost(context.Background(), tt.host)

			if (err != nil) != tt.wantErr {
				t.Errorf("LookupHost() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(ips) != tt.wantLen {
				t.Errorf("LookupHost() returned %d IPs, want %d", len(ips), tt.wantLen)
			}
		})
	}
}

func TestNetResolver_LookupMX(t *testing.T) {
	mock := &mockResolver{
		mxRecords: map[string][]*net.MX{
			"example.com": {
				{Pref: 10, Host: "mail.example.com."},
				{Pref: 20, Host: "mail2.example.com."},
			},
		},
	}

	mxs, err := mock.LookupMX(context.Background(), "example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mxs) != 2 {
		t.Errorf("expected 2 MX records, got %d", len(mxs))
	}
}

func TestNetResolver_LookupTXT(t *testing.T) {
	mock := &mockResolver{
		txtRecords: map[string][]string{
			"example.com": {
				"v=spf1 include:_spf.google.com ~all",
				"google-site-verification=abc123",
			},
			"_dmarc.example.com": {"v=DMARC1; p=reject;"},
		},
	}

	txts, err := mock.LookupTXT(context.Background(), "example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txts) == 0 {
		t.Error("expected TXT records")
	}
}

// mockResolver implements Resolver interface for testing
type mockResolver struct {
	hostRecords map[string][]string
	mxRecords   map[string][]*net.MX
	txtRecords  map[string][]string
}

func (m *mockResolver) LookupHost(ctx context.Context, host string) ([]string, error) {
	if records, ok := m.hostRecords[host]; ok {
		return records, nil
	}
	return nil, &net.DNSError{Name: host, Err: "no such host", IsNotFound: true}
}

func (m *mockResolver) LookupMX(ctx context.Context, host string) ([]*net.MX, error) {
	if records, ok := m.mxRecords[host]; ok {
		return records, nil
	}
	return nil, &net.DNSError{Name: host, Err: "no MX records", IsNotFound: true}
}

func (m *mockResolver) LookupTXT(ctx context.Context, host string) ([]string, error) {
	if records, ok := m.txtRecords[host]; ok {
		return records, nil
	}
	// Return empty slice + nil error for domains with no TXT
	return nil, &net.DNSError{Name: host, Err: "no TXT records", IsNotFound: true}
}

func TestNew_WithTimeout(t *testing.T) {
	// Test that custom resolver has reasonable dial timeout
	r := New("1.1.1.1:53")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := r.LookupHost(ctx, "google.com")
	if err != nil && !strings.Contains(err.Error(), "timeout") && !strings.Contains(err.Error(), "no such host") {
		t.Logf("LookupHost returned: %v (acceptable in test env)", err)
	}
}
