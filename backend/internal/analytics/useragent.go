package analytics

import "strings"

// UAResult holds parsed User-Agent information.
type UAResult struct {
	Browser    string
	OS         string
	DeviceType string // "Desktop", "Mobile", "Tablet"
}

// ParseUserAgent extracts browser, OS, and device type from a User-Agent string
// without any external dependencies.
func ParseUserAgent(ua string) UAResult {
	r := UAResult{
		Browser:    "Unknown",
		OS:         "Unknown",
		DeviceType: "Desktop",
	}
	if ua == "" {
		return r
	}

	lower := strings.ToLower(ua)

	// Detect device type
	switch {
	case strings.Contains(lower, "ipad") || strings.Contains(lower, "tablet") || strings.Contains(lower, "kindle"):
		r.DeviceType = "Tablet"
	case strings.Contains(lower, "mobile") || strings.Contains(lower, "android") && !strings.Contains(lower, "tablet"):
		r.DeviceType = "Mobile"
	case strings.Contains(lower, "iphone"):
		r.DeviceType = "Mobile"
	}

	// Detect OS (check mobile OS before desktop to avoid iPad matching "Mac OS X")
	switch {
	case strings.Contains(lower, "iphone") || strings.Contains(lower, "ipad"):
		r.OS = "iOS"
	case strings.Contains(lower, "android"):
		r.OS = "Android"
	case strings.Contains(lower, "windows"):
		r.OS = "Windows"
	case strings.Contains(lower, "mac os x") || strings.Contains(lower, "macintosh"):
		r.OS = "macOS"
	case strings.Contains(lower, "linux"):
		r.OS = "Linux"
	case strings.Contains(lower, "cros"):
		r.OS = "ChromeOS"
	}

	// Detect browser (order matters: check specific before generic)
	switch {
	case strings.Contains(lower, "edg/") || strings.Contains(lower, "edge/"):
		r.Browser = "Edge"
	case strings.Contains(lower, "opr/") || strings.Contains(lower, "opera"):
		r.Browser = "Opera"
	case strings.Contains(lower, "brave"):
		r.Browser = "Brave"
	case strings.Contains(lower, "vivaldi"):
		r.Browser = "Vivaldi"
	case strings.Contains(lower, "firefox"):
		r.Browser = "Firefox"
	case strings.Contains(lower, "samsungbrowser"):
		r.Browser = "Samsung Internet"
	case strings.Contains(lower, "chrome") && !strings.Contains(lower, "chromium"):
		r.Browser = "Chrome"
	case strings.Contains(lower, "safari") && !strings.Contains(lower, "chrome"):
		r.Browser = "Safari"
	case strings.Contains(lower, "msie") || strings.Contains(lower, "trident"):
		r.Browser = "IE"
	}

	return r
}
