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

type RecordChange struct {
	Zona        string
	Tipo        string
	Name        string
	VL          string
	TTL         int
	Comment     string
	SvcPriority *int
	TargetName  string
	SvcParams   string
	Weight      *int
	Port        *int
	Target      string
	Priority    *int
}

func buildDNSClient() *resty.Client {
	return resty.New().
		SetTimeout(30*time.Second).
		SetBaseURL(env.C.DnsHost).
		SetHeader("X-API-Key", env.C.DnsApiKey).
		SetHeader("Accept", "application/json").
		SetTimeout(6 * time.Second).
		SetRetryCount(2)
}

func normalizeRRSetName(zoneFQDN, name string) string {
	zone := strings.TrimSuffix(zoneFQDN, ".")
	out := name

	if !strings.HasSuffix(out, ".") {
		if !strings.HasSuffix(out, zone) {
			out = fmt.Sprintf("%s.%s.", out, zone)
		} else {
			out += "."
		}
	}

	return out
}

func fetchZone(ctx context.Context, httpc *resty.Client, zoneFQDN string) (Zone, error) {
	resp, err := httpc.R().
		SetContext(ctx).
		Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, zoneFQDN))
	if err != nil {
		return Zone{}, err
	}

	var z Zone

	err = json.Unmarshal(resp.Bytes(), &z)
	if err != nil {
		return Zone{}, errors.New("falha ao analisar resposta do PowerDNS: " + err.Error())
	}

	return z, nil
}

func findRRSetByNameAndType(zone Zone, rrType, rrName string) *rrsetRecord {
	for i := range zone.RRSets {
		rr := &zone.RRSets[i]
		if rr.Type == rrType && rr.Name == rrName {
			return rr
		}
	}

	return nil
}

func buildCommentByContent(rr *rrsetRecord) map[string]string {
	commentByContent := make(map[string]string, len(rr.Records))

	for i, rec := range rr.Records {
		if i < len(rr.Comments) {
			commentByContent[rec.Content] = rr.Comments[i].Content
		}
	}

	return commentByContent
}

func buildRecordsReplacingContent(rr *rrsetRecord, oldContent, newContent string) []Record {
	out := make([]Record, 0, len(rr.Records))

	for _, rec := range rr.Records {
		content := rec.Content
		if rec.Content == oldContent {
			content = newContent
		}

		out = append(out, Record{
			Content:  content,
			Disabled: rec.Disabled,
		})
	}

	return out
}

func patchZoneRRSet(ctx context.Context, httpc *resty.Client, zoneFQDN string, body PDNSZonePatchRequest) error {
	resp, err := httpc.R().
		SetContext(ctx).
		SetBody(body).
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, zoneFQDN))
	if err != nil {
		return err
	}

	if resp.StatusCode() != 204 && resp.StatusCode() != 201 {
		return errors.New("falha ao editar registro: " + resp.String())
	}

	return nil
}

func buildFinalRecordsAndComments(
	rr *rrsetRecord,
	newContent,
	newComment string,
	commentByContent map[string]string,
	now int64,
) ([]Record, []Comment) {
	records := make([]Record, 0, len(rr.Records))
	comments := make([]Comment, 0, len(rr.Records))

	for _, rec := range rr.Records {
		records = append(records, Record{
			Content:  rec.Content,
			Disabled: rec.Disabled,
		})

		commentText := ""
		if rec.Content == newContent {
			commentText = newComment
		} else if existing, ok := commentByContent[rec.Content]; ok {
			commentText = existing
		}

		comments = append(comments, Comment{
			Content:    commentText,
			Account:    "",
			ModifiedAt: now,
		})
	}

	return records, comments
}

func UpdateRecord(ctx context.Context, newRec, oldRec RecordChange, user string) (string, error) {
	httpc := buildDNSClient()

	resNew := normalizeRecordValue(
		newRec.Tipo,
		newRec.VL,
		newRec.Priority,
		newRec.Weight,
		newRec.Port,
		newRec.Target,
		newRec.SvcPriority,
		newRec.TargetName,
		newRec.SvcParams,
	)

	resOld := normalizeRecordValue(
		oldRec.Tipo,
		oldRec.VL,
		oldRec.Priority,
		oldRec.Weight,
		oldRec.Port,
		oldRec.Target,
		oldRec.SvcPriority,
		oldRec.TargetName,
		oldRec.SvcParams,
	)

	name := normalizeRRSetName(newRec.Zona, newRec.Name)

	zoneData, err := fetchZone(ctx, httpc, newRec.Zona)
	if err != nil {
		return "", err
	}

	existingRR := findRRSetByNameAndType(zoneData, newRec.Tipo, name)
	if existingRR == nil {
		return "", errors.New("RRSet não encontrado para edição")
	}

	commentByContent := buildCommentByContent(existingRR)

	recordsStep1 := buildRecordsReplacingContent(existingRR, resOld, resNew)

	ttl1 := existingRR.TTL
	if newRec.TTL > 0 {
		ttl1 = newRec.TTL
	}

	patchBody1 := PDNSZonePatchRequest{
		RRSets: []PDNSRRSetChange{
			{
				Name:       name,
				Type:       newRec.Tipo,
				TTL:        &ttl1,
				ChangeType: "REPLACE",
				Records:    recordsStep1,
			},
		},
	}

	err = patchZoneRRSet(ctx, httpc, newRec.Zona, patchBody1)
	if err != nil {
		return "", err
	}

	zoneData2, err := fetchZone(ctx, httpc, newRec.Zona)
	if err != nil {
		return "", errors.New("falha ao obter zona após edição: " + err.Error())
	}

	finalRR := findRRSetByNameAndType(zoneData2, newRec.Tipo, name)
	if finalRR == nil {
		return "", errors.New("RRSet não encontrado após edição")
	}

	now := time.Now().Unix()
	finalRecords, finalComments := buildFinalRecordsAndComments(finalRR, resNew, newRec.Comment, commentByContent, now)

	ttl2 := finalRR.TTL
	patchBody2 := PDNSZonePatchRequest{
		RRSets: []PDNSRRSetChange{
			{
				Name:       name,
				Type:       newRec.Tipo,
				TTL:        &ttl2,
				ChangeType: "REPLACE",
				Records:    finalRecords,
				Comments:   finalComments,
			},
		},
	}

	resp2, err := httpc.R().SetContext(ctx).SetBody(patchBody2).
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, newRec.Zona))
	if err != nil {
		return "", errors.New("falha ao aplicar comentários do registro editado: " + err.Error())
	}

	if resp2.StatusCode() != 204 && resp2.StatusCode() != 201 {
		return "", errors.New("falha ao aplicar comentários do registro editado: " + resp2.String())
	}

	err = ensureReverseRecordUpdate(ctx, newRec.Tipo, resNew, newRec.Zona, name, resOld)
	if err != nil {
		return "", err
	}

	go logs.InsertLog(
		newRec.Name,
		user,
		"edit_record",
		fmt.Sprintf("Editado registro %s do tipo %s na zona %s", newRec.Name, newRec.Tipo, newRec.Zona),
	)

	return "record updated", nil
}
