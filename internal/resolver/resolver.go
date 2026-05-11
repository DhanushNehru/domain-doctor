package resolver

import (
	"context"
	"fmt"
	"net"
	"time"
)

type Resolver interface {
	LookupHost(ctx context.Context, host string) ([]string, error)
	LookupMX(ctx context.Context, host string) ([]*net.MX, error)
	LookupTXT(ctx context.Context, host string) ([]string, error)
}

type NetResolver struct {
	r *net.Resolver
}

func (n *NetResolver) LookupHost(ctx context.Context, host string) ([]string, error) {
	return n.r.LookupHost(ctx, host)
}

func (n *NetResolver) LookupMX(ctx context.Context, host string) ([]*net.MX, error) {
	return n.r.LookupMX(ctx, host)
}

func (n *NetResolver) LookupTXT(ctx context.Context, host string) ([]string, error) {
	return n.r.LookupTXT(ctx, host)
}

func New(resolverName string) *NetResolver {
	if resolverName == "" {
		return &NetResolver{
			r: net.DefaultResolver,
		}
	}

	return &NetResolver{r: &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 10 * time.Second,
			}

			conn, err := d.DialContext(ctx, network, resolverName)
			if err != nil {
				return nil, fmt.Errorf("dns dial failed: %w", err)
			}

			return conn, nil
		},
	}}
}
