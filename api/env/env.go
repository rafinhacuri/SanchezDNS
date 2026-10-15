//nolint:gochecknoglobals
package env

import (
	"os"
	"strings"
)

type EnvConfig struct {
	Production  bool
	SiteUrl     string
	MongoUrl    string
	RedisUrl    string
	DnsHost     string
	DnsApiKey   string
	DnsServerId string
	FsUser      string
	FsPassword  string
	FsBucket    string
	FsEndpoint  string
	DevUrl      string
	DevKey      string
	DevCert     string
	GeoIPURL    string
	DohResolver []string
}

var C *EnvConfig

func LoadEnv() {
	C = &EnvConfig{
		Production:  os.Getenv("NUXT_PUBLIC_PRODUCTION") == "true",
		SiteUrl:     os.Getenv("NUXT_PUBLIC_SITE_URL"),
		MongoUrl:    os.Getenv("MONGO_URL"),
		RedisUrl:    os.Getenv("REDIS_URL"),
		FsUser:      os.Getenv("FS_USERNAME"),
		FsPassword:  os.Getenv("FS_PASSWORD"),
		FsBucket:    os.Getenv("FS_BUCKET"),
		FsEndpoint:  os.Getenv("FS_ENDPOINT"),
		DnsHost:     os.Getenv("DNS_HOST"),
		DnsApiKey:   os.Getenv("DNS_API_KEY"),
		DnsServerId: os.Getenv("DNS_SERVER_ID"),
		DevUrl:      os.Getenv("DEV_URL"),
		DevKey:      os.Getenv("DEV_KEY"),
		DevCert:     os.Getenv("DEV_CERT"),
		GeoIPURL:    strings.TrimSpace(os.Getenv("GEOIP_URL")),
		DohResolver: lista(os.Getenv("DOH_RESOLVERS")),
	}
}

func lista(valor string) []string {
	var itens []string

	for item := range strings.SplitSeq(valor, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			itens = append(itens, item)
		}
	}

	return itens
}
