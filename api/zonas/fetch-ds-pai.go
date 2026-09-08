package zonas

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"resty.dev/v3"
)

const tipoDs = 43

//nolint:gochecknoglobals
var resolversDoh = []string{
	"https://dns.google/resolve",
	"https://cloudflare-dns.com/dns-query",
}

func FetchDsPai(ctx context.Context, domain string) ([]string, bool, error) {
	nome := strings.TrimSuffix(strings.TrimSpace(domain), ".")
	if nome == "" {
		return nil, false, errors.New("doh: domínio vazio")
	}

	client := resty.New().
		SetTimeout(5*time.Second).
		SetHeader("Accept", "application/dns-json").
		SetRetryCount(1).
		SetRetryWaitTime(200 * time.Millisecond)

	var ultimoErr error

	for _, resolver := range resolversDoh {
		var res dohResponse

		resp, err := client.R().
			SetContext(ctx).
			SetQueryParams(map[string]string{"name": nome, "type": "DS", "do": "1"}).
			SetResult(&res).
			Get(resolver)
		if err != nil {
			ultimoErr = err

			continue
		}

		if resp.IsStatusFailure() {
			ultimoErr = errors.New("doh: " + resp.Status() + " em " + resolver)

			continue
		}

		if res.Status != 0 && res.Status != 3 {
			ultimoErr = fmt.Errorf("doh: rcode %d em %s", res.Status, resolver)

			continue
		}

		ds := make([]string, 0, len(res.Answer))

		for _, answer := range res.Answer {
			if answer.Type == tipoDs && strings.TrimSpace(answer.Data) != "" {
				ds = append(ds, strings.TrimSpace(answer.Data))
			}
		}

		return ds, res.AD, nil
	}

	return nil, false, ultimoErr
}
