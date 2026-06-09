//nolint:contextcheck
package records

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
	"github.com/rafinhacuri/SanchezDNS/api/logs"
)

type PDNSZonePatchRequest struct {
	RRSets []PDNSRRSetChange `json:"rrsets"`
}

type PDNSRRSetChange struct {
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	TTL        *int      `json:"ttl,omitempty"`
	ChangeType string    `json:"changetype"`
	Records    []Record  `json:"records,omitempty"`
	Comments   []Comment `json:"comments,omitempty"`
}

func newPDNSClient() *resty.Client {
	return resty.New().
		SetTimeout(30*time.Second).
		SetBaseURL(env.C.DnsHost).
		SetHeader("X-API-Key", env.C.DnsApiKey).
		SetHeader("Accept", "application/json").
		SetRetryCount(2)
}

func normalizeDeleteName(zona, name string) string {
	zone := strings.TrimSuffix(zona, ".")
	if !strings.HasSuffix(name, ".") {
		if !strings.HasSuffix(name, zone) {
			name = fmt.Sprintf("%s.%s.", name, zone)
		} else {
			name += "."
		}
	}

	return name
}

func mxValueWithoutPriority(vl string) (string, error) {
	parts := strings.Fields(vl)
	if len(parts) < 2 {
		return "", errors.New("valor MX inválido")
	}

	return strings.Join(parts[1:], " "), nil
}

func fetchZoneDelete(ctx context.Context, httpc *resty.Client, zona string) (Zone, error) {
	getResp, err := httpc.R().
		SetContext(ctx).
		Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, zona))
	if err != nil {
		return Zone{}, errors.New("falha ao obter zona: " + err.Error())
	}

	var zoneData Zone

	err = json.Unmarshal(getResp.Bytes(), &zoneData)
	if err != nil {
		return Zone{}, err
	}

	return zoneData, nil
}

func remainingRRSetForDelete(zoneData Zone, tipo, name, res string) ([]Record, []Comment) {
	var (
		remainingRecords  []Record
		remainingComments []Comment
	)

	for _, rr := range zoneData.RRSets {
		if rr.Type != tipo || rr.Name != name {
			continue
		}

		for i, rec := range rr.Records {
			if rec.Content == res {
				continue
			}

			remainingRecords = append(remainingRecords, Record{
				Content:  rec.Content,
				Disabled: rec.Disabled,
			})

			text, account, modified := findCommentForRecordContent(rr.Comments, i, len(rr.Records))
			remainingComments = append(remainingComments, Comment{
				Content:    text,
				Account:    account,
				ModifiedAt: modified,
			})
		}
	}

	return remainingRecords, remainingComments
}

func patchDeleteRRSet(ctx context.Context, httpc *resty.Client, zona, name, tipo string) (*resty.Response, error) {
	return httpc.R().
		SetContext(ctx).
		SetBody(PDNSZonePatchRequest{
			RRSets: []PDNSRRSetChange{
				{
					Name:       name,
					Type:       tipo,
					ChangeType: "DELETE",
				},
			},
		}).
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, zona))
}

func patchReplaceRRSet(
	ctx context.Context,
	httpc *resty.Client,
	zona string,
	name string,
	tipo string,
	ttl int,
	records []Record,
	comments []Comment,
) (*resty.Response, error) {
	return httpc.R().
		SetContext(ctx).
		SetBody(PDNSZonePatchRequest{
			RRSets: []PDNSRRSetChange{
				{
					Name:       name,
					Type:       tipo,
					TTL:        &ttl,
					ChangeType: "REPLACE",
					Records:    records,
					Comments:   comments,
				},
			},
		}).
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, zona))
}

func validatePDNSPatchResponse(resp *resty.Response) error {
	if resp == nil {
		return errors.New("falha ao excluir registro: resposta vazia")
	}

	if resp.StatusCode() != 204 && resp.StatusCode() != 201 {
		return errors.New("falha ao excluir registro: " + resp.String())
	}

	return nil
}

func deleteReverseAndLog(ctx context.Context, tipo, vl, name, user, zona string) error {
	err := deleteReverseRecord(ctx, tipo, vl)
	if err != nil {
		return errors.New("falha ao excluir registro reverso: " + err.Error())
	}

	go logs.InsertLog(
		name,
		user,
		"delete_record",
		fmt.Sprintf("Excluído registro %s do tipo %s na zona %s", name, tipo, zona))

	return nil
}

func DeleteRecord(
	ctx context.Context,
	zona string,
	tipo string,
	name string,
	vl string,
	ttl int,
	comment string,
	svcPriority *int,
	targetName string,
	svcParams string,
	weight *int,
	port *int,
	target string,
	priority *int,
	user string,
) (string, error) {
	httpc := newPDNSClient()

	if tipo == "MX" {
		mx, err := mxValueWithoutPriority(vl)

		vl = mx

		if err != nil {
			return "", err
		}
	}

	res := normalizeRecordValue(tipo, vl, priority, weight, port, target, svcPriority, targetName, svcParams)
	name = normalizeDeleteName(zona, name)

	zoneData, err := fetchZoneDelete(ctx, httpc, zona)
	if err != nil {
		return "", err
	}

	remainingRecords, remainingComments := remainingRRSetForDelete(zoneData, tipo, name, res)

	if len(remainingRecords) == 0 {
		resp, err := patchDeleteRRSet(ctx, httpc, zona, name, tipo)
		if err != nil {
			return "", err
		}

		err = validatePDNSPatchResponse(resp)
		if err != nil {
			return "", err
		}

		err = deleteReverseAndLog(ctx, tipo, vl, name, user, zona)
		if err != nil {
			return "", err
		}

		return "record deleted", nil
	}

	resp, err := patchReplaceRRSet(ctx, httpc, zona, name, tipo, ttl, remainingRecords, remainingComments)
	if err != nil {
		return "", errors.New("falha ao excluir registro: " + err.Error())
	}

	err = validatePDNSPatchResponse(resp)
	if err != nil {
		return "", err
	}

	err = deleteReverseAndLog(ctx, tipo, vl, name, user, zona)
	if err != nil {
		return "", err
	}

	return "record deleted", nil
}
