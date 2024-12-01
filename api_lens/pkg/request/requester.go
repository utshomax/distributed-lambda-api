package request

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"api-lens/pkg/config"
	"api-lens/pkg/dns"
	metricstypes "api-lens/pkg/types"

	"api-lens/pkg/httpstat"
)

func SendRequest(ctx context.Context, config config.RequestConfig) metricstypes.RequestMetrics {
	var requestMetrics metricstypes.RequestMetrics

	parsedURL, err := url.Parse(config.URL)
	if err != nil {
		return metricstypes.RequestMetrics{
			Success:      false,
			ErrorMessage: err.Error(),
		}
	}

	dnsResult, err := dns.Resolve(ctx, parsedURL.Hostname(), config.DNSServers)
	if err != nil {
		return metricstypes.RequestMetrics{
			Success:      false,
			ErrorMessage: fmt.Sprintf("DNS Lookup Error: %v", err),
		}
	}
	requestMetrics.DNSLookup = float64(dnsResult.ResolutionTime.Milliseconds())

	// Create httpstat result
	var result httpstat.Result
	ctx = httpstat.WithHTTPStat(ctx, &result)

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, config.Method, parsedURL.String(), nil)
	if err != nil {
		return metricstypes.RequestMetrics{
			Success:      false,
			ErrorMessage: err.Error(),
		}
	}

	// Add headers
	for k, v := range config.Headers {
		req.Header.Add(k, v)
	}

	// Custom dialer to use resolved IP
	dialer := &net.Dialer{
		Timeout: 30 * time.Second,
	}

	transport := &http.Transport{
		DisableKeepAlives: config.DisableKeepAlive,
		TLSClientConfig:   &tls.Config{},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				host = addr
				port = "80"
				if parsedURL.Scheme == "https" {
					port = "443"
				}
			}

			// Only replace the hostname part with the resolved IP
			if host == parsedURL.Hostname() {
				if config.DisableDNSCache {
					// Perform new DNS lookup for each request
					dnsResult, err := dns.Resolve(ctx, parsedURL.Hostname(), config.DNSServers)
					if err != nil {
						return nil, err
					}
					addr = net.JoinHostPort(dnsResult.IP, port)
				} else {
					addr = net.JoinHostPort(dnsResult.IP, port)
				}
			}

			return dialer.DialContext(ctx, network, addr)
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	// Ensure transport is closed after we're done
	defer transport.CloseIdleConnections()

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		return metricstypes.RequestMetrics{
			Success:       false,
			ErrorMessage:  err.Error(),
			DNSLookup:     requestMetrics.DNSLookup,
			TCPConnection: float64(result.TCPConnection.Milliseconds()),
			TLSHandshake:  float64(result.TLSHandshake.Milliseconds()),
		}
	}

	// Ensure response body is always closed
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return metricstypes.RequestMetrics{
			Success:       false,
			ErrorMessage:  fmt.Sprintf("Failed to read response body: %v", err),
			DNSLookup:     requestMetrics.DNSLookup,
			TCPConnection: float64(result.TCPConnection.Milliseconds()),
			TLSHandshake:  float64(result.TLSHandshake.Milliseconds()),
			FirstByte:     float64(result.ServerProcessing.Milliseconds()),
		}
	}

	result.End(time.Now()) // End timing

	requestMetrics.Success = true
	requestMetrics.StatusCode = resp.StatusCode
	requestMetrics.ResponseSize = len(body)
	requestMetrics.TCPConnection = float64(result.TCPConnection.Milliseconds())
	requestMetrics.TLSHandshake = float64(result.TLSHandshake.Milliseconds())
	requestMetrics.FirstByte = float64(result.StartTransfer.Milliseconds())
	requestMetrics.ContentTransfer = float64(result.ContentTransfer.Milliseconds())
	requestMetrics.ServerProcessing = float64(result.ServerProcessing.Milliseconds())
	requestMetrics.TotalTime = float64(result.Total.Milliseconds())
	return requestMetrics
}
