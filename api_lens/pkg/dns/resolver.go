package dns

import (
	"context"
	"fmt"
	"net"
	"time"
)

type DNSResult struct {
	IP             string
	ResolutionTime time.Duration
}

func Resolve(ctx context.Context, hostname string, dnsServers []string) (*DNSResult, error) {
	// Default to Cloudflare's DNS if no servers specified
	if len(dnsServers) == 0 {
		dnsServers = []string{"1.1.1.1"}
	}

	start := time.Now()
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: time.Second * 10}
			return d.DialContext(ctx, network, fmt.Sprintf("%s:53", dnsServers[0]))
		},
	}

	ips, err := r.LookupIP(ctx, "ip4", hostname)
	if err != nil {
		return nil, err
	}

	return &DNSResult{
		IP:             ips[0].String(),
		ResolutionTime: time.Since(start),
	}, nil
}
