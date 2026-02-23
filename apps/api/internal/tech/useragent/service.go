package useragent

import (
	"regexp"
	"strings"
)

// botPatterns contains substrings found in bot/crawler user agents (lowercase).
var botPatterns = []string{
	"bot", "crawler", "spider", "slurp", "archiver", "fetch",
	"scraper", "wget", "curl", "python-requests", "python-urllib",
	"java/", "httpclient", "go-http-client", "libwww", "lwp-",
	"monitoring", "checker", "validator", "scanner", "probe",
	"headlesschrome", "phantomjs",
}

// regexes for version extraction.
var (
	reChromeVer  = regexp.MustCompile(`(?i)Chrome/(\d+\.\d+)`)
	reFirefoxVer = regexp.MustCompile(`(?i)Firefox/(\d+\.\d+)`)
	reSafariVer  = regexp.MustCompile(`(?i)Version/(\d+\.\d+)`)
	reEdgeVer    = regexp.MustCompile(`(?i)Edg(?:e)?/(\d+\.\d+)`)
	reOperaVer   = regexp.MustCompile(`(?i)OPR/(\d+\.\d+)`)
	reIEVer      = regexp.MustCompile(`(?i)(?:MSIE |rv:)(\d+\.\d+)`)

	reWindowsVer = regexp.MustCompile(`Windows NT (\d+\.\d+)`)
	reMacVer     = regexp.MustCompile(`Mac OS X (\d+[._]\d+)`)
	reAndroidVer = regexp.MustCompile(`Android (\d+(?:\.\d+)?)`)
	reiOSVer     = regexp.MustCompile(`OS (\d+[._]\d+)`)
)

// windowsVersions maps NT version strings to marketing names.
var windowsVersions = map[string]string{
	"10.0": "10/11",
	"6.3":  "8.1",
	"6.2":  "8",
	"6.1":  "7",
	"6.0":  "Vista",
	"5.2":  "XP x64",
	"5.1":  "XP",
}

// Service parses user agent strings.
type Service struct{}

// NewService constructs a new Service.
func NewService() *Service { return &Service{} }

// Parse extracts browser, OS, device, and bot information from a UA string.
func (s *Service) Parse(ua string) Result {
	if ua == "" {
		return Result{Device: "unknown"}
	}

	isBot := detectBot(ua)
	browser, browserVersion := detectBrowser(ua)
	os, osVersion := detectOS(ua)
	device := detectDevice(ua, isBot)

	return Result{
		Browser:        browser,
		BrowserVersion: browserVersion,
		OS:             os,
		OSVersion:      osVersion,
		Device:         device,
		IsBot:          isBot,
	}
}

func detectBot(ua string) bool {
	lower := strings.ToLower(ua)
	for _, pattern := range botPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
}

func detectBrowser(ua string) (name, version string) {
	// Order matters: Edge must come before Chrome; Opera before Chrome.
	switch {
	case strings.Contains(ua, "Edg/") || strings.Contains(ua, "Edge/"):
		return "Edge", extractVersion(reEdgeVer, ua)
	case strings.Contains(ua, "OPR/"):
		return "Opera", extractVersion(reOperaVer, ua)
	case strings.Contains(ua, "Firefox/"):
		return "Firefox", extractVersion(reFirefoxVer, ua)
	case strings.Contains(ua, "Chrome/"):
		return "Chrome", extractVersion(reChromeVer, ua)
	case strings.Contains(ua, "Safari/") && strings.Contains(ua, "Version/"):
		return "Safari", extractVersion(reSafariVer, ua)
	case strings.Contains(ua, "Trident/") || strings.Contains(ua, "MSIE "):
		return "Internet Explorer", extractVersion(reIEVer, ua)
	default:
		return "Other", ""
	}
}

func detectOS(ua string) (name, version string) {
	switch {
	case strings.Contains(ua, "Android"):
		return "Android", normalizeVersion(extractVersion(reAndroidVer, ua))
	case strings.Contains(ua, "iPhone") || strings.Contains(ua, "iPad"):
		raw := extractVersion(reiOSVer, ua)
		return "iOS", normalizeVersion(raw)
	case strings.Contains(ua, "Windows NT"):
		nt := extractVersion(reWindowsVer, ua)
		if friendly, ok := windowsVersions[nt]; ok {
			return "Windows", friendly
		}
		return "Windows", nt
	case strings.Contains(ua, "Mac OS X"):
		return "macOS", normalizeVersion(extractVersion(reMacVer, ua))
	case strings.Contains(ua, "Linux"):
		return "Linux", ""
	case strings.Contains(ua, "CrOS"):
		return "ChromeOS", ""
	default:
		return "Other", ""
	}
}

func detectDevice(ua string, isBot bool) string {
	if isBot {
		return "bot"
	}
	lower := strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "iPad") || (strings.Contains(ua, "Android") && !strings.Contains(lower, "mobile")):
		return "tablet"
	case strings.Contains(lower, "mobile") || strings.Contains(ua, "iPhone"):
		return "mobile"
	default:
		return "desktop"
	}
}

func extractVersion(re *regexp.Regexp, ua string) string {
	m := re.FindStringSubmatch(ua)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

// normalizeVersion replaces underscores with dots (e.g. iOS "14_0" → "14.0").
func normalizeVersion(v string) string {
	return strings.ReplaceAll(v, "_", ".")
}
