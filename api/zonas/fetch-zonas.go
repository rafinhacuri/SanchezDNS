package zonas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

type ZoneFetch struct {
	Name   string `json:"name"`
	Serial int64  `json:"serial"`
	Nivel  string `json:"nivel"`
	Dnssec bool   `json:"dnssec"`
}

type ZonesResponse struct {
	Zones []ZoneFetch `json:"zones"`
}

type pdnsZoneListItem struct {
	Name   string `json:"name"`
	Serial int64  `json:"serial"`
	Dnssec *bool  `json:"dnssec,omitempty"`
}

type pdnsZoneDetail struct {
	Dnssec bool `json:"dnssec"`
}

type ZonePermission struct {
	Zona    string   `bson:"zona"`
	Leitura []string `bson:"leitura"`
	Escrita []string `bson:"escrita"`
}

func containsCI(slice []string, value string) bool {
	v := strings.ToLower(strings.TrimSpace(value))
	for _, s := range slice {
		if strings.ToLower(strings.TrimSpace(s)) == v {
			return true
		}
	}

	return false
}

func resolveNivel(level, email, zone string, permsByZone map[string]ZonePermission) (string, bool) {
	if level == "admin" {
		return "ADMINISTRADOR", true
	}

	perm, ok := permsByZone[zone]
	if !ok {
		return "", false
	}

	if containsCI(perm.Escrita, email) {
		return "ESCRITA", true
	}

	if containsCI(perm.Leitura, email) {
		return "LEITURA", true
	}

	return "", false
}

func zoneMatchesType(zoneNameLower, zoneType string) bool {
	isV4 := strings.HasSuffix(zoneNameLower, ".in-addr.arpa.")
	isV6 := strings.HasSuffix(zoneNameLower, ".ip6.arpa.")

	switch zoneType {
	case "normal":
		return !isV4 && !isV6
	case "reverse":
		return isV4
	case "reverse-ipv6":
		return isV6
	default:
		// If an unknown filter is passed, behave like "normal"
		return !isV4 && !isV6
	}
}

func fetchZonesList(ctx context.Context, httpc *resty.Client) ([]pdnsZoneListItem, error) {
	var zones []pdnsZoneListItem

	resp, err := httpc.R().
		SetContext(ctx).
		SetResult(&zones).
		Get(fmt.Sprintf("/api/v1/servers/%s/zones", env.C.DnsServerId))
	if err != nil {
		log.Println(err.Error())

		return nil, errors.New("falha ao buscar zonas")
	}

	if resp.IsError() {
		return nil, errors.New("falha ao buscar zonas")
	}

	return zones, nil
}

func loadPermissionsByZone(ctx context.Context) (map[string]ZonePermission, error) {
	coll := mongo.Dns.Collection("users")

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		log.Println(err.Error())

		return nil, errors.New("falha ao buscar permissões de zonas")
	}

	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			log.Println(err.Error())
		}
	}()

	permsByZone := map[string]ZonePermission{}

	for cursor.Next(ctx) {
		var p ZonePermission

		err := cursor.Decode(&p)
		if err != nil {
			log.Println(err.Error())

			return nil, errors.New("falha ao decodificar permissões de zonas")
		}

		zonaLower := strings.ToLower(strings.TrimSpace(p.Zona))
		permsByZone[zonaLower] = p
	}

	err = cursor.Err()
	if err != nil {
		log.Println(err.Error())

		return nil, errors.New("erro ao iterar permissões de zonas")
	}

	return permsByZone, nil
}

func fetchZoneDNSSEC(ctx context.Context, httpc *resty.Client, z pdnsZoneListItem) (bool, error) {
	if z.Dnssec != nil {
		return *z.Dnssec, nil
	}

	var detail pdnsZoneDetail

	zoneURL := fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, url.PathEscape(z.Name))

	resp, err := httpc.R().
		SetContext(ctx).
		SetResult(&detail).
		Get(zoneURL)
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, errors.New("pdns error")
	}

	return detail.Dnssec, nil
}

func filterReverseSpecial(zones []ZoneFetch) []ZoneFetch {
	finalFiltered := make([]ZoneFetch, 0, len(zones))

	for _, z := range zones {
		if z.Name == "1.1.1.in-addr.arpa." || z.Name == "2.2.2.in-addr.arpa." {
			continue
		}

		finalFiltered = append(finalFiltered, z)
	}

	return finalFiltered
}

func FetchZonas(ctx context.Context, zoneType, email, level string) ([]ZoneFetch, error) {
	httpc := resty.New().
		SetTimeout(30*time.Second).
		SetBaseURL(env.C.DnsHost).
		SetHeader("X-API-Key", env.C.DnsApiKey).
		SetHeader("Accept", "application/json").
		SetRetryCount(2)

	zones, err := fetchZonesList(ctx, httpc)
	if err != nil {
		return []ZoneFetch{}, err
	}

	permsByZone := map[string]ZonePermission{}

	if level != "admin" {
		permsByZone, err = loadPermissionsByZone(ctx)
		if err != nil {
			return []ZoneFetch{}, err
		}
	}

	filtered := make([]ZoneFetch, 0, len(zones))

	for _, z := range zones {
		nameLower := strings.ToLower(strings.TrimSpace(z.Name))

		if !zoneMatchesType(nameLower, zoneType) {
			continue
		}

		nivel, ok := resolveNivel(level, email, nameLower, permsByZone)
		if !ok {
			continue
		}

		dnssec, derr := fetchZoneDNSSEC(ctx, httpc, z)
		if derr != nil {
			log.Println(derr.Error())

			return []ZoneFetch{}, errors.New("falha ao buscar detalhes da zona")
		}

		filtered = append(filtered, ZoneFetch{
			Name:   z.Name,
			Serial: z.Serial,
			Nivel:  nivel,
			Dnssec: dnssec,
		})
	}

	if zoneType == "reverse" {
		return filterReverseSpecial(filtered), nil
	}

	return filtered, nil
}
