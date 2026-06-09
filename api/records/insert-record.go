//nolint:tagliatelle, contextcheck
package records

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
	"github.com/rafinhacuri/SanchezDNS/api/logs"
)

type Record struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
}

type Comment struct {
	Content    string `json:"content"`
	Account    string `json:"account"`
	ModifiedAt int64  `json:"modified_at,omitempty"`
}

type rrsetRecord struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	TTL      int       `json:"ttl"`
	Comments []Comment `json:"comments"`
	Records  []Record  `json:"records"`
}
type Zone struct {
	Name         string        `json:"name"`
	RRSets       []rrsetRecord `json:"rrsets"`
	Serial       int64         `json:"serial"`
	EditedSerial int64         `json:"edited_serial"`
}

func InsertRecord(
	ctx context.Context,
	zona,
	tipo,
	name,
	vl string,
	ttl int,
	comment string,
	svcPriority *int,
	targetName,
	svcParams string,
	weight,
	port *int,
	target string,
	priority *int,
	user string,
) (string, error) {
	httpc := resty.New().
		SetTimeout(30*time.Second).
		SetBaseURL(env.C.DnsHost).
		SetHeader("X-API-Key", env.C.DnsApiKey).
		SetHeader("Accept", "application/json").
		SetRetryCount(2)

	if tipo != "HTTPS" && tipo != "SRV" && (vl == "" || strings.TrimSpace(vl) == "") {
		return "", errors.New("valor do registro é obrigatório para o tipo " + tipo)
	}

	res := normalizeRecordValue(tipo, vl, priority, weight, port, target, svcPriority, targetName, svcParams)

	zone := strings.TrimSuffix(zona, ".")

	if !strings.HasSuffix(name, ".") {
		if !strings.HasSuffix(name, zone) {
			name = fmt.Sprintf("%s.%s.", name, zone)
		} else {
			name += "."
		}
	}

	commentByContent, recordsStep1, ttl, err := prepareInsertStep1(ctx, httpc, zona, tipo, name, res, ttl)
	if err != nil {
		return "", err
	}

	patchBody1 := PDNSZonePatchRequest{
		RRSets: []PDNSRRSetChange{
			{
				Name:       name,
				Type:       tipo,
				TTL:        &ttl,
				ChangeType: "REPLACE",
				Records:    recordsStep1,
			},
		},
	}

	resp, err := httpc.R().SetContext(ctx).SetBody(patchBody1).
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, zona))
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 204 && resp.StatusCode() != 201 {
		return "", errors.New("falha ao adicionar registro: " + resp.String())
	}

	getResp2, err := httpc.R().SetContext(ctx).Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, zona))
	if err != nil {
		return "", errors.New("falha ao obter zona: " + err.Error())
	}

	var zoneData2 Zone

	err = json.Unmarshal(getResp2.Bytes(), &zoneData2)
	if err != nil {
		return "", errors.New("falha ao analisar resposta do PowerDNS: " + err.Error())
	}

	var finalRR *rrsetRecord

	for i := range zoneData2.RRSets {
		rr := &zoneData2.RRSets[i]
		if rr.Type == tipo && rr.Name == name {
			finalRR = rr

			break
		}
	}

	if finalRR == nil {
		return "", errors.New("RRSet não encontrado após inserção")
	}

	now := time.Now().Unix()

	var (
		finalRecords  []Record
		finalComments []Comment
	)

	for _, rec := range finalRR.Records {
		finalRecords = append(finalRecords, Record{Content: rec.Content, Disabled: rec.Disabled})

		var commentText string
		if rec.Content == res {
			commentText = comment
		} else if existing, ok := commentByContent[rec.Content]; ok {
			commentText = existing
		} else {
			commentText = ""
		}

		finalComments = append(finalComments, Comment{
			Content:    commentText,
			Account:    "",
			ModifiedAt: now,
		})
	}

	ttl2 := finalRR.TTL
	patchBody2 := PDNSZonePatchRequest{
		RRSets: []PDNSRRSetChange{
			{
				Name:       name,
				Type:       tipo,
				TTL:        &ttl2,
				ChangeType: "REPLACE",
				Records:    finalRecords,
				Comments:   finalComments,
			},
		},
	}

	resp2, err := httpc.R().SetContext(ctx).SetBody(patchBody2).
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, zona))
	if err != nil {
		return "", errors.New("falha ao adicionar comentário do registro: " + err.Error())
	}

	if resp2.StatusCode() != 204 && resp2.StatusCode() != 201 {
		return "", errors.New("falha ao adicionar comentário do registro: " + resp2.String())
	}

	err = ensureReverseRecord(ctx, tipo, res, zona, name)
	if err != nil {
		log.Println("falha ao garantir registro reverso:", err)
	}

	go logs.InsertLog(
		name,
		user,
		"insert_record",
		fmt.Sprintf("Criado registro %s do tipo %s na zona %s", name, tipo, zona))

	return "record inserted", nil
}

func prepareInsertStep1(
	ctx context.Context,
	httpc *resty.Client,
	zona, tipo, name, res string,
	ttl int,
) (map[string]string, []Record, int, error) {
	zoneData, err := fetchZoneInsert(ctx, httpc, zona)
	if err != nil {
		return nil, nil, 0, err
	}

	existingRR := findRRSet(zoneData, tipo, name)
	commentByContent := buildCommentByContentInsert(existingRR)

	recordsStep1, effectiveTTL := buildRecordsStep1(existingRR, res, ttl)

	return commentByContent, recordsStep1, effectiveTTL, nil
}

func fetchZoneInsert(ctx context.Context, httpc *resty.Client, zona string) (*Zone, error) {
	getResp, err := httpc.R().SetContext(ctx).
		Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, zona))
	if err != nil {
		return nil, err
	}

	var zoneData Zone

	err = json.Unmarshal(getResp.Bytes(), &zoneData)
	if err != nil {
		return nil, err
	}

	return &zoneData, nil
}

func findRRSet(zoneData *Zone, tipo, name string) *rrsetRecord {
	if zoneData == nil {
		return nil
	}

	for i := range zoneData.RRSets {
		rr := &zoneData.RRSets[i]
		if rr.Type == tipo && rr.Name == name {
			return rr
		}
	}

	return nil
}

func buildCommentByContentInsert(existingRR *rrsetRecord) map[string]string {
	commentByContent := make(map[string]string)
	if existingRR == nil {
		return commentByContent
	}

	for i, rec := range existingRR.Records {
		if i < len(existingRR.Comments) {
			commentByContent[rec.Content] = existingRR.Comments[i].Content
		}
	}

	return commentByContent
}

func buildRecordsStep1(existingRR *rrsetRecord, res string, ttl int) ([]Record, int) {
	recordsStep1 := make([]Record, 0, func() int {
		if existingRR == nil {
			return 1
		}

		return len(existingRR.Records) + 1
	}())

	if existingRR != nil {
		for _, rec := range existingRR.Records {
			recordsStep1 = append(recordsStep1, Record{Content: rec.Content, Disabled: rec.Disabled})
		}

		if ttl <= 0 {
			ttl = existingRR.TTL
		}
	}

	recordsStep1 = append(recordsStep1, Record{Content: res, Disabled: false})

	return recordsStep1, ttl
}
