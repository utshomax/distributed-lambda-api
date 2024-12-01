package metricstypes

type RequestMetrics struct {
	Success          bool    `json:"success"`
	StatusCode       int     `json:"statusCode,omitempty"`
	ResponseSize     int     `json:"responseSize,omitempty"`
	DNSLookup        float64 `json:"dnsLookup"`
	TCPConnection    float64 `json:"tcpConnection"`
	TLSHandshake     float64 `json:"tlsHandshake"`
	ServerProcessing float64 `json:"serverProcessing"`
	FirstByte        float64 `json:"firstByte"`
	ContentTransfer  float64 `json:"contentTransfer"`
	TotalTime        float64 `json:"totalTime"`
	ErrorMessage     string  `json:"errorMessage,omitempty"`
}

type MetricsCollection struct {
	Metrics      []RequestMetrics `json:"metrics"`
	RequestCount int              `json:"requestCount"`
}
