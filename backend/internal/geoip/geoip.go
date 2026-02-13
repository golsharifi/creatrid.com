package geoip

import (
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

// Service wraps a MaxMind GeoIP2 database reader.
type Service struct {
	db *geoip2.Reader
}

// New creates a GeoIP service. Returns nil service (not error) if path is empty.
func New(dbPath string) (*Service, error) {
	if dbPath == "" {
		return nil, nil
	}
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &Service{db: db}, nil
}

// Result holds the country and city resolved from an IP address.
type Result struct {
	Country string
	City    string
}

// Lookup resolves country and city from an IP address string.
// Returns an empty Result if the service is nil, the IP is unparseable, or the
// lookup fails.
func (s *Service) Lookup(ipStr string) Result {
	if s == nil || s.db == nil {
		return Result{}
	}
	// Strip port if present
	host := ipStr
	if strings.LastIndex(ipStr, ":") != -1 {
		h, _, err := net.SplitHostPort(ipStr)
		if err == nil {
			host = h
		}
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return Result{}
	}
	record, err := s.db.City(ip)
	if err != nil {
		return Result{}
	}
	return Result{
		Country: record.Country.Names["en"],
		City:    record.City.Names["en"],
	}
}

// Close releases the underlying database resources.
func (s *Service) Close() error {
	if s != nil && s.db != nil {
		return s.db.Close()
	}
	return nil
}
