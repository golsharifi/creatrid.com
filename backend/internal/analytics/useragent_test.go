package analytics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUserAgent(t *testing.T) {
	tests := []struct {
		name     string
		ua       string
		expected UAResult
	}{
		{
			name: "Chrome on Windows",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			expected: UAResult{
				Browser:    "Chrome",
				OS:         "Windows",
				DeviceType: "Desktop",
			},
		},
		{
			name: "Safari on macOS",
			ua:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
			expected: UAResult{
				Browser:    "Safari",
				OS:         "macOS",
				DeviceType: "Desktop",
			},
		},
		{
			name: "Firefox on Linux",
			ua:   "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
			expected: UAResult{
				Browser:    "Firefox",
				OS:         "Linux",
				DeviceType: "Desktop",
			},
		},
		{
			name: "Edge on Windows",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
			expected: UAResult{
				Browser:    "Edge",
				OS:         "Windows",
				DeviceType: "Desktop",
			},
		},
		{
			name: "Chrome on Android Mobile",
			ua:   "Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.43 Mobile Safari/537.36",
			expected: UAResult{
				Browser:    "Chrome",
				OS:         "Android",
				DeviceType: "Mobile",
			},
		},
		{
			name: "Safari on iPad",
			ua:   "Mozilla/5.0 (iPad; CPU OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
			expected: UAResult{
				Browser:    "Safari",
				OS:         "iOS",
				DeviceType: "Tablet",
			},
		},
		{
			name: "Empty string",
			ua:   "",
			expected: UAResult{
				Browser:    "Unknown",
				OS:         "Unknown",
				DeviceType: "Desktop",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseUserAgent(tt.ua)
			assert.Equal(t, tt.expected.Browser, result.Browser, "Browser mismatch")
			assert.Equal(t, tt.expected.OS, result.OS, "OS mismatch")
			assert.Equal(t, tt.expected.DeviceType, result.DeviceType, "DeviceType mismatch")
		})
	}
}
