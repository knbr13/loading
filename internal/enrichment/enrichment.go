package enrichment

import (
	"embed"
	"net"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

//go:embed assets/geoip.mmdb
var assets embed.FS

type EnrichmentResult struct {
	Country string
	Domain  string
}

type Enricher struct {
	geoDB    *geoip2.Reader
	dnsCache sync.Map
	geoCache sync.Map
}

func NewEnricher() (*Enricher, error) {
	dbData, _ := assets.ReadFile("assets/geoip.mmdb")
	var db *geoip2.Reader
	var err error
	if len(dbData) > 0 {
		db, err = geoip2.FromBytes(dbData)
		if err != nil {
			// Don't fail completely if DB is invalid
			db = nil
		}
	}

	return &Enricher{
		geoDB: db,
	}, nil
}

func (e *Enricher) LookupCountry(ip net.IP) string {
	if e.geoDB == nil {
		return "Unknown"
	}

	if val, ok := e.geoCache.Load(ip.String()); ok {
		return val.(string)
	}

	record, err := e.geoDB.Country(ip)
	if err != nil {
		return "Unknown"
	}

	country := record.Country.Names["en"]
	if country == "" {
		country = "Unknown"
	}
	e.geoCache.Store(ip.String(), country)
	return country
}

func (e *Enricher) LookupDNS(ip net.IP) string {
	ipStr := ip.String()
	if val, ok := e.dnsCache.Load(ipStr); ok {
		return val.(string)
	}

	// For now return IP string and start async lookup
	go func() {
		names, err := net.LookupAddr(ipStr)
		if err == nil && len(names) > 0 {
			e.dnsCache.Store(ipStr, names[0])
		}
	}()

	return ipStr
}

func (e *Enricher) Close() {
	if e.geoDB != nil {
		e.geoDB.Close()
	}
}
