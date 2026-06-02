//nolint:gochecknoglobals
package env

import (
	"os"
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
	DevUrl      string
	DevKey      string
	DevCert     string
}

var C *EnvConfig

func LoadEnv() {
	C = &EnvConfig{
		Production:  os.Getenv("NUXT_PUBLIC_PRODUCTION") == "true",
		SiteUrl:     os.Getenv("NUXT_PUBLIC_SITE_URL"),
		MongoUrl:    os.Getenv("MONGO_URL"),
		RedisUrl:    os.Getenv("REDIS_URL"),
		FsUser:      os.Getenv("FS_USER"),
		FsPassword:  os.Getenv("FS_PASSWORD"),
		FsBucket:    os.Getenv("FS_BUCKET"),
		DnsHost:     os.Getenv("DNS_HOST"),
		DnsApiKey:   os.Getenv("DNS_API_KEY"),
		DnsServerId: os.Getenv("DNS_SERVER_ID"),
		DevUrl:      os.Getenv("DEV_URL"),
		DevKey:      os.Getenv("DEV_KEY"),
		DevCert:     os.Getenv("DEV_CERT"),
	}
}
